/*
 * @Author: cloudyi.li
 * @Date: 2023-04-04 19:47:43
 * @LastEditTime: 2023-04-05 15:56:43
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/dao/query/user.go
 */
package query

import (
	"chatserver-api/internal/dao"
	"chatserver-api/internal/model/entity"
	"chatserver-api/pkg/db"
	"context"
)

var _ dao.UserDao = (*userDao)(nil)

type userDao struct {
	ds db.IDataSource
}

func NewUserDao(_ds db.IDataSource) *userDao {
	return &userDao{
		ds: _ds,
	}
}

func (ud *userDao) GetUserByName(ctx context.Context, username string) (*entity.User, error) {
	user := &entity.User{}
	err := ud.ds.Master().Where("username = ?", username).Find(user).Error
	return user, err
}

func (ud *userDao) GetUserById(ctx context.Context, userid int64) (*entity.User, error) {
	user := &entity.User{}
	err := ud.ds.Master().Where("id = ?", userid).Find(user).Error
	return user, err
}
