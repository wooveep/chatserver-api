/*
 * @Author: cloudyi.li
 * @Date: 2023-03-29 12:36:21
 * @LastEditTime: 2023-04-08 16:07:16
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/handler/v1/user/user.go
 */
package user

import (
	"chatserver-api/internal/consts"
	"chatserver-api/internal/model"
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

func (uh *UserHandler) UserGetAvatar() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetInt64(consts.UserID)
		useravatar, err := uh.userSrv.UserGetAvatar(context.TODO(), id)
		if err != nil {
			response.JSON(c, errors.Wrap(err, ecode.NotFoundErr, "未找到头像"), nil)
		} else {
			response.JSON(c, nil, useravatar)
		}
	}
}
func (uh *UserHandler) UserGetInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetInt64(consts.UserID)
		userinfo, err := uh.userSrv.UserGetInfo(context.TODO(), id)
		if err != nil {
			response.JSON(c, errors.Wrap(err, ecode.NotFoundErr, "未找到用户信息"), nil)
		} else {
			response.JSON(c, nil, userinfo)
		}
	}
}

func (uh *UserHandler) UserRegister() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req model.UserRegisterReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			response.JSON(ctx, errors.WithCode(ecode.ValidateErr, err.Error()), nil)
			return
		}
		res, err := uh.userSrv.UserRegister(ctx, req)
		if err != nil {
			response.JSON(ctx, errors.Wrap(err, ecode.Unknown, "未知错误注册失败"), res)
		} else {
			response.JSON(ctx, errors.Wrap(err, ecode.Success, "注册成功"), res)
		}
	}
}
func (uh *UserHandler) UserLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req model.UserLoginReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			response.JSON(ctx, errors.WithCode(ecode.ValidateErr, err.Error()), nil)
			return
		}
		res, err := uh.userSrv.UserLogin(ctx, req.Username, req.Password)
		if err != nil {
			response.JSON(ctx, errors.Wrap(err, ecode.UserLoginErr, "登录失败；账户或密码错误"), nil)
		} else {
			response.JSON(ctx, errors.Wrap(err, ecode.Success, "登录成功"), res)
		}
	}
}

func (uh *UserHandler) UserLogout() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenstr := ctx.GetString(consts.TokenCtx)
		err := uh.userSrv.UserLogout(tokenstr)
		if err != nil {
			response.JSON(ctx, errors.Wrap(err, ecode.UserLoginErr, "登出失败"), nil)
		} else {
			response.JSON(ctx, errors.Wrap(err, ecode.Success, "登出成功"), nil)
		}
	}
}
