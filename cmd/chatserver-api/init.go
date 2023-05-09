/*
 * @Author: cloudyi.li
 * @Date: 2023-03-29 12:54:42
 * @LastEditTime: 2023-05-09 13:30:15
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/cmd/chatserver-api/init.go
 */
package chatserverapi

import (
	"chatserver-api/internal/dao/query"
	"chatserver-api/internal/handler/v1/chat"
	"chatserver-api/internal/handler/v1/preset"
	"chatserver-api/internal/handler/v1/user"
	"chatserver-api/internal/router"
	"chatserver-api/internal/service"
	"chatserver-api/pkg/db"
	"chatserver-api/pkg/tokenize"
)

func InitRouter(ds db.IDataSource, tk tokenize.Tokenizer) Router {
	userDao := query.NewUserDao(ds)
	userService := service.NewUserService(userDao)
	userhandler := user.NewUserHandler(userService)
	chatDao := query.NewChatDao(ds)
	chatService := service.NewChatService(chatDao, tk)
	chathandler := chat.NewChatHandler(chatService)
	presetDao := query.NewPresetsDao(ds)
	presetService := service.NewPresetService(presetDao)
	presetHandler := preset.NewPresetHandler(presetService)
	apiRouter := router.NewApiRouter(userhandler, chathandler, presetHandler)
	return apiRouter
}
