/*
 * @Author: cloudyi.li
 * @Date: 2023-05-31 17:32:37
 * @LastEditTime: 2023-05-31 18:38:01
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/mail/verifymail_test.go
 */
package mail

import (
	"fmt"
	"testing"
)

func Test_verifierEmail(t *testing.T) {
	type args struct {
		email string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test1",
			args: args{email: "11@transwarp.io"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ver := NewVerifier()
			err := ver.VerifierEmail(tt.args.email)
			if err != nil {
				fmt.Println(err)
			}
		})
	}
}
