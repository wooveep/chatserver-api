/*
 * @Author: cloudyi.li
 * @Date: 2023-03-29 10:33:02
 * @LastEditTime: 2023-05-10 12:29:57
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/utils/security/encrypt_test.go
 */
package security

import (
	"testing"
)

func TestPassword(t *testing.T) {
	type args struct {
		Plain string
		Key   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "Test1",
			args: args{
				Plain: "abcde1234",
				Key:   "ABCDABCDABCDABCD",
			},
			want: "abcde1234",
		},
		{
			name: "Test2",
			args: args{
				Plain: "abcde1234",
				Key:   "ABCDABCDABCDABCD",
			},
			want: "abcde1234",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PasswordDecryption(PasswordEncrypt(tt.args.Plain, tt.args.Key), tt.args.Key); got != tt.want {
				t.Errorf("PasswordDecryption() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPasswordDecryption(t *testing.T) {
	type args struct {
		Cipher string
		Key    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test2",
			args: args{
				Cipher: "3ZZ/mQNpzvK1o3zxu0C1oQ==",
				Key:    "ABCDABCDABCDABCD",
			},
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PasswordDecryption(tt.args.Cipher, tt.args.Key); got != tt.want {
				t.Errorf("PasswordDecryption() = %v, want %v", got, tt.want)
			}
		})
	}
}
