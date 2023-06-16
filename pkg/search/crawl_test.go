/*
 * @Author: cloudyi.li
 * @Date: 2023-06-07 09:50:35
 * @LastEditTime: 2023-06-15 22:01:30
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/search/crawl_test.go
 */
package search

import (
	"fmt"
	"testing"
)

func Test_crawlpage(t *testing.T) {
	type args struct {
		u string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			args: args{u: "https://wiki.mbalib.com/wiki/IBM%E5%85%AC%E5%8F%B8"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Printf(crawlPage(tt.args.u))
		})
	}
}
