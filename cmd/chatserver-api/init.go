package chatserverapi

import (
	"chatserver-api/internal/handler/v1/chat"
	"chatserver-api/internal/handler/v1/user"
	"chatserver-api/internal/router"
	"chatserver-api/internal/service"
)

func InitRouter() Router {

	userService := service.NewUserService()
	userhandler := user.NewUserHandler(userService)
	chatService := service.NewChatService()
	chathandler := chat.NewChatHandler(chatService)
	apiRouter := router.NewApiRouter(userhandler, chathandler)
	return apiRouter

}
