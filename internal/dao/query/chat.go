/*
 * @Author: cloudyi.li
 * @Date: 2023-04-05 15:37:14
 * @LastEditTime: 2023-04-13 15:39:47
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/dao/query/chat.go
 */
package query

import (
	"chatserver-api/internal/dao"
	"chatserver-api/internal/model"
	"chatserver-api/internal/model/entity"
	"chatserver-api/pkg/db"
	"context"
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

func (cd *chatDao) ChatRecordSave(ctx context.Context, record *entity.Record) error {
	return cd.ds.Master().Create(record).Error
}

func (cd *chatDao) ChatDeleteOne(ctx context.Context, userid, chatid int64) error {
	return cd.ds.Master().Where("user_id = ?", userid).Where("id = ? ", chatid).Delete(&entity.Chat{}).Error
}

func (cd *chatDao) ChatDeleteAll(ctx context.Context, userid int64) error {
	return cd.ds.Master().Where("user_id = ?", userid).Delete(&entity.Chat{}).Error
}

func (cd *chatDao) ChatDetailGet(ctx context.Context, userid, chatid int64) (model.ChatDetail, error) {
	var detail model.ChatDetail
	err := cd.ds.Master().Joins("Chats", cd.ds.Master().Where(&entity.Chat{Id: chatid, UserId: userid})).Model(&entity.Preset{}).Find(&detail).Error
	return detail, err
}

func (cd *chatDao) ChatRecordGet(ctx context.Context, chatid int64, memory int16) ([]model.RecordOne, error) {
	var recordlist []model.RecordOne
	subQuery1 := cd.ds.Master().Model(&entity.Record{}).Where("chat_id = ?", chatid).Order("id desc").Limit(int(memory)).Select("*")
	err := cd.ds.Master().Table("(?) as a ", subQuery1).Order("id").Find(&recordlist).Error
	return recordlist, err
}

func (cd *chatDao) ChatGetList(ctx context.Context, userid int64) ([]model.ChatOne, error) {
	var chatlist []model.ChatOne
	err := cd.ds.Master().Model(&entity.Chat{}).Where("user_id = ? ", userid).Find(&chatlist).Error
	return chatlist, err
}

func (cd *chatDao) ChatUserVerify(ctx context.Context, userid, chatid int64) (int64, error) {
	var count int64
	err := cd.ds.Master().Model(&entity.Chat{}).Where("id = ? ", chatid).Where("user_id = ?", userid).Count(&count).Error
	return count, err
}

func (cd *chatDao) ChatCostUpdate(ctx context.Context, userid int64, balance float64) error {
	return cd.ds.Master().Model(&entity.User{}).Where("id  = ? ", userid).UpdateColumn("balance", balance).Error
}

func (cd *chatDao) ChatBalanceGet(ctx context.Context, userid int64) (model.UserBalance, error) {
	var userbalance model.UserBalance
	err := cd.ds.Master().Model(&entity.User{}).Where("id  = ? ", userid).Find(&userbalance).Error
	return userbalance, err
}
