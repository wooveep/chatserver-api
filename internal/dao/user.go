/*
 * @Author: cloudyi.li
 * @Date: 2023-04-04 19:43:52
 * @LastEditTime: 2023-05-12 16:16:09
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/dao/user.go
 */
package dao

import (
	"chatserver-api/internal/model/entity"
	"context"
)

type UserDao interface {
	UserGetByName(ctx context.Context, username string) (entity.User, error)
	UserGetById(ctx context.Context, uid int64) (entity.User, error)
	UserCreate(ctx context.Context, user *entity.User) error
	UserVerifyEmail(ctx context.Context, email string) (userId int64, err error)
	UserVerifyUserName(ctx context.Context, username string) (count int64, err error)
	UserUpdateNickName(ctx context.Context, userId int64, nickname string) error
	UserDelete(ctx context.Context, userId int64) error
	UserUpdate(ctx context.Context, user *entity.User) error
}
