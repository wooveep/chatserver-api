/*
 * @Author: cloudyi.li
 * @Date: 2023-04-10 19:41:56
 * @LastEditTime: 2023-05-06 22:41:58
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/model/preset.go
 */
package model

import "gorm.io/datatypes"

type PresetCreateNewReq struct {
	PresetName    string         `json:"preset_name"  validate:"required"`
	PresetContent string         `json:"preset_content"  validate:"required"`
	ModelName     string         `json:"model_name"`
	MaxTokens     int            `json:"max_token"`
	LogitBias     datatypes.JSON `json:"logit_bias"`
	Temperature   float64        `json:"temperature"`
	TopP          float64        `json:"top_p"`
	Presence      float64        `json:"presence"`
	Frequency     float64        `json:"frequency"`
	WithEmbedding bool           `json:"with_embedding"`
}
type PresetCreateNewRes struct {
	PresetId  int64 `json:"preset_id"`
	IsSuccess bool  `json:"is_success"`
}

type PresetGetListRes struct {
	PresetsList []PresetOneRes `json:"preset_list"`
}
type PresetOneRes struct {
	PresetId   string `json:"preset_id"`
	PresetName string `json:"preset_name"`
}

type PresetOne struct {
	PresetId   int64  `gorm:"column:id" json:"preset_id"`
	PresetName string `json:"preset_name"`
}
