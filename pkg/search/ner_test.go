/*
 * @Author: cloudyi.li
 * @Date: 2023-06-06 11:23:44
 * @LastEditTime: 2023-06-26 16:11:18
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/search/ner_test.go
 */
package search

import (
	"chatserver-api/pkg/config"
	"chatserver-api/pkg/logger"
	"reflect"
	"testing"
)

func Test_nerDetec(t *testing.T) {
	c := config.Load("../../configs/config.yml")
	logger.InitLogger(&c.LogConfig, c.AppName)
	type args struct {
		query string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			args: args{
				query: "这个问题应该如何处理呢",
			},
			want: 0,
		},
		{
			name: "test2",
			args: args{
				query: "今天有什么新闻？",
			},
			want: 2,
		},
		{
			name: "test3",
			args: args{
				query: "这个问题处理一下",
			},
			want: 0,
		},
		{
			name: "test4",
			args: args{
				query: "这个事情的导火索是什么呢",
			},
			want: 0,
		},
		{
			name: "test5",
			args: args{
				query: "南京第一医院怎么样",
			},
			want: 1,
		},
		{
			name: "test6",
			args: args{
				query: "南京今天天气怎样",
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := NerDetec(tt.args.query)
			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("nerDetec() = %v, want %v", got, tt.want)
			}
		})
	}
}
