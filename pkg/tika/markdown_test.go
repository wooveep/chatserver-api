/*
 * @Author: cloudyi.li
 * @Date: 2023-05-22 13:41:03
 * @LastEditTime: 2023-05-29 09:39:33
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
			args: args{filename: "../../uploadfile/保守国家秘密法(2010-04-29).md"},
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
