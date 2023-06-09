/*
 * @Author: cloudyi.li
 * @Date: 2023-04-10 19:36:20
 * @LastEditTime: 2023-05-28 15:08:13
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/model/entity/preset.go
 */
package entity

import (
	"chatserver-api/pkg/jtime"

	"gorm.io/datatypes"
	"gorm.io/plugin/soft_delete"
)

type Preset struct {
	Id            int64                 `gorm:"column:id;primary_key;" json:"id"`
	PresetName    string                `gorm:"column:preset_name" json:"preset_name"`
	PresetContent string                `gorm:"column:preset_content" json:"preset_content"`
	PresetTips    string                `gorm:"column:preset_tips" json:"preset_tips"`
	ModelName     string                `gorm:"column:model_name" json:"model_name"`
	MaxTokens     int                   `gorm:"column:max_token" json:"max_token"`
	LogitBias     datatypes.JSON        `gorm:"column:logit_bias" json:"logit_bias"`
	Temperature   float64               `gorm:"column:temperature" json:"temperature"`
	TopP          float64               `gorm:"column:top_p" json:"top_p"`
	Presence      float64               `gorm:"column:presence" json:"presence"`
	Frequency     float64               `gorm:"column:frequency" json:"frequency"`
	WithEmbedding bool                  `grom:"cloumn:with_embedding" json:"with_embedding"`
	Classify      string                `gorm:"column:classify" json:"classify"`
	Extension     int                   `gorm:"column:extension" json:"extension"`
	Privilege     int                   `gorm:"column:privilege" json:"privilege"`
	CreatedAt     jtime.JsonTime        `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     jtime.JsonTime        `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt     jtime.JsonTime        `gorm:"column:deleted_at" json:"deleted_at"`
	IsDel         soft_delete.DeletedAt `gorm:"softDelete:flag,DeletedAtField:DeletedAt"`
	Chats         []Chat                `gorm:"foreignKey:preset_id;references:id"`
}

func (Preset) TableName() string {
	return "public.preset"
}
