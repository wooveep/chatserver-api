/*
 * @Author: cloudyi.li
 * @Date: 2023-03-29 12:37:13
 * @LastEditTime: 2023-04-05 08:56:36
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/service/user.go
 */
package service

import (
	"chatserver-api/internal/dao"
	"chatserver-api/internal/model/entity"
	"context"
)

var _ UserService = (*userService)(nil)

type UserService interface {
	GetByID(ctx context.Context, uid int64) (*entity.User, error)
}

// userService 实现UserService接口
type userService struct {
	ud dao.UserDao
}

func NewUserService(_ud dao.UserDao) *userService {
	return &userService{
		ud: _ud,
	}
}

// GetByName 通过用户名 查找用户
func (us *userService) GetByID(ctx context.Context, uid int64) (*entity.User, error) {
	return us.ud.GetUserById(ctx, uid)
}
