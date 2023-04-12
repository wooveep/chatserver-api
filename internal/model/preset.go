/*
 * @Author: cloudyi.li
 * @Date: 2023-04-10 19:41:56
 * @LastEditTime: 2023-04-12 18:19:47
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
}
type PresetCreateNewRes struct {
	Id        int64 `json:"id"`
	IsSuccess bool  `json:"is_success"`
}

type PresetGetListRes struct {
	PresetsList []PresetOne `json:"prests_list"`
}

type PresetOne struct {
	Id         int64  `json:"id"`
	PresetName string `json:"preset_name"`
}
