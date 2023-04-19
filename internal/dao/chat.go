/*
 * @Author: cloudyi.li
 * @Date: 2023-04-05 15:36:51
 * @LastEditTime: 2023-04-19 15:23:05
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
	ChatRecordUpdate(ctx context.Context, record *entity.Record) error
	ChatRecordGet(ctx context.Context, chatId int64, memory int16) ([]model.RecordOne, error)
	ChatRegenRecordGet(ctx context.Context, chatId, msgid int64, memory int16) ([]model.RecordOne, int64, error)
	ChatGetList(ctx context.Context, userId int64) ([]model.ChatOne, error)
	ChatDetailGet(ctx context.Context, userId, chatId int64) (model.ChatDetail, error)
	ChatDeleteOne(ctx context.Context, userId, chatId int64) error
	ChatDeleteAll(ctx context.Context, userId int64) error
	ChatUserVerify(ctx context.Context, userId, chatId int64) (int64, error)
	ChatCostUpdate(ctx context.Context, userId int64, balance float64) error
	ChatBalanceGet(ctx context.Context, userId int64) (model.UserBalance, error)
	ChatRecordVerify(ctx context.Context, recordid int64) (int64, error)
}
