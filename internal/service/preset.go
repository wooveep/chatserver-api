/*
 * @Author: cloudyi.li
 * @Date: 2023-04-11 10:22:31
 * @LastEditTime: 2023-05-27 10:25:30
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/service/preset.go
 */
package service

import (
	"chatserver-api/internal/consts"
	"chatserver-api/internal/dao"
	"chatserver-api/internal/model"
	"chatserver-api/internal/model/entity"
	"chatserver-api/pkg/cache"
	"chatserver-api/pkg/logger"
	"chatserver-api/pkg/openai"
	"chatserver-api/utils/tools"
	"chatserver-api/utils/uuid"
	"context"
	"encoding/json"
	"strconv"

	"github.com/go-redis/redis/v8"
	"gorm.io/datatypes"
)

var _ PresetService = (*presetService)(nil)

type PresetService interface {
	PresetCreateNew(ctx context.Context, req model.PresetCreateNewReq) (res model.PresetCreateNewRes, err error)
	PresetGetList(ctx context.Context) (res model.PresetGetListRes, err error)
	PresetUpdate(ctx context.Context, req model.PresetUpdateReq) (err error)
}

// userService 实现UserService接口
type presetService struct {
	pd   dao.PresetDao
	iSrv uuid.SnowNode
	rc   *redis.Client
}

func NewPresetService(_pd dao.PresetDao) *presetService {
	return &presetService{
		pd:   _pd,
		iSrv: *uuid.NewNode(2),
		rc:   cache.GetRedisClient(),
	}
}

func (ps *presetService) PresetUpdate(ctx context.Context, req model.PresetUpdateReq) (err error) {
	err = ps.rc.Del(ctx, consts.PresetPrefix+"1").Err()
	if err != nil {
		logger.Errorf("删除PresetList缓存失败:%v", err.Error())
	}
	preset := entity.Preset{}
	presetId, err := strconv.ParseInt(req.PresetId, 10, 64)
	if err != nil {
		return
	}
	preset.Id = presetId
	preset.PresetName = req.PresetName
	preset.PresetContent = req.PresetContent
	preset.PresetTips = req.PresetTips
	preset.Classify = req.Classify
	preset.Frequency = req.Frequency
	preset.LogitBias = req.LogitBias
	preset.WithEmbedding = req.WithEmbedding
	preset.Temperature = req.Temperature
	preset.ModelName = req.ModelName
	preset.MaxTokens = req.MaxTokens
	preset.TopP = req.TopP
	preset.Presence = req.Presence
	err = ps.pd.PresetUpdate(ctx, &preset)
	if err != nil {
		return
	}
	return nil
}
func (ps *presetService) PresetCreateNew(ctx context.Context, req model.PresetCreateNewReq) (res model.PresetCreateNewRes, err error) {
	err = ps.rc.Del(ctx, consts.PresetPrefix+"1").Err()
	if err != nil {
		logger.Errorf("删除PresetList缓存失败:%v", err.Error())
	}
	preset := entity.Preset{}
	preset.Id = ps.iSrv.GenSnowID()
	preset.PresetName = req.PresetName
	preset.PresetContent = req.PresetContent
	preset.PresetTips = req.PresetTips
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
		return res, err
	} else {
		res.PresetId = preset.Id
		res.IsSuccess = true
	}
	return res, nil
}

func (ps *presetService) PresetGetList(ctx context.Context) (res model.PresetGetListRes, err error) {
	jsonbyte, err := ps.rc.Get(ctx, consts.PresetPrefix+"1").Bytes()
	if err == nil {
		err = json.Unmarshal(jsonbyte, &res)
		if err == nil {
			return res, nil
		}
		logger.Errorf(" PresetListRes反序列化失败:%v", err.Error())
	} else {
		if err != redis.Nil {
			logger.Errorf("Redis连接异常:%v", err.Error())
		}
		logger.Debugf(" PresetList缓存不存在:%v", err.Error())
	}
	var presetOne model.PresetOneRes
	var presetlistRes []model.PresetOneRes
	presetlist, err := ps.pd.PresetGetList(ctx)
	if err != nil {
		return
	}
	for _, v := range presetlist {
		presetOne.PresetId = strconv.FormatInt(v.PresetId, 10)
		presetOne.PresetName = v.PresetName
		presetOne.PresetTips = v.PresetTips
		presetlistRes = append(presetlistRes, presetOne)
	}
	res.PresetsList = presetlistRes
	jsonbyte, err = json.Marshal(res)
	if err != nil {
		logger.Errorf("PresetListRes序列化失败:%v", err.Error())
		return res, nil
	}
	err = ps.rc.Set(ctx, consts.PresetPrefix+"1", jsonbyte, 0).Err()
	if err != nil {
		logger.Errorf("PresetListRes存储Cache失败:%v", err.Error())
		return res, nil
	}
	return
}
