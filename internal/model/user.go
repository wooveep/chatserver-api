/*
 * @Author: cloudyi.li
 * @Date: 2023-04-04 15:44:35
 * @LastEditTime: 2023-05-12 16:56:08
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
	Username string `json:"username" validate:"required,username"  label:"用户名"`
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
	Username  string         `json:"username"`
	Nickname  string         `json:"nickname"`
	Email     string         `json:"email"`
	Phone     string         `json:"phone"`
	AvatarUrl string         `json:"avatar_url"`
	ExpiredAt jtime.JsonTime `json:"expired_at"`
	Balance   float64        `json:"balance"`
}

type UserUpdateNickNameReq struct {
	Nickname string `json:"nickname" validate:"required"  label:"用户昵称"`
}
type UserUpdateNickNameRes struct {
	IsChanged bool `json:"is_changed"`
}

type UserBalance struct {
	Balance float64 `json:"balance"`
}

type UserPasswordModifyReq struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
}

type UserPasswordResetReq struct {
	TempCode    string `json:"temp_code" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
}

type UserForgetReq struct {
	Email string `json:"email"`
}
