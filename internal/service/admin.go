/*
 * @Author: cloudyi.li
 * @Date: 2023-05-15 13:30:31
 * @LastEditTime: 2023-05-21 19:38:30
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/service/admin.go
 */
package service

import (
	"chatserver-api/internal/consts"
	"chatserver-api/internal/dao"
	"chatserver-api/internal/model"
	"chatserver-api/internal/model/entity"
	"chatserver-api/utils/uuid"
	"strconv"

	"github.com/gin-gonic/gin"
)

var _ AdminService = (*adminService)(nil)

type AdminService interface {
	AdminVerify(ctx *gin.Context) bool
	CdKeyGenerate(ctx *gin.Context, number int, CardId int64) (res model.CdKeyGenerateRes, err error)
	GiftCardUpdate(ctx *gin.Context, req model.GiftCardUpdate) error
	GiftCardCreate(ctx *gin.Context, req model.GiftCardCreate) error
}

// userService 实现UserService接口
type adminService struct {
	kd   dao.CDkeyDao
	ud   dao.UserDao
	aSrv uuid.SnowNode
}

func NewAdminService(_kd dao.CDkeyDao, _ud dao.UserDao) *adminService {
	return &adminService{
		kd:   _kd,
		ud:   _ud,
		aSrv: *uuid.NewNode(5),
	}
}

func (as *adminService) AdminVerify(ctx *gin.Context) bool {
	userId := ctx.GetInt64(consts.UserID)
	role, err := as.ud.UserGetRole(ctx, userId)
	if err != nil || role != consts.Administrator {
		return false
	} else {
		return true
	}
}

func (as *adminService) CdKeyGenerate(ctx *gin.Context, number int, CardId int64) (res model.CdKeyGenerateRes, err error) {
	var cdkey entity.CdKey
	var cdkeylist []entity.CdKey
	var codekey []string
	for i := 0; i < number; i++ {
		keyId := as.aSrv.GenSnowID()
		code := uuid.IdToCode(keyId)
		cdkey.Id = keyId
		cdkey.CodeKey = code
		cdkey.GiftCardId = CardId
		cdkeylist = append(cdkeylist, cdkey)
		codekey = append(codekey, code)
	}
	err = as.kd.CdKeyGenerate(ctx, cdkeylist)
	res.CodeKey = codekey
	return
}

func (as *adminService) GiftCardCreate(ctx *gin.Context, req model.GiftCardCreate) error {
	var giftcard entity.GiftCard
	giftcard.Id = as.aSrv.GenSnowID()
	giftcard.CardAmount = req.CardAmount
	giftcard.CardDiscount = req.CardDiscount
	giftcard.CardName = req.CardName
	giftcard.CardBuyLink = req.CardLink
	giftcard.CardComment = req.CardComment
	return as.kd.GiftCardCreate(ctx, &giftcard)
}

func (as *adminService) GiftCardUpdate(ctx *gin.Context, req model.GiftCardUpdate) error {
	var giftcard entity.GiftCard
	cardId, err := strconv.ParseInt(req.CardId, 10, 64)
	if err != nil {
		return err
	}
	giftcard.Id = cardId
	if req.CardAmount != 0 {
		giftcard.CardAmount = req.CardAmount
	}
	if req.CardDiscount != 0 {
		giftcard.CardDiscount = req.CardDiscount
	}
	if req.CardName != "" {
		giftcard.CardName = req.CardName
	}
	if req.CardLink != "" {
		giftcard.CardBuyLink = req.CardLink
	}
	if req.CardComment != "" {
		giftcard.CardComment = req.CardComment
	}
	return as.kd.GiftCardUpdate(ctx, &giftcard)
}

// func (as *adminService) GiftCardListGet(ctx *gin.Context){

// }
