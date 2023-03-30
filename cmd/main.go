/*
 * @Author: cloudyi.li
 * @Date: 2023-02-15 11:26:59
 * @LastEditTime: 2023-03-29 12:56:38
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/cmd/main.go
 */
package main

import (
	chatserverapi "chatserver-api/cmd/chatserver-api"
	"chatserver-api/di/config"
	"chatserver-api/di/logger"
	"chatserver-api/internal/middleware"
)

func main() {
	var args chatserverapi.Args
	args = chatserverapi.LoadArgsValid()
	c := config.Load(args.Config)
	logger.InitLogger(&c.LogConfig, c.AppName)
	srv := chatserverapi.NewHttpServer(config.AppConfig)
	srvRouter := chatserverapi.InitRouter()
	srv.Run(middleware.NewMiddleware(), srvRouter)
}
