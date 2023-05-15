/*
 * @Author: cloudyi.li
 * @Date: 2023-04-11 09:50:36
 * @LastEditTime: 2023-05-15 11:57:58
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/dao/preset.go
 */
package dao

import (
	"chatserver-api/internal/model"
	"chatserver-api/internal/model/entity"
	"context"
)

type PresetDao interface {
	PresetCreateNew(ctx context.Context, preset *entity.Preset) error
	PresetGetList(ctx context.Context) ([]model.PresetOne, error)
	PresetUpdate(ctx context.Context, preset *entity.Preset) error
}
