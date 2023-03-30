/*
 * @Author: cloudyi.li
 * @Date: 2023-03-29 12:36:21
 * @LastEditTime: 2023-03-29 13:44:57
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/handler/v1/user/user.go
 */
package user

import (
	"chatserver-api/internal/constant"
	"chatserver-api/internal/service"
	"chatserver-api/pkg/errors"
	"chatserver-api/pkg/errors/ecode"
	"chatserver-api/pkg/response"
	"context"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userSrv service.UserService
}

func NewUserHandler(_userSrv service.UserService) *UserHandler {
	return &UserHandler{
		userSrv: _userSrv,
	}
}

func (uh *UserHandler) GetAvatar() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetInt64(constant.UserID)
		user, err := uh.userSrv.GetByID(context.TODO(), id)
		if err != nil {
			response.JSON(c, errors.Wrap(err, ecode.NotFoundErr, "用户信息为空"), nil)
		} else {
			response.JSON(c, nil, user)
		}
	}
}
