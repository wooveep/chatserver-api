/*
 * @Author: cloudyi.li
 * @Date: 2023-04-05 15:37:14
 * @LastEditTime: 2023-04-12 20:36:26
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

func (cd *chatDao) ChatDetailGet(ctx context.Context, userid, chatid int64) (model.ChatDetail, error) {
	var detail model.ChatDetail
	err := cd.ds.Master().Joins("Chats", "id = ?", chatid).Model(&entity.Preset{}).Where("user_id = ? ", userid).Scan(&detail).Error
	return detail, err
}

func (cd *chatDao) ChatRecordGet(ctx context.Context, userid int64, chatid int64, memory int16) ([]model.RecordOne, error) {
	var recordlist []model.RecordOne
	err := cd.ds.Master().Joins("Records", "chat_id = ?", chatid).Model(&entity.Chat{}).Where("user_id = ? ", userid).Order("id desc").Limit(int(memory)).Find(&recordlist).Error
	return recordlist, err
}

func (cd *chatDao) ChatGetList(ctx context.Context, userid int64) ([]model.ChatOne, error) {
	var chatlist []model.ChatOne
	err := cd.ds.Master().Model(&entity.Chat{}).Where("user_id = ? ", userid).Find(&chatlist).Error
	return chatlist, err
}
