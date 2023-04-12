/*
 * @Author: cloudyi.li
 * @Date: 2023-04-05 15:36:51
 * @LastEditTime: 2023-04-12 20:12:28
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/dao/chat.go
 */
package dao

import (
	"chatserver-api/internal/model"
	"chatserver-api/internal/model/entity"
	"context"
)

type ChatDao interface {
	ChatCreateNew(ctx context.Context, chat *entity.Chat) error
	ChatRecordSave(ctx context.Context, record *entity.Record) error
	ChatRecordGet(ctx context.Context, userid int64, chatid int64, memory int16) ([]model.RecordOne, error)
	// ChatRecordGetList(ctx context.Context, userid int64, chatid int64, memory int16) ([]entity.ChatRecord, error)
	ChatGetList(ctx context.Context, userid int64) ([]model.ChatOne, error)
	ChatDetailGet(ctx context.Context, userid, chatid int64) (model.ChatDetail, error)
}
