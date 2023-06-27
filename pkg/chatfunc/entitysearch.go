/*
 * @Author: cloudyi.li
 * @Date: 2023-06-18 16:56:17
 * @LastEditTime: 2023-06-19 12:44:16
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/chatfunc/entitysearch.go
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

var FuncEntitySearch = openai.FunctionDefine{
	Name:        "EntitySearch",
	Description: "When the content sent by the user contains special entities, call this function to perform a knowledge graph search.",
	Parameters: &openai.FunctionParams{
		Type: openai.JSONSchemaTypeObject,
		Properties: map[string]*openai.JSONSchemaDefine{
			"query": {
				Type:        openai.JSONSchemaTypeString,
				Description: "the entity name",
			},
			"etype": {
				Type: openai.JSONSchemaTypeString,
				Enum: []string{
					"Action",
					"BioChemEntity",
					"CreativeWork",
					"Event",
					"Intangible",
					"MedicalEntity",
					"Organization",
					"Person",
					"Place",
					"Product"},
				Description: "the entity type",
			},
		},
		Required: []string{"query", "etype"},
	},
}

func EntitySearch(ctx context.Context, query string, etype string) string {
	var content string
	rc := cache.GetRedisClient()
	content, err := rc.Get(ctx, consts.SearchCachePrefix+security.Md5(query+etype)).Result()
	if err != nil {
		if err != redis.Nil {
			logger.Errorf("Redis异常%v", err)
		}
		logger.Debugf("EntitySearch无缓存")
		// wcli := search.NewWeatherClient("3d042ce3865a4fe7842ce3865a0fe725")
		// url := wcli.Weathermakeapiurl(lat, lng, "m")
		// content, err = wcli.WeatherdoGetForecast5(url)
		content, err = search.Entity(ctx, query, etype)
		if err != nil {
			logger.Errorf("获取EntitySearch无缓存异常：%v", err)
		} else {
			err = rc.Set(ctx, consts.SearchCachePrefix+security.Md5(query+etype), content, 72*time.Hour).Err()
			if err != nil {
				logger.Errorf("Redis异常%v", err)
			}
		}
	}
	return content
}
