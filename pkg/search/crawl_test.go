/*
 * @Author: cloudyi.li
 * @Date: 2023-06-07 09:50:35
 * @LastEditTime: 2023-06-21 21:49:15
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/search/crawl_test.go
 */
package search

import (
	"chatserver-api/pkg/config"
	"chatserver-api/pkg/logger"
	"fmt"
	"testing"
)

func Test_crawlpage(t *testing.T) {
	c := config.Load("../../configs/config.yml")
	logger.InitLogger(&c.LogConfig, c.AppName)
	type args struct {
		u string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "test2",
			args: args{u: "https://www.pconline.com.cn/focus/1627/16275165.html"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Printf(crawlPage(tt.args.u))
		})
	}
}
