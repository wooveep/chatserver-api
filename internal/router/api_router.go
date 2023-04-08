/*
 * @Author: cloudyi.li
 * @Date: 2023-03-29 11:51:00
 * @LastEditTime: 2023-04-08 16:05:31
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/router/api_router.go
 */
// Created on 2023/3/4.
// @author tony
// email xmgtony@gmail.com
// description

package router

import (
	"chatserver-api/internal/handler/v1/chat"
	"chatserver-api/internal/handler/v1/user"
	"chatserver-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

type ApiRouter struct {
	userHandler *user.UserHandler
	chatHandler *chat.ChatHandler
}

func NewApiRouter(
	userHandler *user.UserHandler,
	chatHandler *chat.ChatHandler,
) *ApiRouter {
	return &ApiRouter{
		userHandler: userHandler,
		chatHandler: chatHandler,
	}
}

// Load 实现了server/http.go:40
func (ar *ApiRouter) Load(g *gin.Engine) {
	// login
	g.POST("/login", ar.userHandler.UserLogin())
	g.POST("/register", ar.userHandler.UserRegister())
	ug := g.Group("/user", middleware.AuthToken())
	{
		ug.GET("/avatar-url", ar.userHandler.UserGetAvatar())
		ug.GET("/info", ar.userHandler.UserGetInfo())
		ug.GET("/logout", ar.userHandler.UserLogout())

	}
	cg := g.Group("/chat", middleware.AuthToken())
	{
		cg.POST("/chatting", middleware.Stream(), ar.chatHandler.ChattingStreamSend())
		cg.POST("/new", ar.chatHandler.CreateNewChat())
	}
}
