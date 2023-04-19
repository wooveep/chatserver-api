/*
 * @Author: cloudyi.li
 * @Date: 2023-04-05 15:37:14
 * @LastEditTime: 2023-04-19 14:41:47
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/dao/query/chat_test.go
 */
package query

import (
	"chatserver-api/internal/model"
	"chatserver-api/pkg/config"
	"chatserver-api/pkg/db"
	"chatserver-api/pkg/logger"
	"context"
	"fmt"
	"testing"
)

func Test_chatDao_ChatRegenRecordGet(t *testing.T) {
	c := config.Load("../../../configs/config.yml")
	logger.InitLogger(&c.LogConfig, c.AppName)
	// confi := config.DBConfig{
	// 	Dbname:   "whatserver",
	// 	Host:     "192.168.10.251",
	// 	Port:     "5432",
	// 	Username: "whatserver",
	// 	Password: "whatserver123",
	// }
	ds := db.NewDefaultPostGre(c.DBConfig)
	defer ds.Close()
	type args struct {
		ctx    context.Context
		chatId int64
		msgid  int64
		memory int16
	}
	tests := []struct {
		name    string
		args    args
		wantr   []model.RecordOne
		wantid  int64
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				ctx:    context.Background(),
				chatId: 1647868345921310720,
				msgid:  1648535109717987328,
				memory: 5,
			},
			wantr:   []model.RecordOne{},
			wantid:  1648525207868018688,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cd := &chatDao{
				ds: ds,
			}
			got, gotid, err := cd.ChatRegenRecordGet(tt.args.ctx, tt.args.chatId, tt.args.msgid, tt.args.memory)
			if (err != nil) != tt.wantErr {
				t.Errorf("chatDao.ChatRegenRecordGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for i, v := range got {
				fmt.Printf("id:%d,msgid:%d,msg:%s\n", i, v.Id, v.Message)
			}
			fmt.Printf("id:%d\n", gotid)
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("chatDao.ChatRegenRecordGet() = %v, want %v", got, tt.want)
			// }
		})
	}
}
