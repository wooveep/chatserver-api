/*
 * @Author: cloudyi.li
 * @Date: 2023-06-16 21:16:47
 * @LastEditTime: 2023-06-26 17:08:40
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/chatfunc/chatfunc_test.go
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
		// {name: "tset2",
		// 	args: args{
		// 		name:      "GetWeather",
		// 		arguments: "{\n  \"location\": \"南京\"\n}",
		// 	}},
		{
			name: "tset2",
			args: args{
				name:      "EntitySearch",
				arguments: "{\n  \"query\": \"中国国务院总理\",\n  \"etype\": \"Person\"\n}",
			},
		},
		// {
		// 	name: "test5",
		// 	args: args{
		// 		name:      "GoogleSearch",
		// 		arguments: "{\n  \"query\": \"GPT-4由8个MoE模型组成\"\n}",
		// 	},
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println(ChatFuncProcess(context.Background(), tt.args.name, tt.args.arguments))
		})
	}
}

func Test_queryExtra(t *testing.T) {
	c := config.Load("../../configs/config.yml")
	logger.InitLogger(&c.LogConfig, c.AppName)
	cache.InitRedis(c.RedisConfig)
	type args struct {
		message string
	}
	tests := []struct {
		name string
		args args
		// want string
	}{
		// TODO: Add test cases.
		{
			name: "test2",
			args: args{
				message: "上海明天会下雨嘛",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Print(queryExtra(tt.args.message))
		})
	}
}

func TestCustomFuncExtension(t *testing.T) {
	c := config.Load("../../configs/config.yml")
	logger.InitLogger(&c.LogConfig, c.AppName)
	cache.InitRedis(c.RedisConfig)

	type args struct {
		ctx           context.Context
		queryoriginal string
	}
	tests := []struct {
		name string
		args args
		// wantContent string
		// wantErr     bool
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			args: args{
				ctx:           context.Background(),
				queryoriginal: "上海明天会下雨嘛",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotContent, _ := CustomFuncExtension(tt.args.ctx, tt.args.queryoriginal)
			fmt.Print(gotContent)
		})
	}
}
