/*
 * @Author: cloudyi.li
 * @Date: 2023-04-12 05:39:35
 * @LastEditTime: 2023-04-12 05:58:35
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/utils/tools/tokenzi_test.go
 */
package tools

import "testing"

func Test_tokenzi(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "test1",
			args: args{"aaa"},
			want: 4,
		},
		{
			name: "test2",
			args: args{"你好"},
			want: 8,
		},
		{
			name: "test3",
			args: args{"Hello!"},
			want: 8,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tokenzi(tt.args.str); got != tt.want {
				t.Errorf("tokenzi() = %v, want %v", got, tt.want)
			}
		})
	}
}
