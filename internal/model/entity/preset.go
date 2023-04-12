/*
 * @Author: cloudyi.li
 * @Date: 2023-04-10 19:36:20
 * @LastEditTime: 2023-04-12 20:06:45
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/model/entity/preset.go
 */
package entity

import (
	"chatserver-api/pkg/jtime"

	"gorm.io/datatypes"
)

type Preset struct {
	Id            int64          `gorm:"column:id;primary_key;" json:"id"`
	PresetName    string         `gorm:"column:preset_name" json:"preset_name"`
	PresetContent string         `gorm:"column:preset_content" json:"preset_content"`
	ModelName     string         `gorm:"column:model_name" json:"model_name"`
	MaxTokens     int            `gorm:"column:max_token" json:"max_token"`
	LogitBias     datatypes.JSON `gorm:"column:logit_bias" json:"logit_bias"`
	Temperature   float64        `gorm:"column:temperature" json:"temperature"`
	TopP          float64        `gorm:"column:top_p" json:"top_p"`
	Presence      float64        `gorm:"column:presence" json:"presence"`
	Frequency     float64        `gorm:"column:frequency" json:"frequency"`
	CreatedAt     jtime.JsonTime `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     jtime.JsonTime `gorm:"column:updated_at" json:"updated_at"`
	Chats         []Chat         `gorm:"foreignKey:preset_id;references:id"`
}

func (Preset) TableName() string {
	return "public.preset"
}
