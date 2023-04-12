/*
 * @Author: cloudyi.li
 * @Date: 2023-04-04 19:47:43
 * @LastEditTime: 2023-04-10 17:24:30
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

func (ud *userDao) UserRegisterNew(ctx context.Context, user *entity.User) error {
	return ud.ds.Master().Create(user).Error
}
func (ud *userDao) UserGetByName(ctx context.Context, username string) (*entity.User, error) {
	user := &entity.User{}
	err := ud.ds.Master().Where("username = ?", username).Find(user).Error
	return user, err
}

func (ud *userDao) UserGetById(ctx context.Context, userid int64) (*entity.User, error) {
	user := &entity.User{}
	err := ud.ds.Master().Where("id = ?", userid).Find(user).Error
	return user, err
}
func (ud *userDao) UserGetByEmail(ctx context.Context, email string) (*entity.User, error) {
	user := &entity.User{}
	err := ud.ds.Master().Where("email = ?", email).Find(user).Error
	return user, err
}

func (ud *userDao) UserVerifyEmail(ctx context.Context, email string) (count int64, err error) {
	err = ud.ds.Master().Model(&entity.User{}).Where("email = ?", email).Count(&count).Error
	return
}
func (ud *userDao) UserVerifyUserName(ctx context.Context, username string) (count int64, err error) {
	err = ud.ds.Master().Model(&entity.User{}).Where("username = ?", username).Count(&count).Error
	return
}

func (ud *userDao) UserUpdateNickName(ctx context.Context, userid int64, nickname string) error {
	return ud.ds.Master().Model(&entity.User{}).Where("id = ?", userid).Update("nickname", nickname).Error
}
