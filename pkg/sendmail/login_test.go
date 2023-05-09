/*
 * @Author: cloudyi.li
 * @Date: 2023-05-09 14:27:26
 * @LastEditTime: 2023-05-09 16:30:40
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/sendmail/login_test.go
 */
package sendmail

import (
	"chatserver-api/pkg/config"
	"chatserver-api/pkg/logger"
	"testing"
)

func TestLoginMail(t *testing.T) {
	c := config.Load("../../configs/config.yml")
	logger.InitLogger(&c.LogConfig, c.AppName)
	type args struct {
		mailTo         []string
		activationCode string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				mailTo:         []string{"wooveep@outlook.com"},
				activationCode: "asdfasdfasdf",
			},
			wantErr: false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := LoginMail(tt.args.mailTo, tt.args.activationCode); (err != nil) != tt.wantErr {
				t.Errorf("LoginMail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
