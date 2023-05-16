/*
 * @Author: cloudyi.li
 * @Date: 2023-04-04 19:47:43
 * @LastEditTime: 2023-05-15 16:53:33
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

func (ud *userDao) UserCreate(ctx context.Context, user *entity.User) error {
	return ud.ds.Master().Create(user).Error
}
func (ud *userDao) UserUpdate(ctx context.Context, user *entity.User) error {
	return ud.ds.Master().Updates(user).Error
}
func (ud *userDao) UserDelete(ctx context.Context, userId int64) error {
	return ud.ds.Master().Delete(&entity.User{Id: userId}).Error
}
func (ud *userDao) UserGetRole(ctx context.Context, userId int64) (int, error) {
	var role int
	err := ud.ds.Master().Model(&entity.User{}).Where("id = ?", userId).Select("role").Find(&role).Error
	return role, err
}

func (ud *userDao) UserGetByName(ctx context.Context, username string) (entity.User, error) {
	var user entity.User
	err := ud.ds.Master().Where("username = ?", username).Find(&user).Error
	return user, err
}

func (ud *userDao) UserGetById(ctx context.Context, userId int64) (entity.User, error) {
	var user entity.User
	err := ud.ds.Master().Where("id = ?", userId).Find(&user).Error
	return user, err
}
func (ud *userDao) UserGetByEmail(ctx context.Context, email string) (entity.User, error) {
	var user entity.User
	err := ud.ds.Master().Where("email = ?", email).Find(&user).Error
	return user, err
}

func (ud *userDao) UserVerifyEmail(ctx context.Context, email string) (userId int64, err error) {

	err = ud.ds.Master().Model(&entity.User{}).Where("email = ?", email).Select("id").Find(&userId).Error
	return
}
func (ud *userDao) UserVerifyUserName(ctx context.Context, username string) (count int64, err error) {
	err = ud.ds.Master().Model(&entity.User{}).Where("username = ?", username).Count(&count).Error
	return
}

func (ud *userDao) UserUpdateNickName(ctx context.Context, userId int64, nickname string) error {
	return ud.ds.Master().Model(&entity.User{}).Where("id = ?", userId).Update("nickname", nickname).Error
}
