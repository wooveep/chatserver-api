/*
 * @Author: cloudyi.li
 * @Date: 2023-03-29 12:54:42
 * @LastEditTime: 2023-05-19 13:59:08
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/cmd/chatserver-api/init.go
 */
package chatserverapi

import (
	"chatserver-api/internal/dao/query"
	"chatserver-api/internal/handler/v1/admin"
	"chatserver-api/internal/handler/v1/chat"
	"chatserver-api/internal/handler/v1/preset"
	"chatserver-api/internal/handler/v1/user"
	"chatserver-api/internal/router"
	"chatserver-api/internal/service"
	"chatserver-api/pkg/db"
)

func InitRouter(ds db.IDataSource) Router {
	// tk := tokenize.NewTokenizer()
	userDao := query.NewUserDao(ds)
	cdkeyDao := query.NewCDkeyDao(ds)
	userService := service.NewUserService(userDao, cdkeyDao)
	userhandler := user.NewUserHandler(userService)
	chatDao := query.NewChatDao(ds)
	chatService := service.NewChatService(chatDao)
	chathandler := chat.NewChatHandler(chatService)
	presetDao := query.NewPresetsDao(ds)
	presetService := service.NewPresetService(presetDao)
	presetHandler := preset.NewPresetHandler(presetService)
	adminService := service.NewAdminService(cdkeyDao, userDao)
	adminHandler := admin.NewAdminHandler(adminService)
	apiRouter := router.NewApiRouter(userhandler, chathandler, presetHandler, adminHandler)
	return apiRouter
}
