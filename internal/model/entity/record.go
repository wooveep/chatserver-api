/*
 * @Author: cloudyi.li
 * @Date: 2023-04-12 06:09:03
 * @LastEditTime: 2023-05-18 13:04:23
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/model/entity/record.go
 */
package entity

import (
	"chatserver-api/pkg/jtime"

	"gorm.io/plugin/soft_delete"
)

// UserId      int64          `gorm:"column:user_id" json:"user_id"`

type Record struct {
	Id           int64                 `gorm:"column:id;primary_key;" json:"id"`
	ChatId       int64                 `gorm:"column:chat_id" json:"chat_id" `
	Sender       string                `gorm:"column:sender" json:"sender" `
	Message      string                `gorm:"column:message" json:"message" `
	MessageHash  string                `gorm:"column:message_hash" json:"message_hash"`
	MessageToken int                   `gorm:"column:message_token" json:"message_token"`
	CreatedAt    jtime.JsonTime        `gorm:"column:created_at" json:"created_at" `
	UpdatedAt    jtime.JsonTime        `gorm:"column:updated_at" json:"updated_at" `
	DeletedAt    jtime.JsonTime        `gorm:"column:deleted_at" json:"deleted_at" `
	IsDel        soft_delete.DeletedAt `gorm:"softDelete:flag,DeletedAtField:DeletedAt"`
}

func (Record) TableName() string {
	return "public.record"
}
