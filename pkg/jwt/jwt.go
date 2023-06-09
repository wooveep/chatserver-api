/*
 * @Author: cloudyi.li
 * @Date: 2023-03-29 11:57:25
 * @LastEditTime: 2023-05-31 15:02:16
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/jwt/jwt.go
 */

package jwt

import (
	"chatserver-api/pkg/cache"
	"chatserver-api/pkg/config"
	"chatserver-api/pkg/logger"
	"chatserver-api/utils/security"
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
)

// CustomClaims 在标准声明中加入用户id
type CustomClaims struct {
	UserId int64 `json:"user_id"`
	RoleId int   `json:"role_id"`
	jwt.RegisteredClaims
}

func BuildClaims(exp time.Time, uid int64, rid int) *CustomClaims {
	return &CustomClaims{
		UserId: uid,
		RoleId: rid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    config.AppConfig.AppName,
		},
	}
}

// GenToken 生成jwt token
func GenToken(c *CustomClaims, secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	ss, err := token.SignedString([]byte(secretKey))
	return ss, err
}

// ParseToken 解析jwt token
func ParseToken(jwtStr, secretKey string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(jwtStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, err
	} else {
		return nil, err
	}
}
func getBlackListKey(token string) string {
	return "jwt_black_list:" + security.Md5(token)
}

func JoinBlackList(ctx context.Context, tokenstr string, secretKey string) (err error) {
	token, err := jwt.ParseWithClaims(tokenstr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return err
	}
	nowUnix := time.Now().Unix()
	timer := time.Duration(token.Claims.(*CustomClaims).ExpiresAt.Unix()-nowUnix) * time.Second
	rc := cache.GetRedisClient()
	err = rc.SetNX(ctx, getBlackListKey(token.Raw), nowUnix, timer).Err()
	return
}

func IsInBlackList(ctx context.Context, token string) bool {
	rc := cache.GetRedisClient()
	joinUnixStr, err := rc.Get(ctx, getBlackListKey(token)).Result()
	if err != nil {
		if err != redis.Nil {
			logger.Errorf("Redis连接异常:%v", err.Error())
		}
		return false
	}
	joinUnix, err := strconv.ParseInt(joinUnixStr, 10, 64)
	if time.Now().Unix()-joinUnix < config.AppConfig.JwtConfig.JwtBlacklistGracePeriod {
		return false
	}
	return true
}
