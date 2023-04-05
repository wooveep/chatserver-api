/*
 * @Author: cloudyi.li
 * @Date: 2023-04-04 19:43:52
 * @LastEditTime: 2023-04-04 19:48:02
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/dao/user.go
 */
package dao

import (
	"chatserver-api/internal/model/entity"
	"context"
)

type UserDao interface {
	GetUserByName(ctx context.Context, name string) (*entity.User, error)
	GetUserById(ctx context.Context, uid int64) (*entity.User, error)
}
