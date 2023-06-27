/*
 * @Author: cloudyi.li
 * @Date: 2023-06-01 08:53:25
 * @LastEditTime: 2023-06-27 15:08:41
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/search/search.go
 */
package search

// 对接Google API。通过query查询关键词
import (
	"chatserver-api/pkg/config"
	"chatserver-api/pkg/logger"
	"chatserver-api/pkg/openai"
	"context"
	"net/http"
	"strings"
	"sync"

	customsearch "google.golang.org/api/customsearch/v1"
	"google.golang.org/api/googleapi/transport"
)

func SummaryContent(question, message string, retry int) string {
	if message == "" {
		return ""
	}
	var req openai.ChatCompletionRequest
	var chatMessages []openai.ChatCompletionMessage
	var systemPreset, userMessage openai.ChatCompletionMessage
	systemPreset.Role = openai.ChatMessageRoleSystem
	systemStr := `You need to help users organize and filter information obtained through search engines. User will send the '''Search KeyWord''' and '''Search Content''' to you. You are responsible for organizing and summarizing the content related to '''Search KeyWord''' from the '''Search Content'''. If the search keyword is not mentioned in the search content, **must ** reply '''[NO_CONTENT]'''. Answering user's questions in Chinese.`
	systemPreset.Content = systemStr
	userMessage.Role = openai.ChatMessageRoleUser
	userMessage.Content = "## Search KeyWord:" + question + "\n" + "## Search Content" + message
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
		if retry < 3 {
			SummaryContent(question, message, retry+1)
		} else {
			return ""
		}
	}
	if len(resp.Choices) != 0 {
		return resp.Choices[0].Message.Content
	} else {
		return ""
	}
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
	// rc := cache.GetRedisClient()
	// cacheresult, err := rc.Get(ctx, consts.SearchCachePrefix+security.Md5(query)).Result()
	// if err == nil {
	// 	logger.Debug("获取缓存返回")
	// 	return cacheresult, nil
	// }

	var searchResult []searchOne
	var textcontent string
	client := &http.Client{Transport: &transport.APIKey{Key: googlecfg.ApiKey}}
	svc, err := customsearch.New(client)
	if err != nil {
		logger.Errorf("%s", err)
		return "", err
	}

	var resp *customsearch.Search
	if classify == "News" {
		resp, err = svc.Cse.List().Cx(googlecfg.CxId).Num(5).Cr("zh-CN").Sort("date").DateRestrict("y[1]").Q(query).Do()
	} else {
		resp, err = svc.Cse.List().Cx(googlecfg.CxId).Num(5).Cr("zh-CN").Q(query).Do()
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
	if len(searchResult) == 0 {
		return "", nil
	}

	wg := sync.WaitGroup{}
	var lock sync.Mutex
	for _, v := range searchResult {
		wg.Add(1)
		go func(v searchOne) {
			defer wg.Done()
			logger.Debug("[CustomSearch]", logger.Pair("关键词", query), logger.Pair("结果标题", v.Title), logger.Pair("结果链接", v.Link))
			content := crawlPage(v.Link)
			var summary string
			if content != "" {
				summary = SummaryContent(query, "Title:"+v.Title+"\n"+"Snippet"+v.Snippet+"\n"+"Content:"+content, 0)
			}
			lock.Lock()
			if !strings.Contains(summary, "[NO_CONTENT]") {
				textcontent += "Title:\n" + v.Title + "\n" + "\n" + "Snippet" + v.Snippet + "\n" + "Content:\n" + summary + "\n" + "Web Link:\n" + v.Link + "\n"
			}
			lock.Unlock()
		}(v)

	}
	wg.Wait()

	// err = rc.Set(ctx, consts.SearchCachePrefix+security.Md5(query), textcontent, 30*time.Minute).Err()
	return textcontent, err
}
