/*
 * @Author: cloudyi.li
 * @Date: 2023-06-07 09:50:35
 * @LastEditTime: 2023-06-07 17:04:42
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/search/crawl_test.go
 */
package search

import "testing"

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
			args: args{u: "http://www.mafengwo.cn/gonglve/ziyouxing/2606.html"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			crawlPage(tt.args.u)
		})
	}
}
