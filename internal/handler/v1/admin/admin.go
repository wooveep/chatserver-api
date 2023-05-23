/*
 * @Author: cloudyi.li
 * @Date: 2023-05-12 22:43:23
 * @LastEditTime: 2023-05-22 10:31:09
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/handler/v1/admin/admin.go
 */
package admin

import (
	"chatserver-api/internal/model"
	"chatserver-api/internal/service"
	"chatserver-api/pkg/errors"
	"chatserver-api/pkg/errors/ecode"
	"chatserver-api/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	aSrv service.AdminService
}

func NewAdminHandler(_aSrv service.AdminService) *AdminHandler {

	ah := &AdminHandler{
		aSrv: _aSrv,
	}
	return ah
}

func (ah *AdminHandler) AdminGenNewCDkey() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := model.CdKeyGenerateReq{}
		if err := ctx.ShouldBindJSON(&req); err != nil {
			response.JSON(ctx, errors.WithCode(ecode.ValidateErr, err.Error()), nil)
			return
		}
		if !ah.aSrv.AdminVerify(ctx) {
			response.JSON(ctx, errors.WithCode(ecode.PermissionErr, "权限错误"), nil)
			return
		}
		cardId, err := strconv.ParseInt(req.GiftCardId, 10, 64)
		if err != nil {
			response.JSON(ctx, errors.WithCode(ecode.ValidateErr, err.Error()), nil)
			return
		}
		res, err := ah.aSrv.CdKeyGenerate(ctx, req.KeyNumber, cardId)
		if err != nil {
			response.JSON(ctx, errors.WithCode(ecode.CreatErr, "错误"), nil)
			return
		}
		response.JSON(ctx, nil, res)
	}
}

func (ah *AdminHandler) AdminCreateGiftCard() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := model.GiftCardCreate{}
		if err := ctx.ShouldBindJSON(&req); err != nil {
			response.JSON(ctx, errors.WithCode(ecode.ValidateErr, err.Error()), nil)
			return
		}
		if !ah.aSrv.AdminVerify(ctx) {
			response.JSON(ctx, errors.WithCode(ecode.PermissionErr, "权限错误"), nil)
			return
		}
		if err := ah.aSrv.GiftCardCreate(ctx, req); err != nil {
			response.JSON(ctx, errors.WithCode(ecode.CreatErr, "错误"), nil)
			return
		}
		response.JSON(ctx, nil, nil)
	}
}

func (ah *AdminHandler) AdminUpdateGiftCard() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := model.GiftCardUpdate{}
		if err := ctx.ShouldBindJSON(&req); err != nil {
			response.JSON(ctx, errors.WithCode(ecode.ValidateErr, err.Error()), nil)
			return
		}
		if !ah.aSrv.AdminVerify(ctx) {
			response.JSON(ctx, errors.WithCode(ecode.PermissionErr, "权限错误"), nil)
			return
		}
		if err := ah.aSrv.GiftCardUpdate(ctx, req); err != nil {
			response.JSON(ctx, errors.WithCode(ecode.CreatErr, "错误"), nil)
			return
		}
		response.JSON(ctx, nil, nil)
	}
}
