/*
 * @Author: cloudyi.li
 * @Date: 2023-03-29 11:51:00
 * @LastEditTime: 2023-03-29 16:37:57
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
	g.GET("/user/avatar-url", ar.userHandler.GetAvatar())

	g.POST("/chat/chatting", middleware.Stream(), ar.chatHandler.ChattingStreamSend())
}
