/*
 * @Author: cloudyi.li
 * @Date: 2023-03-29 12:37:13
 * @LastEditTime: 2023-05-12 16:51:28
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/service/user.go
 */
package service

import (
	"chatserver-api/internal/consts"
	"chatserver-api/internal/dao"
	"chatserver-api/internal/model"
	"chatserver-api/internal/model/entity"
	"chatserver-api/pkg/active"
	"chatserver-api/pkg/avatar"
	"chatserver-api/pkg/config"
	"chatserver-api/pkg/jtime"
	"chatserver-api/pkg/jwt"
	"chatserver-api/pkg/logger"
	"chatserver-api/pkg/mail"
	"chatserver-api/utils/security"
	"chatserver-api/utils/uuid"
	"context"
	"encoding/base64"
	"errors"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var _ UserService = (*userService)(nil)

type UserService interface {
	UserLogout(ctx context.Context, tokenstr string) error
	UserGetByID(ctx context.Context, uid int64) (user entity.User, err error)
	UserRegister(ctx *gin.Context, req model.UserRegisterReq) (res model.UserRegisterRes, err error)
	UserGetAvatar(ctx context.Context, userId int64) (res model.UserGetAvatarRes, err error)
	UserLogin(ctx context.Context, username, password string) (res model.UserLoginRes, err error)
	UserRefresh(ctx *gin.Context) (res model.UserLoginRes, err error)
	UserGetInfo(ctx context.Context, userId int64) (res model.UserGetInfoRes, err error)
	UserVerifyEmail(ctx *gin.Context, email string) (res model.UserVerifyEmailRes, err error)
	UserVerifyUserName(ctx context.Context, username string) (res model.UserVerifyUserNameRes, err error)
	UserUpdateNickName(ctx context.Context, userId int64, nickname string) (res model.UserUpdateNickNameRes, err error)
	UserActiveGen(ctx *gin.Context) (err error)
	UserActiveChange(ctx *gin.Context) (err error)
	UserDelete(ctx *gin.Context) error
	UserPasswordVerify(ctx *gin.Context, password string) (Isvalid bool)
	UserPasswordModify(ctx *gin.Context, password string) (err error)
	UserPasswordForget(ctx *gin.Context) (err error)
	UserTempCodeVerify(ctx *gin.Context, tempcode string) (Isvalid bool)
	UserTempCodeGen(ctx *gin.Context) (tempcode string, email string, nickname string, err error)
}

// userService 实现UserService接口
type userService struct {
	ud   dao.UserDao
	kd   dao.CDkeyDao
	iSrv uuid.SnowNode
}

func NewUserService(_ud dao.UserDao, _kd dao.CDkeyDao) *userService {
	return &userService{
		ud:   _ud,
		kd:   _kd,
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

	if !security.ValidatePassword(security.PasswordDecryption(password, consts.CBCKEY), userInfo.Password) {
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

func (us *userService) UserRefresh(ctx *gin.Context) (res model.UserLoginRes, err error) {
	userId := ctx.GetInt64(consts.UserID)
	tokenStr := ctx.GetString(consts.TokenCtx)
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
	err = jwt.JoinBlackList(ctx, tokenStr, config.AppConfig.JwtConfig.Secret)
	if err != nil {
		logger.Infof("加入黑名单失败%s", userInfo.Username)
	}
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
	pattern := "^(http://|https://)"
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
		res.AvatarUrl = config.AppConfig.ExternalURL + user.AvatarUrl

	}
	return res, err
}

func (us *userService) UserGetInfo(ctx context.Context, userId int64) (res model.UserGetInfoRes, err error) {
	user, err := us.ud.UserGetById(ctx, userId)
	if err != nil {
		return res, err
	}
	pattern := "^(http://|https://)"
	match, err := regexp.MatchString(pattern, user.AvatarUrl)
	if err != nil {
		logger.Error(err.Error())
		return res, err
	}
	if match {
		res.AvatarUrl = user.AvatarUrl

	} else {
		res.AvatarUrl = config.AppConfig.ExternalURL + user.AvatarUrl

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
	ctx.Set(consts.UserID, user.Id)
	user.Username = req.Username
	user.Nickname = req.Username
	user.RegisteredIp = ctx.ClientIP()
	user.Email = req.Email
	user.ExpiredAt = jtime.JsonTime(time.Now().AddDate(0, 0, 7))
	user.Balance = 10
	user.IsActive = false
	user.AvatarUrl, err = avatar.GenNewAvatar(security.Md5WithSalt(req.Username, req.Email))
	if err != nil {
		return res, err
	}
	user.Password, err = security.Encrypt(security.PasswordDecryption(req.Password, consts.CBCKEY))
	if err != nil {
		return res, err
	}
	err = us.ud.UserCreate(ctx, &user)
	if err != nil {
		return res, err

	}
	res.IsSuccess = true
	return
}

func (us *userService) UserPasswordVerify(ctx *gin.Context, password string) (Isvalid bool) {
	userId := ctx.GetInt64(consts.UserID)
	Isvalid = false
	userInfo, err := us.ud.UserGetById(ctx, userId)
	if err != nil {
		logger.Infof("查询用户失败%s", err)
		return
	}
	if !security.ValidatePassword(security.PasswordDecryption(password, consts.CBCKEY), userInfo.Password) {
		// err = errors.New("Password Error")
		logger.Infof("密码错误%s", userInfo.Username)
		return
	}
	return true
}

func (us *userService) UserPasswordModify(ctx *gin.Context, password string) (err error) {
	user := entity.User{}
	userId := ctx.GetInt64(consts.UserID)
	user.Id = userId
	user.Password, err = security.Encrypt(security.PasswordDecryption(password, consts.CBCKEY))
	if err != nil {
		return err
	}
	err = us.ud.UserUpdate(ctx, &user)
	return err
}

func (us *userService) UserPasswordForget(ctx *gin.Context) (err error) {
	tempcode, email, nikcname, err := us.UserTempCodeGen(ctx)
	//19+16 35
	err = mail.SendForgetCode(email, nikcname, tempcode)
	return
}

func (us *userService) UserVerifyEmail(ctx *gin.Context, email string) (res model.UserVerifyEmailRes, err error) {
	UserId, err := us.ud.UserVerifyEmail(ctx, email)
	if err != nil {
		return
	}
	if UserId != 0 {
		res.Isvalid = false
		ctx.Set(consts.UserID, UserId)
	} else {
		res.Isvalid = true
	}
	logger.Debugf("邮箱校验信息：%d", UserId)
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

func (us *userService) UserLogout(ctx context.Context, tokenstr string) error {
	return jwt.JoinBlackList(ctx, tokenstr, config.AppConfig.JwtConfig.Secret)
}

func (us *userService) UserDelete(ctx *gin.Context) error {
	userId := ctx.GetInt64(consts.UserID)
	return us.ud.UserDelete(ctx, userId)
}

func (us *userService) UserTempCodeGen(ctx *gin.Context) (tempcode string, email string, nickname string, err error) {
	userId := ctx.GetInt64(consts.UserID)
	userInfo, err := us.ud.UserGetById(ctx, userId)
	code, err := active.ActiveCodeGen(ctx, userId)
	if err != nil {
		return

	}
	email = userInfo.Email
	nickname = userInfo.Nickname
	tempcode = base64.StdEncoding.EncodeToString([]byte(code + "|" + userInfo.Username))
	return
}

func (us *userService) UserTempCodeVerify(ctx *gin.Context, tempcode string) (Isvalid bool) {
	Isvalid = false
	codeStr, err := base64.StdEncoding.DecodeString(tempcode)
	if err != nil {
		return
	}
	codelist := strings.Split(string(codeStr), "|")
	if len(codelist) < 2 {
		// err = errors.New("Active Failed")
		return
	}
	code := codelist[0]
	username := codelist[1]
	userInfo, err := us.ud.UserGetByName(ctx, username)
	if err != nil {
		return
	}
	ctx.Set(consts.UserID, userInfo.Id)
	active := active.ActiveCodeCompare(ctx, code, userInfo.Id)
	if !active {
		Isvalid = false
	} else {
		Isvalid = true
	}
	return
}

func (us *userService) UserActiveGen(ctx *gin.Context) (err error) {
	tempcode, email, nikcname, err := us.UserTempCodeGen(ctx)
	//19+16 35
	err = mail.SendActiceCode(email, nikcname, tempcode)
	return
}

func (us *userService) UserActiveChange(ctx *gin.Context) (err error) {
	userId := ctx.GetInt64(consts.UserID)
	user := entity.User{}
	user.Id = userId
	user.IsActive = true
	err = us.ud.UserUpdate(ctx, &user)
	return err
}

func (us *userService) UserCDkeyPay(ctx *gin.Context, codekey string) {

}
