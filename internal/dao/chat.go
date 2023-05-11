/*
 * @Author: cloudyi.li
 * @Date: 2023-04-05 15:36:51
 * @LastEditTime: 2023-05-11 09:20:49
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/dao/chat.go
 */
package dao

import (
	"chatserver-api/internal/model"
	"chatserver-api/internal/model/entity"
	"chatserver-api/pkg/pgvector"
	"context"
)

type ChatDao interface {
	ChatCreateNew(ctx context.Context, chat *entity.Chat) error
	ChatUpdate(ctx context.Context, chat *entity.Chat) error
	ChatRecordSave(ctx context.Context, record *entity.Record) error
	ChatRecordClear(ctx context.Context, chatId int64) error
	DocEmbeddingSave(ctx context.Context, docs *entity.Documents) error
	ChatRecordUpdate(ctx context.Context, record *entity.Record) error
	ChatRecordGet(ctx context.Context, chatId int64, memory int16) ([]model.RecordOne, error)
	ChatRegenRecordGet(ctx context.Context, chatId, msgid int64, memory int16) ([]model.RecordOne, int64, error)
	ChatListGet(ctx context.Context, userId int64) ([]model.ChatOne, error)
	ChatDetailGet(ctx context.Context, userId, chatId int64) (model.ChatDetail, error)
	ChatDeleteOne(ctx context.Context, userId, chatId int64) error
	ChatDeleteAll(ctx context.Context, userId int64) error
	ChatUserVerify(ctx context.Context, userId, chatId int64) (int64, error)
	ChatCostUpdate(ctx context.Context, userId int64, balance float64) error
	ChatBalanceGet(ctx context.Context, userId int64) (model.UserBalance, error)
	ChatRecordVerify(ctx context.Context, recordid int64) (int64, error)
	ChatEmbeddingCompare(ctx context.Context, question pgvector.Vector, classify string) ([]model.DocsCompare, error)
}
