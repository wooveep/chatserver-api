/*
 * @Author: cloudyi.li
 * @Date: 2023-03-29 11:51:00
 * @LastEditTime: 2023-04-21 10:31:24
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
	"chatserver-api/internal/handler/v1/preset"
	"chatserver-api/internal/handler/v1/user"
	"chatserver-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

type ApiRouter struct {
	userHandler   *user.UserHandler
	chatHandler   *chat.ChatHandler
	presetHandler *preset.PresetHandler
}

func NewApiRouter(
	userHandler *user.UserHandler,
	chatHandler *chat.ChatHandler,
	presetHandler *preset.PresetHandler,
) *ApiRouter {
	return &ApiRouter{
		userHandler:   userHandler,
		chatHandler:   chatHandler,
		presetHandler: presetHandler,
	}
}

// Load 实现了server/http.go:40
func (ar *ApiRouter) Load(g *gin.Engine) {
	// login
	g.POST("/login", ar.userHandler.UserLogin())
	g.POST("/register", ar.userHandler.UserRegister())
	g.POST("/checkemail", ar.userHandler.UserVerifyEmail())
	g.POST("/checkusername", ar.userHandler.UserVerifyUserName())
	ug := g.Group("/user", middleware.AuthToken())
	{
		ug.GET("/avatar-url", ar.userHandler.UserGetAvatar())
		ug.GET("/info", ar.userHandler.UserGetInfo())
		ug.GET("/logout", ar.userHandler.UserLogout())
		ug.POST("/changenickname", ar.userHandler.UserUpdateNickName())
		ug.GET("/refresh", ar.userHandler.UserRefresh())

	}
	cg := g.Group("/chat", middleware.AuthToken())
	{
		cg.POST("/chatting", middleware.Stream(), ar.chatHandler.ChatChatting())
		cg.POST("/regenerate", middleware.Stream(), ar.chatHandler.ChatRegenerateg())
		cg.POST("/new", ar.chatHandler.ChatCreateNew())
		cg.GET("/list", ar.chatHandler.ChatListGet())
		cg.POST("/detail", ar.chatHandler.ChatDetailGet())
		cg.GET("/history", ar.chatHandler.ChatRecordHistory())
		cg.DELETE("/delete", ar.chatHandler.ChatDelete())
		cg.POST("/update", ar.chatHandler.ChatUpdate())
	}
	pg := g.Group("/preset", middleware.AuthToken())
	{
		pg.POST("/new", ar.presetHandler.PresetCreateNew())
		pg.GET("/list", ar.presetHandler.PresetGetList())
	}
}
