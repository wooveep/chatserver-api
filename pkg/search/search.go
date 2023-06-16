/*
 * @Author: cloudyi.li
 * @Date: 2023-06-01 08:53:25
 * @LastEditTime: 2023-06-15 22:35:01
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
	systemPreset.Content = `Extract information is limited to 1000 words. Please answer in Chinese and do not add any additional prompts.`
	userMessage.Role = openai.ChatMessageRoleUser
	userMessage.Content = message
	chatMessages = append(chatMessages, systemPreset, userMessage)
	req.Model = "gpt-3.5-turbo-16k-0613"
	req.MaxTokens = 2000
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

func CustomSearch(ctx context.Context, query string, classify string) (string, error) {
	googlecfg := config.AppConfig.GoogelConfig
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

	switch classify {
	case "Entity":
		{
			resp, err = svc.Cse.List().Cx(googlecfg.CxId).Num(3).Cr("zh-CN").Lr("lang_zh-CN").SiteSearch("wikipedia.org").SiteSearchFilter("i").Q(query).Do()
			if err != nil {
				logger.Errorf("%s", err)
				return "", err
			}
		}
	case "News":
		{
			resp, err = svc.Cse.List().Cx(googlecfg.CxId).Num(3).Cr("zh-CN").Lr("lang_zh-CN").DateRestrict("m[1]").Sort("date").Q(query).Do()
			if err != nil {
				logger.Errorf("%s", err)
				return "", err
			}
		}
	default:
		{
			resp, err = svc.Cse.List().Cx(googlecfg.CxId).Num(4).Cr("zh-CN").DateRestrict("y[3]").Sort("date").Q(query).Do()
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
	wg := sync.WaitGroup{}

	var lock sync.Mutex

	for _, v := range searchResult[:lenresult-1] {
		wg.Add(1)
		go func(v searchOne) {
			defer wg.Done()
			content := crawlPage(v.Link)
			// var summary string
			// if content != "" {
			// 	summary = summaryContent("Title:" + v.Title + "\n" + "Link" + v.Link + "\n" + "Content:" + content)
			// }
			lock.Lock()
			textcontent += "Title:\n" + v.Title + "\n" + "Snippet:\n" + v.Snippet + "\n" + "Content:\n" + content + "\n" + "Web Link:\n" + v.Link + "\n"
			lock.Unlock()
		}(v)

	}

	wg.Wait()

	err = rc.Set(ctx, consts.QuerySearchPrefix+security.Md5(query), textcontent, 30*time.Minute).Err()
	return textcontent, err
}
