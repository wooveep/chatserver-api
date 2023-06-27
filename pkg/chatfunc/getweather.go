/*
 * @Author: cloudyi.li
 * @Date: 2023-06-15 09:23:25
 * @LastEditTime: 2023-06-21 16:24:41
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/chatfunc/getweather.go
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
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var FuncGetWeather = openai.FunctionDefine{
	Name:        "GetWeather",
	Description: "Get the current weather",
	Parameters: &openai.FunctionParams{
		Type: openai.JSONSchemaTypeObject,
		Properties: map[string]*openai.JSONSchemaDefine{
			"location": {
				Type:        openai.JSONSchemaTypeString,
				Description: "The city and state, e.g. San Francisco, CA",
			},
		},
		Required: []string{"location"},
	},
}

func GetWeather(ctx context.Context, address string) string {
	// lat, lng, _ := search.GeoSearch(address)
	// w, err := owm.NewForecast("5", "C", "zh_cn", "") // valid options for first parameter are "5" and "16"
	// if err != nil {
	// 	logger.Error(err.Error())
	// }

	// w.DailyByCoordinates(
	// 	&owm.Coordinates{
	// 		Longitude: lng,
	// 		Latitude:  lat,
	// 	},
	// 	5,
	// )
	// fmt.Println(w.ForecastWeatherJson.(*owm.Forecast5WeatherData))
	rc := cache.GetRedisClient()
	lat, lng, fmtadd := search.GeoSearch(address)
	var content string
	content, err := rc.Get(ctx, consts.SearchCachePrefix+security.Md5(fmtadd+"weather")).Result()
	if err != nil {
		if err != redis.Nil {
			logger.Errorf("Redis异常%v", err)
		}
		logger.Debugf("weather天气无缓存")
		wcli := search.NewWeatherClient("3d042ce3865a4fe7842ce3865a0fe725")
		url := wcli.Weathermakeapiurl(lat, lng, "m")
		content, err = wcli.WeatherdoGetForecast5(url)
		if err != nil {
			logger.Errorf("获取天气API异常：%v", err)
		} else {
			err = rc.Set(ctx, consts.SearchCachePrefix+security.Md5(fmtadd+"weather"), content, 6*time.Hour).Err()
			if err != nil {
				logger.Errorf("Redis异常%v", err)
			}
		}
	}
	result := fmt.Sprintf("Address: %s\nCurrent_Date_Time:%s\n%s\nWebLink: https://weather.com/zh-CN/weather/today/l/%.2f,%.2f?par=google", fmtadd, time.Now().Format(consts.TimeLayout), content, lat, lng)
	return result
}
