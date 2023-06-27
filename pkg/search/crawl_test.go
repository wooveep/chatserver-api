/*
 * @Author: cloudyi.li
 * @Date: 2023-06-07 09:50:35
 * @LastEditTime: 2023-06-26 11:51:20
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
			args: args{u: "https://www.kanzhun.com/baike_salary/1X170w~~/"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Printf(crawlPage(tt.args.u))
		})
	}
}
