/*
 * @Author: cloudyi.li
 * @Date: 2023-04-04 15:44:35
 * @LastEditTime: 2023-04-08 15:18:07
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/model/user.go
 */
package model

import "chatserver-api/pkg/jtime"

type UserLoginReq struct {
	Username string `json:"username" validate:"required"  label:"用户名"`
	Password string `json:"password" validate:"required"  label:"密码"`
}
type UserLoginRes struct {
	Token    string         `json:"token"`
	ExpireAt jtime.JsonTime `json:"expire_at"`
}
type UserRegisterReq struct {
	Username string `json:"username" validate:"required"  label:"用户名"`
	Password string `json:"password" validate:"required"  label:"密码"`
	Email    string `json:"email" validate:"required"  label:"邮箱地址"`
}

type UserRegisterRes struct {
	IsSuccess bool `json:"is_success"`
}
type UserVerifyUserNameReq struct {
	Username string `json:"username" validate:"required"  label:"用户名"`
}

type UserVerifyUserNameRes struct {
	Isvalid bool `json:"is_valid"`
}

type UserVerifyEmailReq struct {
	Email string `json:"email" validate:"required"  label:"邮箱地址"`
}

type UserVerifyEmailRes struct {
	Isvalid bool `json:"is_valid"`
}

type UserGetAvatarRes struct {
	AvatarUrl string `json:"avatar_url"`
}

type UserGetInfoRes struct {
	Username  string  `json:"username"`
	Nickname  string  `json:"nickname"`
	Email     string  `json:"email"`
	Phone     string  `json:"phone"`
	AvatarUrl string  `json:"avatar_url"`
	Balance   float64 `json:"balance"`
}
