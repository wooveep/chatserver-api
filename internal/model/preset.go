/*
 * @Author: cloudyi.li
 * @Date: 2023-04-10 19:41:56
 * @LastEditTime: 2023-06-15 05:54:57
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/model/preset.go
 */
package model

import "gorm.io/datatypes"

type PresetCreateNewReq struct {
	PresetName    string         `json:"preset_name"  validate:"required"`
	PresetContent string         `json:"preset_content"  validate:"required"`
	PresetTips    string         `json:"preset_tips"  validate:"required"`
	ModelName     string         `json:"model_name"`
	MaxTokens     int            `json:"max_token"`
	LogitBias     datatypes.JSON `json:"logit_bias"`
	Temperature   float64        `json:"temperature"`
	TopP          float64        `json:"top_p"`
	Presence      float64        `json:"presence"`
	Frequency     float64        `json:"frequency"`
	WithEmbedding bool           `json:"with_embedding"`
	Classify      string         `json:"classify"`
	Extension     int            `json:"extension" validate:"required"`
	Privilege     int            `json:"privilege" validate:"required"`
}
type PresetCreateNewRes struct {
	PresetId  int64 `json:"preset_id"`
	IsSuccess bool  `json:"is_success"`
}

type PresetUpdateReq struct {
	PresetId      string         `json:"preset_id"  validate:"required"`
	PresetName    string         `json:"preset_name"`
	PresetContent string         `json:"preset_content"`
	PresetTips    string         `json:"preset_tips"`
	ModelName     string         `json:"model_name"`
	MaxTokens     int            `json:"max_token"`
	LogitBias     datatypes.JSON `json:"logit_bias"`
	Temperature   float64        `json:"temperature"`
	TopP          float64        `json:"top_p"`
	Presence      float64        `json:"presence"`
	Frequency     float64        `json:"frequency"`
	WithEmbedding bool           `json:"with_embedding"`
	Classify      string         `json:"classify"`
	Extension     int            `json:"extension"`
	Privilege     int            `json:"privilege" validate:"required"`
}

type PresetGetListRes struct {
	PresetsList []PresetOneRes `json:"preset_list"`
}
type PresetOneRes struct {
	PresetId   string `json:"preset_id"`
	PresetName string `json:"preset_name"`
	PresetTips string `json:"preset_tips"`
}

type PresetOne struct {
	PresetId   int64  `gorm:"column:id" json:"preset_id"`
	PresetName string `json:"preset_name"`
	PresetTips string `json:"preset_tips"`
}
