/*
 * @Author: cloudyi.li
 * @Date: 2023-03-29 12:37:13
 * @LastEditTime: 2023-04-21 10:23:50
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
	"os"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
)

var _ UserService = (*userService)(nil)

type UserService interface {
	UserLogout(tokenstr string) error
	UserGetByID(ctx context.Context, uid int64) (user entity.User, err error)
	UserRegister(ctx *gin.Context, req model.UserRegisterReq) (res model.UserRegisterRes, err error)
	UserGetAvatar(ctx context.Context, userId int64) (res model.UserGetAvatarRes, err error)
	UserLogin(ctx context.Context, username, password string) (res model.UserLoginRes, err error)
	UserRefresh(ctx context.Context, userId int64) (res model.UserLoginRes, err error)
	UserGetInfo(ctx context.Context, userId int64) (res model.UserGetInfoRes, err error)
	UserVerifyEmail(ctx context.Context, email string) (res model.UserVerifyEmailRes, err error)
	UserVerifyUserName(ctx context.Context, username string) (res model.UserVerifyUserNameRes, err error)
	UserUpdateNickName(ctx context.Context, userId int64, nickname string) (res model.UserUpdateNickNameRes, err error)
}

// userService 实现UserService接口
type userService struct {
	ud   dao.UserDao
	iSrv uuid.SnowNode
}

func NewUserService(_ud dao.UserDao) *userService {
	return &userService{
		ud:   _ud,
		iSrv: *uuid.NewNode(3),
	}
}

func (us *userService) UserLogin(ctx context.Context, username, password string) (res model.UserLoginRes, err error) {
	userInfo, err := us.ud.UserGetByName(ctx, username)
	if err != nil {
		logger.Infof("查询用户失败%s", err)
		return res, err
	}
	if userInfo.IsActive != true {
		err = errors.New("用户未激活")
		return res, err
	}
	if !security.ValidatePassword(password, userInfo.Password) {
		err = errors.New("Password Error")
		logger.Infof("密码错误%s", username)
		return res, err
	}
	if userInfo.ExpiredAt.GetUnixTime() < time.Now().Unix() {
		err = errors.New("用户授权过期")
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

func (us *userService) UserRefresh(ctx context.Context, userId int64) (res model.UserLoginRes, err error) {
	userInfo, err := us.ud.UserGetById(ctx, userId)
	if err != nil {
		logger.Infof("查询用户失败%s", err)
		return res, err
	}
	if userInfo.IsActive != true {
		err = errors.New("用户未激活")
		return res, err
	}
	if userInfo.ExpiredAt.GetUnixTime() < time.Now().Unix() {
		err = errors.New("用户授权过期")
		return res, err
	}
	expireAt := time.Now().Add(time.Duration(config.AppConfig.JwtConfig.JwtTtl) * time.Second)
	claims := jwt.BuildClaims(expireAt, userInfo.Id)
	token, err := jwt.GenToken(claims, config.AppConfig.JwtConfig.Secret)
	if err != nil {
		logger.Infof("JWTTOKEN生成错误%s", userInfo.Username)

		return res, err
	}
	res.Token = token
	res.ExpireAt = jtime.JsonTime(expireAt)
	return res, err
}

// GetByName 通过用户名 查找用户
func (us *userService) UserGetByID(ctx context.Context, uid int64) (user entity.User, err error) {
	return us.ud.UserGetById(ctx, uid)
}

func (us *userService) UserGetAvatar(ctx context.Context, userId int64) (res model.UserGetAvatarRes, err error) {
	user, err := us.ud.UserGetById(ctx, userId)
	if err != nil {
		return res, err
	}
	logger.Debugf("获取用户头像UID:%d,%s", userId, user.AvatarUrl)
	pattern := "^http://.*"
	match, err := regexp.MatchString(pattern, user.AvatarUrl)
	if err != nil {
		logger.Error(err.Error())
		return res, err
	}
	if match {
		res.AvatarUrl = user.AvatarUrl

	} else {
		if _, err := os.Stat(user.AvatarUrl); os.IsNotExist(err) {
			return res, err
		}
		res.AvatarUrl = config.AppConfig.AvatarURL + user.AvatarUrl

	}
	return res, err
}

func (us *userService) UserGetInfo(ctx context.Context, userId int64) (res model.UserGetInfoRes, err error) {
	user, err := us.ud.UserGetById(ctx, userId)
	if err != nil {
		return res, err
	}
	pattern := "^http://.*"
	match, err := regexp.MatchString(pattern, user.AvatarUrl)
	if err != nil {
		logger.Error(err.Error())
		return res, err
	}
	if match {
		res.AvatarUrl = user.AvatarUrl

	} else {
		res.AvatarUrl = config.AppConfig.AvatarURL + user.AvatarUrl

	}
	res.Balance = user.Balance
	res.Email = user.Email
	res.Nickname = user.Nickname
	res.Username = user.Username
	res.Phone = user.Phone
	res.ExpiredAt = jtime.JsonTime(user.ExpiredAt)
	return res, err
}

func (us *userService) UserRegister(ctx *gin.Context, req model.UserRegisterReq) (res model.UserRegisterRes, err error) {
	user := entity.User{}
	res.IsSuccess = false
	user.Id = us.iSrv.GenSnowID()
	user.Username = req.Username
	user.Nickname = req.Username
	user.RegisteredIp = ctx.ClientIP()
	user.Email = req.Email
	user.ExpiredAt = jtime.JsonTime(time.Now().AddDate(0, 1, 0))
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

func (us *userService) UserVerifyEmail(ctx context.Context, email string) (res model.UserVerifyEmailRes, err error) {
	count, err := us.ud.UserVerifyEmail(ctx, email)
	if err != nil {
		return
	}
	if count != 0 {
		res.Isvalid = false
	} else {
		res.Isvalid = true
	}
	return
}

func (us *userService) UserVerifyUserName(ctx context.Context, username string) (res model.UserVerifyUserNameRes, err error) {
	count, err := us.ud.UserVerifyUserName(ctx, username)
	if err != nil {
		return
	}
	if count != 0 {
		res.Isvalid = false
	} else {
		res.Isvalid = true
	}
	return
}

func (us *userService) UserUpdateNickName(ctx context.Context, userId int64, nickname string) (res model.UserUpdateNickNameRes, err error) {
	err = us.ud.UserUpdateNickName(ctx, userId, nickname)
	if err != nil {
		res.IsChanged = false
	} else {
		res.IsChanged = true
	}
	return
}

func (us *userService) UserLogout(tokenstr string) error {
	return jwt.JoinBlackList(tokenstr, config.AppConfig.JwtConfig.Secret)
}
