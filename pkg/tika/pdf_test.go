/*
 * @Author: cloudyi.li
 * @Date: 2023-05-01 12:55:30
 * @LastEditTime: 2023-05-31 14:44:33
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/tika/pdf_test.go
 */
package tika

import (
	"fmt"
	"testing"
)

func Test_readPdf2(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "test1",
			args:    args{"hyperbase.pdf"},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readPdf2(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("readPdf2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if got != tt.want {
			// 	t.Errorf("readPdf2() = %v, want %v", got, tt.want)
			// }
			fmt.Print(got)
		})
	}
}

func Test_readPdf(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "test1",
			args:    args{"hyperbase.pdf"},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readPdf(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("readPdf() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Print(got)
		})
	}
}

// func Test_readPd3f(t *testing.T) {
// 	type args struct {
// 		path string
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 	}{
// 		// TODO: Add test cases.
// 		{
// 			name: "test1",
// 			args: args{"/Users/cloudyi/goProject/chatserver-api/uploadfile/创远信科_SK-6000详细技术手册.pdf"},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			a, err := ReadPd3f(tt.args.path)
// 			if err != nil {
// 				fmt.Printf(err.Error())
// 			}
// 			for _, v := range a {
// 				fmt.Println(v)
// 			}
// 		})
// 	}
// }
