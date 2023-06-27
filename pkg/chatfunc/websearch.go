/*
 * @Author: cloudyi.li
 * @Date: 2023-06-21 16:28:30
 * @LastEditTime: 2023-06-27 13:47:43
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/chatfunc/websearch.go
 */
package chatfunc

import (
	"chatserver-api/internal/consts"
	"chatserver-api/pkg/cache"
	"chatserver-api/pkg/logger"
	"chatserver-api/pkg/openai"
	"chatserver-api/pkg/search"
	"chatserver-api/utils/security"
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

var FuncGoogleSearch = openai.FunctionDefine{
	Name:        "GoogleSearch",
	Description: "When the content sent by the user contains some newly occurred events, news or other things that you don't know, use this function to call the Google search engine.",
	Parameters: &openai.FunctionParams{
		Type: openai.JSONSchemaTypeObject,
		Properties: map[string]*openai.JSONSchemaDefine{
			"query": {
				Type:        openai.JSONSchemaTypeString,
				Description: "Extracted query statements from user questions that are applicable for search engines.",
			},
			"classify": {
				Type: openai.JSONSchemaTypeString,
				Enum: []string{
					"News",
					"Custom",
				},
				Description: "Determine what category the content to be searched belongs to.",
			},
		},
		Required: []string{"query", "classify"},
	},
}

func GoogleSearch(ctx context.Context, query string, classify string) string {
	var content string
	rc := cache.GetRedisClient()
	content, err := rc.Get(ctx, consts.SearchCachePrefix+security.Md5(query)).Result()
	if err != nil {
		if err != redis.Nil {
			logger.Errorf("Redis异常%v", err)
		}
		logger.Debugf("EntitySearch无缓存")
		// wcli := search.NewWeatherClient("3d042ce3865a4fe7842ce3865a0fe725")
		// url := wcli.Weathermakeapiurl(lat, lng, "m")
		// content, err = wcli.WeatherdoGetForecast5(url)
		content, err = search.CustomSearch(ctx, query, classify)
		if err != nil {
			logger.Errorf("获取EntitySearch无缓存异常：%v", err)
		} else {
			err = rc.Set(ctx, consts.SearchCachePrefix+security.Md5(query), content, 72*time.Hour).Err()
			if err != nil {
				logger.Errorf("Redis异常%v", err)
			}
		}
	}
	return content
}
