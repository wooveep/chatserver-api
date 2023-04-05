/*
 * @Author: cloudyi.li
 * @Date: 2023-03-29 12:54:42
 * @LastEditTime: 2023-04-05 15:53:21
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/cmd/chatserver-api/init.go
 */
package chatserverapi

import (
	"chatserver-api/internal/dao/query"
	"chatserver-api/internal/handler/v1/chat"
	"chatserver-api/internal/handler/v1/user"
	"chatserver-api/internal/router"
	"chatserver-api/internal/service"
	"chatserver-api/pkg/db"
)

func InitRouter(ds db.IDataSource) Router {
	userDao := query.NewUserDao(ds)
	userService := service.NewUserService(userDao)
	userhandler := user.NewUserHandler(userService)
	chatDao := query.NewChatDao(ds)
	chatService := service.NewChatService(chatDao)
	chathandler := chat.NewChatHandler(chatService)
	apiRouter := router.NewApiRouter(userhandler, chathandler)
	return apiRouter

}
