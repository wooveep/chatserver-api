/*
 * @Author: cloudyi.li
 * @Date: 2023-05-15 11:39:37
 * @LastEditTime: 2023-05-15 11:43:07
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/model/entity/invite.go
 */
package entity

import (
	"chatserver-api/pkg/jtime"

	"gorm.io/plugin/soft_delete"
)

type Invite struct {
	Id          int64                 `gorm:"column:id;primary_key;" json:"id"`
	UserId      int64                 `gorm:"column:user_id" json:"user_id"`
	InviteCode  string                `gorm:"column:invite_code" json:"invite_code"`
	Number      int                   `gorm:"column:number" json:"number"`
	TotalReward float64               `gorm:"column:total_reward" json:"total_reward"`
	CreatedAt   jtime.JsonTime        `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   jtime.JsonTime        `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt   jtime.JsonTime        `gorm:"column:deleted_at" json:"deleted_at" `
	IsDel       soft_delete.DeletedAt `gorm:"softDelete:flag,DeletedAtField:DeletedAt"`
}

func (Invite) TableName() string {
	return "public.invite"
}
