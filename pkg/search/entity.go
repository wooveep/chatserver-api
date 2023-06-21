/*
 * @Author: cloudyi.li
 * @Date: 2023-06-17 21:58:47
 * @LastEditTime: 2023-06-20 14:39:54
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/search/entity.go
 */
package search

import (
	"chatserver-api/pkg/config"
	"chatserver-api/pkg/logger"
	"context"
	"encoding/json"
	"fmt"

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

//	{
//
// "name": "EntitySearch",
// "description": "When the content sent by the user contains special entities, call this function to perform a knowledge graph search.",
//
//	"parameters": {
//		"type": "object",
//		"properties": {
//			"query": {
//				"type": "string",
//				"description": "the entity name"
//			},
//			"type": {
//				"type": "string",
//				"enum": [
//					"Action",
//					"BioChemEntity",
//					"CreativeWork",
//					"Event",
//					"Intangible",
//					"MedicalEntity",
//					"Organization",
//					"Person",
//					"Place",
//					"Product"
//				],
//				"description": "the entity type"
//			}
//		},
//		"required": [
//			"query",
//			"type"
//		]
//	}
//
//	},
func Entity(ctx context.Context, query string, etype string) (string, error) {
	googlecfg := config.AppConfig.GoogelConfig
	var element []EntityElement
	var result string
	kgsearchService, err := kgsearch.NewService(ctx, option.WithAPIKey(googlecfg.ApiKey))
	// client := &http.Client{Transport: &transport.APIKey{Key: googlecfg.ApiKey}}
	// kvc, err := kgsearch.New(client)
	if err != nil {
		logger.Errorf("%s", err)
		return "", err
	}
	resp, err := kgsearchService.Entities.Search().Languages("zh").Types(etype).Query(query).Limit(5).Do()
	// resp, err := kvc.Entities.Search().Languages("zh").Query(query).Limit(2).Do()
	// if err != nil {
	// 	logger.Errorf("%s", err)
	// 	return "", err
	// }
	data, err := json.Marshal(resp.ItemListElement)
	if err != nil {
		logger.Errorf("反序列化%v", err)
	}
	err = json.Unmarshal([]byte(data), &element)
	if err != nil {
		logger.Errorf("序列化%v", err)
	}
	// for _, v := range element {
	// 	// a := element.(EntityElement)
	// 	fmt.Println(v.Result.Name)
	// 	fmt.Println(v.Result.DetailedDescription.URL)
	// 	fmt.Println(v.Result.DetailedDescription.ArticleBody)
	// 	fmt.Println(v.Result.Image.ContentURL)

	// }
	for _, v := range element {
		result += fmt.Sprintf("{\"Entity\":\"%s\",\n\"Description\"%s\",\n\"Describe\":\"%s\",\n\"Weblink\":\"%s\"\n}", v.Result.Name, v.Result.Description, v.Result.DetailedDescription.ArticleBody, v.Result.DetailedDescription.URL)
	}
	if len(element) == 0 {
		logger.Debugf("entity无结果")
		if etype != "Thing" {
			content, err := Entity(ctx, query, "Thing")
			return content, err
		} else {
			return "", nil
		}
	} else {
		text := crawlPage(element[0].Result.DetailedDescription.URL)
		result += fmt.Sprintf("Text:%s", text)
		return result, nil
	}
}
