/*
 * @Author: cloudyi.li
 * @Date: 2023-04-05 15:39:43
 * @LastEditTime: 2023-05-08 17:28:53
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/model/entity/chat.go
 */
package entity

import "chatserver-api/pkg/jtime"

type Chat struct {
	Id        int64          `gorm:"column:id;primary_key;" json:"id"`
	UserId    int64          `gorm:"column:user_id" json:"user_id"`
	PresetId  int64          `gorm:"column:preset_id" json:"preset_id"`
	ChatName  string         `gorm:"column:chat_name" json:"chat_name"`
	CreatedAt jtime.JsonTime `gorm:"column:created_at" json:"created_at"`
	UpdatedAt jtime.JsonTime `gorm:"column:updated_at" json:"updated_at"`
	Records   []Record       `gorm:"foreignKey:chat_id;references:id"`
}

func (Chat) TableName() string {
	return "public.chat"
}
