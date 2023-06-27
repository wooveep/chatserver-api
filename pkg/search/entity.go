/*
 * @Author: cloudyi.li
 * @Date: 2023-06-17 21:58:47
 * @LastEditTime: 2023-06-27 11:08:50
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/search/entity.go
 */
package search

import (
	"chatserver-api/pkg/config"
	"chatserver-api/pkg/logger"
	"context"
	"encoding/json"
	"strings"
	"sync"

	kgsearch "google.golang.org/api/kgsearch/v1"
	"google.golang.org/api/option"
)

type EntityElement struct {
	ResultScore float64 `json:"resultScore"`
	Result      struct {
		ID    string `json:"@id"`
		Image struct {
			URL        string `json:"url"`
			ContentURL string `json:"contentUrl"`
		} `json:"image,omitempty"`
		DetailedDescription struct {
			URL         string `json:"url"`
			License     string `json:"license"`
			ArticleBody string `json:"articleBody"`
		} `json:"detailedDescription"`
		Type        []string `json:"@type"`
		Description string   `json:"description"`
		Name        string   `json:"name"`
		URL         string   `json:"url,omitempty"`
	} `json:"result"`
	Type string `json:"@type"`
}

func Entity(ctx context.Context, query string, etype string) (string, error) {
	googlecfg := config.AppConfig.GoogelConfig
	var element []EntityElement
	var result string
	kgsearchService, err := kgsearch.NewService(ctx, option.WithAPIKey(googlecfg.ApiKey))
	if err != nil {
		logger.Errorf("%s", err)
		return "", err
	}
	resp, err := kgsearchService.Entities.Search().Languages("zh").Types(etype).Query(query).Limit(5).Do()

	data, err := json.Marshal(resp.ItemListElement)
	if err != nil {
		logger.Errorf("反序列化%v", err)
	}
	err = json.Unmarshal([]byte(data), &element)
	if err != nil {
		logger.Errorf("序列化%v", err)
	}

	// for _, v := range element {
	// 	result += fmt.Sprintf("{\"Entity\":\"%s\",\n\"Description\"%s\",\n\"Describe\":\"%s\",\n\"Weblink\":\"%s\"\n}", v.Result.Name, v.Result.Description, v.Result.DetailedDescription.ArticleBody, v.Result.DetailedDescription.URL)
	// }
	if len(element) == 0 {
		logger.Debugf("entity无结果")
		if etype != "Thing" {
			content, err := Entity(ctx, query, "Thing")
			return content, err
		} else {
			return "", nil
		}
	} else {
		wg := sync.WaitGroup{}
		var lock sync.Mutex
		for _, v := range element {
			wg.Add(1)
			go func(v EntityElement) {
				defer wg.Done()
				logger.Debug("[EntitySearch]", logger.Pair("关键词", query), logger.Pair("实体名称", v.Result.Name), logger.Pair("实体链接", v.Result.DetailedDescription.URL))
				content := crawlPage(v.Result.DetailedDescription.URL)
				var summary string
				if content != "" {
					summary = SummaryContent(query, "Entity:"+v.Result.Name+"\n"+"Describe"+v.Result.DetailedDescription.ArticleBody+"\n"+"Content:"+content, 0)
				}
				lock.Lock()
				if !strings.Contains(summary, "[NO_CONTENT]") {
					result += "Entity:\n" + v.Result.Name + "\n" + "Description:\n" + v.Result.Description + "\n" + v.Result.DetailedDescription.ArticleBody + "\n" + "Content:\n" + summary + "\n" + "Web Link:\n" + v.Result.DetailedDescription.URL + "\n"
				}
				lock.Unlock()
			}(v)

		}
		wg.Wait()
		return result, nil
	}
}
