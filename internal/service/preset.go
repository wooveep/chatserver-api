/*
 * @Author: cloudyi.li
 * @Date: 2023-04-11 10:22:31
 * @LastEditTime: 2023-05-10 21:29:58
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/service/preset.go
 */
package service

import (
	"chatserver-api/internal/dao"
	"chatserver-api/internal/model"
	"chatserver-api/internal/model/entity"
	"chatserver-api/pkg/openai"
	"chatserver-api/utils/tools"
	"chatserver-api/utils/uuid"
	"context"
	"strconv"

	"gorm.io/datatypes"
)

var _ PresetService = (*presetService)(nil)

type PresetService interface {
	PresetCreateNew(ctx context.Context, req *model.PresetCreateNewReq) (res model.PresetCreateNewRes, err error)
	PresetGetList(ctx context.Context) (res model.PresetGetListRes, err error)
}

// userService 实现UserService接口
type presetService struct {
	pd   dao.PresetDao
	iSrv uuid.SnowNode
}

func NewPresetService(_pd dao.PresetDao) *presetService {
	return &presetService{
		pd:   _pd,
		iSrv: *uuid.NewNode(2),
	}
}

func (ps *presetService) PresetCreateNew(ctx context.Context, req *model.PresetCreateNewReq) (res model.PresetCreateNewRes, err error) {
	preset := entity.Preset{}
	preset.Id = ps.iSrv.GenSnowID()

	preset.PresetName = req.PresetName
	preset.PresetContent = req.PresetContent
	preset.ModelName = tools.DefaultValue(req.ModelName, openai.GPT3Dot5Turbo).(string)
	preset.MaxTokens = tools.DefaultValue(req.MaxTokens, 500).(int)
	preset.Temperature = tools.DefaultValue(req.Temperature, 0.7).(float64)
	preset.TopP = tools.DefaultValue(req.TopP, 1.0).(float64)
	preset.LogitBias = datatypes.JSON(tools.DefaultValue(req.LogitBias.String(), "").(string))
	preset.Frequency = tools.DefaultValue(req.Frequency, 0.1).(float64)
	preset.Presence = tools.DefaultValue(req.Presence, 0.2).(float64)
	preset.WithEmbedding = tools.DefaultValue(req.WithEmbedding, false).(bool)
	preset.Classify = tools.DefaultValue(req.Classify, "").(string)
	err = ps.pd.PresetCreateNew(ctx, &preset)
	if err != nil {
		res.IsSuccess = false
	} else {
		res.PresetId = preset.Id
		res.IsSuccess = true
	}
	return
}

func (ps *presetService) PresetGetList(ctx context.Context) (res model.PresetGetListRes, err error) {
	var presetOne model.PresetOneRes
	var presetlistRes []model.PresetOneRes
	presetlist, err := ps.pd.PresetGetList(ctx)
	if err != nil {
		return
	}
	for _, v := range presetlist {
		presetOne.PresetId = strconv.FormatInt(v.PresetId, 10)
		presetOne.PresetName = v.PresetName
		presetlistRes = append(presetlistRes, presetOne)
	}
	res.PresetsList = presetlistRes
	return
}
