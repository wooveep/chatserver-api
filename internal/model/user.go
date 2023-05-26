/*
 * @Author: cloudyi.li
 * @Date: 2023-04-04 15:44:35
 * @LastEditTime: 2023-05-26 13:09:51
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/model/user.go
 */
package model

import "chatserver-api/pkg/jtime"

type UserLoginReq struct {
	Username string `json:"username" validate:"required"  label:"用户名"`
	Password string `json:"password" validate:"required"  label:"密码"`
	Captcha  string `json:"captcha" validate:"required"  label:"验证码"`
}
type UserLoginRes struct {
	Token   string `json:"token"`
	TimeOut int    `json:"timeout"`
}
type UserRegisterReq struct {
	Username   string `json:"username" validate:"required,username"  label:"用户名"`
	Password   string `json:"password" validate:"required"  label:"密码"`
	Email      string `json:"email" validate:"required"  label:"邮箱地址"`
	Captcha    string `json:"captcha" validate:"required"  label:"验证码"`
	InviteCode string `json:"invite_code" label:"邀请码"`
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
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
}

type UserInfo struct {
	Username string `gorm:"column:username" json:"username"`
	Nickname string `gorm:"column:nickname" json:"nickname"`
	Password string `gorm:"column:password" json:"password"`
	Email    string `gorm:"column:email" json:"email"`
	Phone    string `gorm:"column:phone" json:"phone"`
	Role     int    `gorm:"column:role" json:"role"`
	IsActive bool   `gorm:"column:is_active" json:"is_active"`
}
type UserAvatarRes struct {
	Avatar string `json:"avatar"`
}

type UserActiveReq struct {
	ActiveCode string `form:"active_code"  validate:"required"`
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
	Email   string `json:"email"`
	Captcha string `json:"captcha" validate:"required"  label:"验证码"`
}

type UserInviteLinkRes struct {
	InviteLink   string  `json:"invite_link"`
	InviteNumber int     `json:"invite_number"`
	InviteReward float64 `json:"invite_reward"`
}

type CaptchaRes struct {
	Image string `json:"image"`
}

type UserBillGetReq struct {
	Page     int `form:"page"`
	PageSize int `form:"pagesize" validate:"required"`
	Date     int `form:"date"`
}

type UserBillRes struct {
	CreatedAt   jtime.JsonTime `gorm:"column:created_at" json:"change_time"`
	CostChange  float64        `gorm:"column:cost_change" json:"change"`
	Balance     float64        `gorm:"column:balance" json:"balance"`
	CostComment string         `gorm:"column:cost_comment" json:"comment"`
}

type UserBillListRes struct {
	BillList []UserBillRes `json:"bill_list"`
}
