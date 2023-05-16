/*
 * @Author: cloudyi.li
 * @Date: 2023-05-15 13:30:31
 * @LastEditTime: 2023-05-15 16:57:29
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

	"github.com/gin-gonic/gin"
)

var _ AdminService = (*adminService)(nil)

type AdminService interface {
	AdminVerify(ctx *gin.Context) bool
	CdKeyGenerate(ctx *gin.Context, number, amount int) (res model.CdKeyGenerateRes, err error)
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

func (as *adminService) CdKeyGenerate(ctx *gin.Context, number, amount int) (res model.CdKeyGenerateRes, err error) {
	var cdkey entity.CdKey
	var cdkeylist []entity.CdKey
	var codekey []string
	for i := 0; i < number; i++ {
		keyId := as.aSrv.GenSnowID()
		code := uuid.IdToCode(keyId, 16, 32)
		cdkey.Id = keyId
		cdkey.CodeKey = code
		cdkey.Amount = amount
		cdkeylist = append(cdkeylist, cdkey)
		codekey = append(codekey, code)
	}
	err = as.kd.CdKeyGenerate(ctx, cdkeylist)
	res.CodeKey = codekey
	return
}
