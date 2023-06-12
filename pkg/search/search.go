/*
 * @Author: cloudyi.li
 * @Date: 2023-06-01 08:53:25
 * @LastEditTime: 2023-06-12 19:06:08
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/search/search.go
 */
package search

// 对接Google API。通过query查询关键词
import (
	"chatserver-api/internal/consts"
	"chatserver-api/pkg/cache"
	"chatserver-api/pkg/config"
	"chatserver-api/pkg/logger"
	"chatserver-api/pkg/openai"
	"chatserver-api/utils/security"
	"context"
	"net/http"
	"sync"
	"time"

	customsearch "google.golang.org/api/customsearch/v1"
	"google.golang.org/api/googleapi/transport"
)

func summaryContent(message string) string {
	if message == "" {
		return ""
	}
	var req openai.ChatCompletionRequest
	var chatMessages []openai.ChatCompletionMessage
	var systemPreset, userMessage openai.ChatCompletionMessage
	systemPreset.Role = openai.ChatMessageRoleSystem
	systemPreset.Content = `Extract information is limited to 500 words. Please answer in Chinese and do not add any additional prompts.`
	userMessage.Role = openai.ChatMessageRoleUser
	userMessage.Content = message
	chatMessages = append(chatMessages, systemPreset, userMessage)
	req.Model = "gpt-3.5-turbo"
	req.MaxTokens = 700
	req.Messages = chatMessages
	client, err := openai.NewClient()
	if err != nil {
		logger.Errorf("%s", err)
		return ""
	}
	resp, err := client.CreateChatCompletion(req)
	if err != nil {
		logger.Errorf("%s", err)
		return ""
	}
	content := resp.Choices[0].Message.Content
	return content
}

func queryExtra(message string) string {
	if message == "" {
		return ""
	}
	var req openai.ChatCompletionRequest
	var chatMessages []openai.ChatCompletionMessage
	var systemPreset, userMessage openai.ChatCompletionMessage
	systemPreset.Role = openai.ChatMessageRoleSystem
	systemPreset.Content = `Your task is to rewrite user-submitted content into efficient search engine query statements. Example:''' 
	user: "搜索一下台湾最新的选举情况"
	assistant: "台湾 最新 选举 情况"
	user: "南京今天天气怎么样呢"
	assistant: "南京 今天 天气"
	''' `
	userMessage.Role = openai.ChatMessageRoleUser
	userMessage.Content = message
	chatMessages = append(chatMessages, systemPreset, userMessage)
	req.Model = "gpt-3.5-turbo"
	req.MaxTokens = 100
	req.Messages = chatMessages
	client, err := openai.NewClient()
	if err != nil {
		logger.Errorf("%s", err)
		return ""
	}
	resp, err := client.CreateChatCompletion(req)
	if err != nil {
		logger.Errorf("%s", err)
		return ""
	}
	content := resp.Choices[0].Message.Content
	return content
}

// Message拼装

// 搜索
type searchOne struct {
	Title   string
	Snippet string
	Link    string
}

func CustomSearch(ctx context.Context, queryoriginal string) (string, error) {
	googlecfg := config.AppConfig.GoogelConfig
	ner, keyword := nerDetec(queryoriginal)
	if ner == 0 {
		logger.Debug("没有实体返回")
		return "", nil
	}
	query := queryExtra(queryoriginal)
	rc := cache.GetRedisClient()
	cacheresult, err := rc.Get(ctx, consts.QuerySearchPrefix+security.Md5(query)).Result()
	if err == nil {
		logger.Debug("获取缓存返回")
		return cacheresult, nil
	}

	var searchResult []searchOne
	var textcontent string
	client := &http.Client{Transport: &transport.APIKey{Key: googlecfg.ApiKey}}
	svc, err := customsearch.New(client)
	if err != nil {
		logger.Errorf("%s", err)
		return "", err
	}

	var resp *customsearch.Search

	switch ner {
	case 2:
		{
			resp, err = svc.Cse.List().Cx(googlecfg.CxId).Num(8).Sort("date").Cr("zh-CN").DateRestrict("d[2]").ExactTerms(keyword).Q(query).Do()
			if err != nil {
				logger.Errorf("%s", err)
				return "", err
			}
		}
	default:
		{
			resp, err = svc.Cse.List().Cx(googlecfg.CxId).Num(8).Cr("zh-CN").DateRestrict("y[3]").ExactTerms(keyword).Sort("date").Q(query).Do()
			if err != nil {
				logger.Errorf("%s", err)
				return "", err
			}
		}
	}

	for _, result := range resp.Items {
		searchone := searchOne{}
		searchone.Title = result.Title
		searchone.Snippet = result.Snippet
		searchone.Link = result.Link
		searchResult = append(searchResult, searchone)
	}

	lenresult := len(searchResult)
	if lenresult == 0 {
		return "", nil
	}
	if lenresult > 3 {
		lenresult = 3
	}
	wg := sync.WaitGroup{}

	var lock sync.Mutex

	for _, v := range searchResult[:lenresult-1] {
		wg.Add(1)
		go func(v searchOne) {
			defer wg.Done()
			content := crawlPage(v.Link)
			var summary string
			if content != "" {
				summary = summaryContent("Title:" + v.Title + "\n" + "Link" + v.Link + "\n" + "Content:" + content)
			}
			lock.Lock()
			textcontent += "Title:\n" + v.Title + "\n" + "Snippet:\n" + v.Snippet + "\n" + "Content:\n" + summary + "\n" + "Web Link:\n" + v.Link + "\n"
			lock.Unlock()
		}(v)

	}

	wg.Wait()

	err = rc.Set(ctx, consts.QuerySearchPrefix+security.Md5(query), textcontent, 30*time.Minute).Err()
	return textcontent, err
}
