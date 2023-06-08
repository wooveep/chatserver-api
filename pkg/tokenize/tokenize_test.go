/*
 * @Author: cloudyi.li
 * @Date: 2023-05-08 14:04:14
 * @LastEditTime: 2023-06-08 07:01:02
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/tokenize/tokenize_test.go
 */
package tokenize

import (
	"fmt"
	"testing"
)

func TestGetKeyword(t *testing.T) {
	tk := NewTokenizer("../../dict")
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
				s: "南京第一医院怎么样",
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

func Test_tokenizer_GetSearch(t *testing.T) {
	tk := NewTokenizer("../../dict")
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		// wantSearch string
	}{
		// TODO: Add test cases.
		{
			name: "tset1",
			args: args{
				s: "南京有哪些好玩的",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := tk.GetSearch(tt.args.s)
			fmt.Println(a)
		})
	}
}
