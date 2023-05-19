/*
 * @Author: cloudyi.li
 * @Date: 2023-05-15 11:39:37
 * @LastEditTime: 2023-05-18 14:10:14
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/model/entity/invite.go
 */
package entity

import (
	"chatserver-api/pkg/jtime"
)

type Invite struct {
	Id           int64          `gorm:"column:id;primary_key;" json:"id"`
	UserId       int64          `gorm:"column:user_id" json:"user_id"`
	InviteCode   string         `gorm:"column:invite_code" json:"invite_code"`
	InviteNumber int            `gorm:"column:invite_number" json:"invite_number"`
	CreatedAt    jtime.JsonTime `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    jtime.JsonTime `gorm:"column:updated_at" json:"updated_at"`
}

func (Invite) TableName() string {
	return "public.invite"
}
