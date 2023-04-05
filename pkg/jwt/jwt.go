/*
 * @Author: cloudyi.li
 * @Date: 2023-03-29 11:57:25
 * @LastEditTime: 2023-04-05 15:54:24
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/jwt/jwt.go
 */

package jwt

import (
	"chatserver-api/pkg/config"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// CustomClaims 在标准声明中加入用户id
type CustomClaims struct {
	UserId int64 `json:"user_id"`
	jwt.RegisteredClaims
}

func BuildClaims(exp time.Time, uid int64) *CustomClaims {
	return &CustomClaims{
		UserId: uid,
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
