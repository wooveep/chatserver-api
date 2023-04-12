/*
 * @Author: cloudyi.li
 * @Date: 2023-04-11 11:55:58
 * @LastEditTime: 2023-04-11 12:26:45
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/handler/v1/preset/preset.go
 */
package preset

import (
	"chatserver-api/internal/model"
	"chatserver-api/internal/service"
	"chatserver-api/pkg/errors"
	"chatserver-api/pkg/errors/ecode"
	"chatserver-api/pkg/response"
	"context"

	"github.com/gin-gonic/gin"
)

type PresetHandler struct {
	pSrv service.PresetService
}

func NewPresetHandler(_pSrv service.PresetService) *PresetHandler {

	ph := &PresetHandler{
		pSrv: _pSrv,
	}
	return ph
}

func (ph *PresetHandler) PresetCreateNew() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req model.PresetCreateNewReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			response.JSON(ctx, errors.WithCode(ecode.ValidateErr, err.Error()), nil)
			return
		}
		res, err := ph.pSrv.PresetCreateNew(context.TODO(), &req)
		if err != nil {
			response.JSON(ctx, errors.Wrap(err, ecode.Unknown, "创建失败"), nil)
		} else {
			response.JSON(ctx, nil, res)

		}

	}
}

func (ph *PresetHandler) PresetGetList() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		res, err := ph.pSrv.PresetGetList(context.TODO())
		if err != nil {
			response.JSON(ctx, errors.Wrap(err, ecode.Unknown, "查询失败"), nil)
		} else {
			response.JSON(ctx, nil, res)
		}
	}
}
