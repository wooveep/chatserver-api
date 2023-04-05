/*
 * @Author: cloudyi.li
 * @Date: 2023-04-05 15:39:43
 * @LastEditTime: 2023-04-05 15:43:09
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/model/entity/chat.go
 */
package entity

import "chatserver-api/pkg/time"

type ChatSession struct {
	Id          int            `gorm:"column:id" json:"id"`
	UserId      int            `gorm:"column:user_id" json:"user_id"`
	PresetId    int            `gorm:"column:preset_id" json:"preset_id"`
	SessionName string         `gorm:"column:session_name" json:"session_name"`
	MemoryLevel int            `gorm:"column:memory_level" json:"memory_level"`
	CreatedAt   *time.JsonTime `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   *time.JsonTime `gorm:"column:updated_at" json:"updated_at"`
}

func (ChatSession) TableName() string {
	return "public.chatsession"
}
