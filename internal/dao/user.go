/*
 * @Author: cloudyi.li
 * @Date: 2023-04-04 19:43:52
 * @LastEditTime: 2023-05-26 10:17:17
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/dao/user.go
 */
package dao

import (
	"chatserver-api/internal/model"
	"chatserver-api/internal/model/entity"
	"context"
)

type UserDao interface {
	UserGetByName(ctx context.Context, username string) (entity.User, error)
	UserGetById(ctx context.Context, userId int64) (model.UserInfo, error)
	UserCreate(ctx context.Context, user *entity.User) error
	UserVerifyEmail(ctx context.Context, email string) (userId int64, err error)
	UserVerifyUserName(ctx context.Context, username string) (count int64, err error)
	UserUpdateNickName(ctx context.Context, userId int64, nickname string) error
	UserDelete(ctx context.Context, userId int64) error
	UserUpdate(ctx context.Context, user *entity.User) error
	UserGetRole(ctx context.Context, userId int64) (int, error)
	UserGetAvatar(ctx context.Context, userId int64) (string, error)
	UserGetBalance(ctx context.Context, userId int64) (float64, error)
	UserInviteGen(ctx context.Context, invite *entity.Invite) error
	UserInviteGetByUser(ctx context.Context, userId int64) (entity.Invite, error)
	UserInviteGetByCode(ctx context.Context, code string) (entity.Invite, error)
	UserInviteUpdate(ctx context.Context, invite *entity.Invite) error
	UserBillCreate(ctx context.Context, bill *entity.Bill) error
	UserBillGet(ctx context.Context, userId int64, page, pagesize int, start, end string) ([]model.UserBillRes, error)
}
