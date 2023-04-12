/*
 * @Author: cloudyi.li
 * @Date: 2023-03-29 14:11:49
 * @LastEditTime: 2023-04-12 20:10:26
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/model/chat.go
 */
package model

import (
	"chatserver-api/pkg/jtime"

	"gorm.io/datatypes"
)

type ChatCreateNewReq struct {
	ChatName    string `json:"chatname" validate:"required" label:"会话名称"`
	PresetId    int64  `json:"presetid" validate:"required" label:"预设ID"`
	MemoryLevel int16  `json:"memorylevel" validate:"required" label:"消息记忆"`
}
type ChatCreateNewRes struct {
	ChatId int64 `json:"chatid"`
}

type ChatChattingReq struct {
	ChatId  int64  `json:"chatid" validate:"required" label:"会话ID"`
	Message string `json:"message" validate:"required" label:"消息"`
}

type ChatChattingRes struct {
}

type ChatListRes struct {
	ChatList []ChatOne `json:"chat_list"`
}
type ChatOne struct {
	Id        int64          `json:"id"`
	ChatName  string         `json:"chat_name"`
	CreatedAt jtime.JsonTime `json:"created_at"`
}

type ChatDetail struct {
	PresetName    string         `gorm:"column:preset_name" json:"preset_name"`
	PresetContent string         `gorm:"column:preset_content" json:"preset_content"`
	ModelName     string         `gorm:"column:model_name" json:"model_name"`
	MaxTokens     int            `gorm:"column:max_token" json:"max_token"`
	LogitBias     datatypes.JSON `gorm:"column:logit_bias" json:"logit_bias"`
	Temperature   float64        `gorm:"column:temperature" json:"temperature"`
	TopP          float64        `gorm:"column:top_p" json:"top_p"`
	Presence      float64        `gorm:"column:presence" json:"presence"`
	Frequency     float64        `gorm:"column:frequency" json:"frequency"`
	MemoryLevel   int16          `gorm:"column:Chats__memory_level" json:"memorylevel"`
}

type ChatDetailReq struct {
	Id int64 `json:"id"`
}
type ChatDetailRes struct {
	PresetContent string `json:"preset_content"`
	ModelName     string `json:"model_name"`
	MaxTokens     int    `json:"max_token"`
	MemoryLevel   int16  `json:"memorylevel"`
}
