/*
 * @Author: cloudyi.li
 * @Date: 2023-06-01 08:53:25
 * @LastEditTime: 2023-06-08 15:08:34
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

const (
	apiKey = "AIzaSyB1PJ1jZ5hzGRD5ZB1iJ1EhZEntunS5S4c"
	cx     = "708f239208e3942ee"
	// query  = "南京有哪些好玩的"
)

// 实体检测

// 页面爬取

// 页面摘要提取
func summaryContent(message string) string {
	if message == "" {
		return ""
	}
	var req openai.ChatCompletionRequest
	var chatMessages []openai.ChatCompletionMessage
	var systemPreset, userMessage openai.ChatCompletionMessage
	systemPreset.Role = openai.ChatMessageRoleSystem
	systemPreset.Content = `Summarize the content of the webpage based on its title, extract the main information and list the detailed data included,and keep it within 300 words.please respond in Chinese.`
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

// Message拼装

// 搜索
type searchOne struct {
	Title   string
	Snippet string
	Link    string
}

func CustomSearch(ctx context.Context, query string) (string, error) {
	googlecfg := config.AppConfig.GoogelConfig
	ner, keyword := nerDetec(query)
	if ner == 0 {
		logger.Debug("没有实体返回")
		return "", nil
	}
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
		resp, err = svc.Cse.List().Cx(googlecfg.CxId).Num(10).Sort("date").Cr("zh-CN").DateRestrict("d[2]").ExactTerms(keyword).Q(query).Do()
	default:
		resp, err = svc.Cse.List().Cx(googlecfg.CxId).Num(10).Cr("zh-CN").DateRestrict("y[3]").Sort("date").ExactTerms(keyword).Q(query).Do()
	}

	if err != nil {
		logger.Errorf("%s", err)
		return "", err
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
	if lenresult > 6 {
		lenresult = 6
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
