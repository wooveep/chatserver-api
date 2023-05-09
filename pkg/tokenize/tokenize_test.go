/*
 * @Author: cloudyi.li
 * @Date: 2023-05-08 14:04:14
 * @LastEditTime: 2023-05-09 13:41:49
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/tokenize/tokenize_test.go
 */
package tokenize

import (
	"fmt"
	"testing"
)

func TestGetKeyword(t *testing.T) {
	tk := NewTokenizer()
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		// wantWords []string
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			args: args{
				s: "hyperbase是什么 怎么连接，如何批量提交",
			},
		},
		{
			name: "test2",
			args: args{
				s: "hyperbase如何调优 Rowkey如何设计",
			},
		},
		{
			name: "test3",
			args: args{
				s: "inceptor 如何调优 如何进行CPU和内存的设置 应该选择什么队列 ",
			},
		},
		{
			name: "test4",
			args: args{
				s: "中新赛克移动互联网设备如何配置复合规则 	可以举例说明么  V6规则如何配置",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotWords := tk.GetKeyword(tt.args.s)
			fmt.Println(gotWords)
		})
	}
}
