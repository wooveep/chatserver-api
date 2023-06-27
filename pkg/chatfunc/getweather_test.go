/*
 * @Author: cloudyi.li
 * @Date: 2023-06-15 09:23:25
 * @LastEditTime: 2023-06-15 20:54:27
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/chatfunc/getweather_test.go
 */
package chatfunc

import (
	"chatserver-api/pkg/cache"
	"chatserver-api/pkg/config"
	"chatserver-api/pkg/logger"
	"context"
	"fmt"
	"testing"
)

func TestGetWeather(t *testing.T) {
	c := config.Load("../../configs/config.yml")
	logger.InitLogger(&c.LogConfig, c.AppName)
	cache.InitRedis(c.RedisConfig)

	type args struct {
		address string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			args: args{
				address: "南京",
			},
		},
		{
			name: "test2",
			args: args{
				address: "上海",
			},
		},
		{
			name: "test3",
			args: args{
				address: "平潭",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Print(GetWeather(context.Background(), tt.args.address))
		})
	}
}
