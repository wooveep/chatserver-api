/*
 * @Author: cloudyi.li
 * @Date: 2023-06-16 21:16:47
 * @LastEditTime: 2023-06-16 21:29:39
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/chatfunc/chatfunc_test.go
 */
package chatfunc

import (
	"chatserver-api/pkg/cache"
	"chatserver-api/pkg/config"
	"chatserver-api/pkg/logger"
	"fmt"
	"testing"
)

func TestChatFuncProcess(t *testing.T) {
	c := config.Load("../../configs/config.yml")
	logger.InitLogger(&c.LogConfig, c.AppName)
	cache.InitRedis(c.RedisConfig)

	type args struct {
		name      string
		arguments string
	}
	tests := []struct {
		name string
		args args
		// wantContent string
	}{
		// TODO: Add test cases.
		{name: "tset2",
			args: args{
				name:      "GetWeather",
				arguments: "{\n  \"location\": \"南京\"\n}",
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println(ChatFuncProcess(tt.args.name, tt.args.arguments))
		})
	}
}
