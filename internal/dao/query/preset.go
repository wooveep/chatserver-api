/*
 * @Author: cloudyi.li
 * @Date: 2023-04-11 09:49:32
 * @LastEditTime: 2023-05-15 11:57:45
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/dao/query/preset.go
 */
package query

import (
	"chatserver-api/internal/dao"
	"chatserver-api/internal/model"
	"chatserver-api/internal/model/entity"
	"chatserver-api/pkg/db"
	"context"
)

var _ dao.PresetDao = (*presetDao)(nil)

type presetDao struct {
	ds db.IDataSource
}

func NewPresetsDao(_ds db.IDataSource) *presetDao {
	return &presetDao{
		ds: _ds,
	}
}

func (pd *presetDao) PresetCreateNew(ctx context.Context, preset *entity.Preset) error {
	return pd.ds.Master().Create(preset).Error
}

func (pd *presetDao) PresetUpdate(ctx context.Context, preset *entity.Preset) error {
	return pd.ds.Master().Updates(preset).Error
}

func (pd *presetDao) PresetGetList(ctx context.Context) ([]model.PresetOne, error) {
	var presetlist []model.PresetOne
	err := pd.ds.Master().Model(&entity.Preset{}).Find(&presetlist).Error
	return presetlist, err
}
