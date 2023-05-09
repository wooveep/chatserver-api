/*
 * @Author: cloudyi.li
 * @Date: 2023-04-12 06:09:03
 * @LastEditTime: 2023-05-06 21:29:19
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/model/entity/record.go
 */
package entity

import "chatserver-api/pkg/jtime"

// UserId      int64          `gorm:"column:user_id" json:"user_id"`

type Record struct {
	Id          int64          `gorm:"column:id;primary_key;" json:"id"`
	ChatId      int64          `gorm:"column:chat_id" json:"chat_id" `
	Sender      string         `gorm:"column:sender" json:"sender" `
	Message     string         `gorm:"column:message" json:"message" `
	MessageHash string         `gorm:"column:message_hash" json:"message_hash"`
	CreatedAt   jtime.JsonTime `gorm:"column:created_at" json:"created_at" `
	UpdatedAt   jtime.JsonTime `gorm:"column:updated_at" json:"updated_at" `
}

func (Record) TableName() string {
	return "public.record"
}
