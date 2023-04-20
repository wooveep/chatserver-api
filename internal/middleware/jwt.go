/*
 * @Author: cloudyi.li
 * @Date: 2023-03-29 11:55:27
 * @LastEditTime: 2023-04-20 15:05:29
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/middleware/jwt.go
 */
// Created on 2021/5/4.
// @author tony
// email xmgtony@gmail.com
// description jwt中间件

package middleware

import (
	"chatserver-api/internal/consts"
	"chatserver-api/pkg/config"
	"chatserver-api/pkg/errors"
	"chatserver-api/pkg/errors/ecode"
	"chatserver-api/pkg/jwt"
	"chatserver-api/pkg/response"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

// 请求头的形式为 Authorization: Bearer token
const authorizationHeader = "Authorization"

// AuthToken 鉴权，验证用户token是否有效
func AuthToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenstr, err := getJwtFromHeader(c)
		if err != nil {
			response.JSON(c, errors.Wrap(err, ecode.RequireAuthErr, "invalid token"), nil)
			c.Abort()
			return
		}
		if jwt.IsInBlackList(tokenstr) {
			response.JSON(c, errors.WithCode(ecode.RequireAuthErr, "invalid token"), nil)
			c.Abort()
			return
		}
		// 验证token是否正确
		claims, err := jwt.ParseToken(tokenstr, config.AppConfig.JwtConfig.Secret)
		if err != nil {
			response.JSON(c, errors.Wrap(err, ecode.RequireAuthErr, "invalid token"), nil)
			c.Abort()
			return
		}
		c.Set(consts.UserID, claims.UserId)
		c.Set(consts.TokenCtx, tokenstr)
		c.Next()
	}
}

func getJwtFromHeader(c *gin.Context) (string, error) {
	aHeader := c.Request.Header.Get(authorizationHeader)
	if len(aHeader) == 0 {
		return "", fmt.Errorf("token is empty")
	}
	strs := strings.SplitN(aHeader, " ", 2)
	if len(strs) != 2 || strs[0] != "Bearer" {
		return "", fmt.Errorf("token 不符合规则")
	}
	return strs[1], nil
}
