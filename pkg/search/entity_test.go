/*
 * @Author: cloudyi.li
 * @Date: 2023-06-17 21:58:47
 * @LastEditTime: 2023-06-20 22:31:48
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/search/entity_test.go
 */
package search

import (
	"chatserver-api/pkg/cache"
	"chatserver-api/pkg/config"
	"chatserver-api/pkg/logger"
	"context"
	"fmt"
	"testing"
)

func TestEntitySearch(t *testing.T) {
	c := config.Load("../../configs/config.yml")
	logger.InitLogger(&c.LogConfig, c.AppName)
	cache.InitRedis(c.RedisConfig)
	type args struct {
		// ctx   context.Context
		query string
		etype string
	}
	tests := []struct {
		name string
		args args
		// want    string
		// wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			args: args{
				query: "南京邮电大学",
				etype: "Organization",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println(Entity(context.Background(), tt.args.query, tt.args.etype))
		})
	}
}
