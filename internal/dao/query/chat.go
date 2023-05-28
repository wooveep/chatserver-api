/*
 * @Author: cloudyi.li
 * @Date: 2023-04-05 15:37:14
 * @LastEditTime: 2023-05-28 17:46:45
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/dao/query/chat.go
 */
package query

import (
	"chatserver-api/internal/dao"
	"chatserver-api/internal/model"
	"chatserver-api/internal/model/entity"
	"chatserver-api/pkg/db"
	"chatserver-api/pkg/pgvector"
	"context"

	"gorm.io/gorm/clause"
)

var _ dao.ChatDao = (*chatDao)(nil)

type chatDao struct {
	ds db.IDataSource
}

func NewChatDao(_ds db.IDataSource) *chatDao {
	return &chatDao{
		ds: _ds,
	}
}

func (cd *chatDao) ChatCreateNew(ctx context.Context, chat *entity.Chat) error {
	return cd.ds.Master().Create(chat).Error
}

func (cd *chatDao) ChatUpdate(ctx context.Context, chat *entity.Chat) error {
	return cd.ds.Master().Updates(chat).Error
}

func (cd *chatDao) ChatRecordSave(ctx context.Context, record *entity.Record) error {
	return cd.ds.Master().Create(record).Error
}

func (cd *chatDao) ChatRecordClear(ctx context.Context, chatId int64) error {
	return cd.ds.Master().Where("chat_id = ?", chatId).Delete(&entity.Record{}).Error
}

func (cd *chatDao) ChatRecordIdGet(ctx context.Context, chatId int64) ([]int64, error) {
	var recordlist []int64
	err := cd.ds.Master().Where("chat_id = ?", chatId).Model(&entity.Record{}).Select("id").Find(&recordlist).Error
	return recordlist, err
}

func (cd *chatDao) DocEmbeddingSave(ctx context.Context, docs *entity.Documents) error {
	return cd.ds.Master().Create(docs).Error
}

func (cd *chatDao) ChatRecordVerify(ctx context.Context, recordid int64) (int64, error) {
	var count int64
	err := cd.ds.Master().Model(&entity.Record{}).Where("id = ? ", recordid).Count(&count).Error
	return count, err
}

func (cd *chatDao) ChatRecordUpdate(ctx context.Context, record *entity.Record) error {
	return cd.ds.Master().Updates(record).Error
}

func (cd *chatDao) ChatDeleteOne(ctx context.Context, userId, chatId int64) error {
	return cd.ds.Master().Select("Records").Delete(&entity.Chat{Id: chatId}).Error
}

func (cd *chatDao) ChatDeleteAll(ctx context.Context, userId int64) error {
	var chatIds []int64
	err := cd.ds.Master().Model(&entity.Chat{}).Where("user_id = ?", userId).Select("id").Find(&chatIds).Error
	if err != nil {
		return err
	}
	for _, v := range chatIds {
		err := cd.ds.Master().Select("Records").Delete(&entity.Chat{Id: v}).Error
		if err != nil {
			return err
		}
	}
	return err
}

func (cd *chatDao) ChatDetailGet(ctx context.Context, userId, chatId int64) (model.ChatDetail, error) {
	var detail model.ChatDetail
	err := cd.ds.Master().InnerJoins("Chats", cd.ds.Master().Where(&entity.Chat{Id: chatId, UserId: userId})).Model(&entity.Preset{}).Find(&detail).Error
	return detail, err
}

func (cd *chatDao) ChatRecordGet(ctx context.Context, chatId int64, memory int16) ([]model.RecordOne, error) {
	var recordlist []model.RecordOne
	subQuery1 := cd.ds.Master().Model(&entity.Record{}).Where("chat_id = ?", chatId).Order("id desc").Limit(int(memory)).Select("*")
	err := cd.ds.Master().Table("(?) as a ", subQuery1).Order("id").Find(&recordlist).Error
	return recordlist, err
}

func (cd *chatDao) ChatRegenRecordGet(ctx context.Context, chatId, msgid int64, memory int16) ([]model.RecordOne, int64, error) {
	var recordlist []model.RecordOne
	var answerid int64
	err := cd.ds.Master().Model(&entity.Record{}).Where("chat_id = ?", chatId).Where("id > ? ", msgid).Order("id").Limit(1).Select("id").Find(&answerid).Error
	subQuery2 := cd.ds.Master().Model(&entity.Record{}).Where("chat_id = ?", chatId).Where("id <= ? ", msgid).Order("id desc").Limit(int(memory) + 1).Select("*")
	err = cd.ds.Master().Table("(?) as a ", subQuery2).Order("id").Find(&recordlist).Error
	return recordlist, answerid, err
}

func (cd *chatDao) ChatListGet(ctx context.Context, userId int64) ([]model.ChatOne, error) {
	var chatlist []model.ChatOne
	err := cd.ds.Master().Model(&entity.Chat{}).Where("user_id = ? ", userId).Order("id").Find(&chatlist).Error
	return chatlist, err
}

func (cd *chatDao) ChatUserVerify(ctx context.Context, userId, chatId int64) (int64, error) {
	var count int64
	err := cd.ds.Master().Model(&entity.Chat{}).Where("id = ? ", chatId).Where("user_id = ?", userId).Count(&count).Error
	return count, err
}

func (cd *chatDao) ChatCostUpdate(ctx context.Context, userId int64, balance float64) error {
	return cd.ds.Master().Model(&entity.User{}).Where("id  = ? ", userId).UpdateColumn("balance", balance).Error
}

func (cd *chatDao) ChatBalanceGet(ctx context.Context, userId int64) (model.UserBalance, error) {
	var userbalance model.UserBalance
	err := cd.ds.Master().Model(&entity.User{}).Where("id  = ? ", userId).Find(&userbalance).Error
	return userbalance, err
}

func (cd *chatDao) ChatEmbeddingCompare(ctx context.Context, question pgvector.Vector, classify string) ([]model.DocsCompare, error) {
	var docsbody []model.DocsCompare
	err := cd.ds.Master().Model(&entity.Documents{}).Clauses(clause.OrderBy{Expression: clause.Expr{SQL: "Embedding <=> ? ", Vars: []interface{}{question}}}).Limit(3).Where("classify = ?", classify).Find(&docsbody).Error
	return docsbody, err
}
