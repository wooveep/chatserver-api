/*
 * @Author: cloudyi.li
 * @Date: 2023-04-05 15:37:14
 * @LastEditTime: 2023-04-05 15:56:29
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/dao/query/chat.go
 */
package query

import (
	"chatserver-api/internal/dao"
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

func (cd *chatDao) ChatCreateNew(ctx context.Context, chat *entity.ChatSession) error {
	return cd.ds.Master().Create(chat).Error
}
