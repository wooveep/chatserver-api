/*
 * @Author: cloudyi.li
 * @Date: 2023-05-10 09:58:37
 * @LastEditTime: 2023-05-10 10:44:06
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/model/entity/cdkey.go
 */
package entity

import (
	"chatserver-api/pkg/jtime"

	"gorm.io/plugin/soft_delete"
)

type CdKey struct {
	Id         int64                 `gorm:"column:id;primary_key;" json:"id"`
	UseId      int64                 `gorm:"column:user_id" json:"user_id"`
	InviteCode string                `gorm:"column:invite_code" json:"invite_code"`
	InviteNum  int                   `gorm:"column:invite_num" json:"invite_num"`
	ExpireAt   jtime.JsonTime        `gorm:"column:expired_at" json:"expired_at"`
	CreatedAt  jtime.JsonTime        `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  jtime.JsonTime        `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt  jtime.JsonTime        `gorm:"column:deleted_at" json:"deleted_at" `
	IsDel      soft_delete.DeletedAt `gorm:"softDelete:flag,DeletedAtField:DeletedAt"`
}

func (CdKey) TableName() string {
	return "public.cdkey"
}
