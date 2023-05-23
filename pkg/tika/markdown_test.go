/*
 * @Author: cloudyi.li
 * @Date: 2023-05-22 13:41:03
 * @LastEditTime: 2023-05-22 16:15:24
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/tika/markdown_test.go
 */
package tika

import (
	"fmt"
	"testing"
)

func TestProcessMarkDown(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test1",
			args: args{filename: "../../uploadfile/科学技术普及法(2002-06-29).md"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ProcessMarkDown(tt.args.filename)
			if err != nil {
				fmt.Print(err)
			}
			fmt.Println(got)
		})
	}
}
