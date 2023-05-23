/*
 * @Author: cloudyi.li
 * @Date: 2023-03-29 11:51:00
 * @LastEditTime: 2023-05-22 10:30:23
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/router/api.go
 */
// Created on 2023/3/4.
// @author tony
// email xmgtony@gmail.com
// description

package router

import (
	"chatserver-api/internal/handler/v1/admin"
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
	adminHandler  *admin.AdminHandler
}

func NewApiRouter(
	userHandler *user.UserHandler,
	chatHandler *chat.ChatHandler,
	presetHandler *preset.PresetHandler,
	adminHandler *admin.AdminHandler,
) *ApiRouter {
	return &ApiRouter{
		userHandler:   userHandler,
		chatHandler:   chatHandler,
		presetHandler: presetHandler,
		adminHandler:  adminHandler,
	}
}

// Load 实现了server/http.go:40
func (ar *ApiRouter) Load(g *gin.Engine) {
	// login
	g.POST("/login", ar.userHandler.UserLogin())
	g.POST("/register", ar.userHandler.UserRegister())
	g.POST("/checkemail", ar.userHandler.UserVerifyEmail())
	g.POST("/checkusername", ar.userHandler.UserVerifyUserName())
	g.GET("/active", ar.userHandler.UserActive())
	g.POST("/forget", ar.userHandler.UserPasswordForget())
	g.POST("/resetpassword", ar.userHandler.UserPasswordReset())
	// g.GET("/test", ar.chatHandler.TestJieba())
	ug := g.Group("/user", middleware.AuthToken())
	{
		ug.GET("/avatar", ar.userHandler.UserGetAvatar())
		ug.GET("/info", ar.userHandler.UserGetInfo())
		ug.GET("/logout", ar.userHandler.UserLogout())
		ug.POST("/changenickname", ar.userHandler.UserUpdateNickName())
		ug.GET("/refresh", ar.userHandler.UserRefresh())
		ug.POST("/updatepassword", ar.userHandler.UserPasswordModify())
		ug.POST("/cdkeypay", ar.userHandler.UserCDkeyPay())
		ug.GET("/giftcard", ar.userHandler.UserGiftCardListGet())
		ug.GET("/invitelink", ar.userHandler.UserInviteLinkGet())

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
		cg.DELETE("/clear", ar.chatHandler.ChatRecordClear())
		cg.POST("/update", ar.chatHandler.ChatUpdate())
	}
	pg := g.Group("/preset", middleware.AuthToken())
	{
		pg.POST("/new", ar.presetHandler.PresetCreateNew())
		pg.GET("/list", ar.presetHandler.PresetGetList())
	}
	eg := g.Group("/embedding", middleware.AuthToken())
	{
		eg.POST("/file", ar.chatHandler.ChatEmbeddingFile())
		eg.POST("/string", ar.chatHandler.ChatEmbeddingString())
	}
	ag := g.Group("/admin", middleware.AuthToken())
	{
		ag.POST("/cdkeygen", ar.adminHandler.AdminGenNewCDkey())
		ag.POST("/cardcreate", ar.adminHandler.AdminCreateGiftCard())
		ag.POST("/cardupdate", ar.adminHandler.AdminUpdateGiftCard())
	}
}
