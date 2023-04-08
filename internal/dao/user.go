/*
 * @Author: cloudyi.li
 * @Date: 2023-04-04 19:43:52
 * @LastEditTime: 2023-04-08 14:51:30
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/dao/user.go
 */
package dao

import (
	"chatserver-api/internal/model/entity"
	"context"
)

type UserDao interface {
	UserGetByName(ctx context.Context, username string) (*entity.User, error)
	UserGetById(ctx context.Context, uid int64) (*entity.User, error)
	UserRegisterNew(ctx context.Context, user *entity.User) error
}
