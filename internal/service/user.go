/*
 * @Author: cloudyi.li
 * @Date: 2023-03-29 12:37:13
 * @LastEditTime: 2023-04-08 16:02:52
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/service/user.go
 */
package service

import (
	"chatserver-api/internal/dao"
	"chatserver-api/internal/model"
	"chatserver-api/internal/model/entity"
	"chatserver-api/pkg/avatar"
	"chatserver-api/pkg/config"
	"chatserver-api/pkg/jtime"
	"chatserver-api/pkg/jwt"
	"chatserver-api/pkg/logger"
	"chatserver-api/utils/security"
	"chatserver-api/utils/uuid"
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
)

var _ UserService = (*userService)(nil)

type UserService interface {
	UserLogout(tokenstr string) error
	UserGetByID(ctx context.Context, uid int64) (user *entity.User, err error)
	UserRegister(ctx *gin.Context, req model.UserRegisterReq) (res model.UserRegisterRes, err error)
	UserGetAvatar(ctx context.Context, userid int64) (res model.UserGetAvatarRes, err error)
	UserLogin(ctx context.Context, username, password string) (res model.UserLoginRes, err error)
	UserGetInfo(ctx context.Context, userid int64) (res model.UserGetInfoRes, err error)
}

// userService 实现UserService接口
type userService struct {
	ud dao.UserDao
}

func NewUserService(_ud dao.UserDao) *userService {
	return &userService{
		ud: _ud,
	}
}

func (us *userService) UserLogin(ctx context.Context, username, password string) (res model.UserLoginRes, err error) {
	userInfo, err := us.ud.UserGetByName(ctx, username)
	if err != nil {
		logger.Infof("查询用户失败%s", err)
		return res, err
	}
	if !security.ValidatePassword(password, userInfo.Password) {
		err = errors.New("Password Error")
		logger.Infof("密码错误%s", username)
		return res, err
	}
	expireAt := time.Now().Add(time.Duration(config.AppConfig.JwtConfig.JwtTtl) * time.Second)
	claims := jwt.BuildClaims(expireAt, userInfo.Id)
	token, err := jwt.GenToken(claims, config.AppConfig.JwtConfig.Secret)
	if err != nil {
		logger.Infof("JWTTOKEN生成错误%s", username)

		return res, err
	}
	res.Token = token
	res.ExpireAt = jtime.JsonTime(expireAt)
	return res, err
}

// GetByName 通过用户名 查找用户
func (us *userService) UserGetByID(ctx context.Context, uid int64) (user *entity.User, err error) {
	return us.ud.UserGetById(ctx, uid)
}

func (us *userService) UserGetAvatar(ctx context.Context, userid int64) (res model.UserGetAvatarRes, err error) {
	user, err := us.ud.UserGetById(ctx, userid)
	if err != nil {
		return res, err
	}
	logger.Debugf("获取用户头像UID:%d,%s", userid, user.AvatarUrl)
	res.AvatarUrl = user.AvatarUrl
	return res, err
}

func (us *userService) UserGetInfo(ctx context.Context, userid int64) (res model.UserGetInfoRes, err error) {
	user, err := us.ud.UserGetById(ctx, userid)
	if err != nil {
		return res, err
	}
	res.AvatarUrl = user.AvatarUrl
	res.Balance = user.Balance
	res.Email = user.Email
	res.Nickname = user.Nickname
	res.Username = user.Username
	res.Phone = user.Phone

	return res, err
}
func (us *userService) UserRegister(ctx *gin.Context, req model.UserRegisterReq) (res model.UserRegisterRes, err error) {
	user := entity.User{}
	res.IsSuccess = false
	user.Id, err = uuid.GenID()
	if err != nil {
		return res, err
	}
	user.Username = req.Username
	user.Nickname = req.Username
	user.RegisteredIp = ctx.ClientIP()
	user.Email = req.Email
	user.AvatarUrl, err = avatar.GenNewAvatar(security.Md5WithSalt(req.Username, req.Email))
	if err != nil {
		return res, err
	}
	user.Password, err = security.Encrypt(req.Password)
	if err != nil {
		return res, err
	}
	err = us.ud.UserRegisterNew(ctx, &user)
	if err != nil {
		return res, err

	}
	res.IsSuccess = true
	return
}

func (us *userService) UserVerifyUserName() {

}
func (us *userService) UserVerifyEmail() {

}
func (us *userService) UserLogout(tokenstr string) error {
	return jwt.JoinBlackList(tokenstr, config.AppConfig.JwtConfig.Secret)

}
