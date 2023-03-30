/*
 * @Author: cloudyi.li
 * @Date: 2023-03-29 12:37:13
 * @LastEditTime: 2023-03-29 12:52:47
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/service/user.go
 */
package service

import (
	"context"
)

var _ UserService = (*userService)(nil)

type UserService interface {
	GetByID(ctx context.Context, uid int64) (string, error)
}

// userService 实现UserService接口
type userService struct {
}

func NewUserService() *userService {
	return &userService{}
}

// GetByName 通过用户名 查找用户
func (us *userService) GetByID(ctx context.Context, uid int64) (string, error) {
	return "admin", nil
}
