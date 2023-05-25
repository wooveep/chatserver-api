/*
 * @Author: cloudyi.li
 * @Date: 2023-04-04 19:47:43
 * @LastEditTime: 2023-05-25 16:23:22
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/dao/query/user.go
 */
package query

import (
	"chatserver-api/internal/dao"
	"chatserver-api/internal/model"
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

func (ud *userDao) UserBillCreate(ctx context.Context, bill *entity.Bill) error {
	return ud.ds.Master().Create(bill).Error
}

func (ud *userDao) UserGetRole(ctx context.Context, userId int64) (int, error) {
	var role int
	err := ud.ds.Master().Model(&entity.User{}).Where("id = ?", userId).Select("role").Find(&role).Error
	return role, err
}

func (ud *userDao) UserInviteGen(ctx context.Context, invite *entity.Invite) error {
	return ud.ds.Master().Create(invite).Error
}

func (ud *userDao) UserInviteUpdate(ctx context.Context, invite *entity.Invite) error {
	return ud.ds.Master().Updates(invite).Error
}

func (ud *userDao) UserInviteGetByCode(ctx context.Context, code string) (entity.Invite, error) {
	var inviteCode entity.Invite
	err := ud.ds.Master().Where("invite_code = ?", code).Find(&inviteCode).Error
	return inviteCode, err
}

func (ud *userDao) UserInviteGetByUser(ctx context.Context, userId int64) (entity.Invite, error) {
	var inviteCode entity.Invite
	err := ud.ds.Master().Where("user_id = ?", userId).Find(&inviteCode).Error
	return inviteCode, err
}

func (ud *userDao) UserGetBalance(ctx context.Context, userId int64) (float64, error) {
	var balance float64
	err := ud.ds.Master().Model(&entity.User{}).Where("id = ?", userId).Select("balance").Find(&balance).Error
	return balance, err
}

func (ud *userDao) UserGetByName(ctx context.Context, username string) (entity.User, error) {
	var user entity.User
	err := ud.ds.Master().Where("username = ?", username).Find(&user).Error
	return user, err
}

func (ud *userDao) UserGetById(ctx context.Context, userId int64) (model.UserInfo, error) {
	var user model.UserInfo
	err := ud.ds.Master().Model(&entity.User{}).Where("id = ?", userId).Find(&user).Error
	return user, err
}
func (ud *userDao) UserGetAvatar(ctx context.Context, userId int64) (string, error) {
	var avatar string
	err := ud.ds.Master().Model(&entity.User{}).Where("id = ?", userId).Select("avatar_url").Find(&avatar).Error
	return avatar, err
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
