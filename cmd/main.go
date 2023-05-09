/*
 * @Author: cloudyi.li
 * @Date: 2023-02-15 11:26:59
 * @LastEditTime: 2023-05-09 13:27:09
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/cmd/main.go
 */
package main

import (
	chatserverapi "chatserver-api/cmd/chatserver-api"
	"chatserver-api/pkg/cache"
	"chatserver-api/pkg/config"
	"chatserver-api/pkg/db"
	"chatserver-api/pkg/logger"
	"chatserver-api/pkg/tokenize"

	"chatserver-api/internal/middleware"
)

func main() {
	var args chatserverapi.Args
	args = chatserverapi.LoadArgsValid()
	c := config.Load(args.Config)
	logger.InitLogger(&c.LogConfig, c.AppName)
	ds := db.NewDefaultPostGre(c.DBConfig)
	cache.InitRedis(c.RedisConfig)
	tk := tokenize.NewTokenizer()

	srv := chatserverapi.NewHttpServer(config.AppConfig)
	srv.RegisterOnShutdown(func() {
		if ds != nil {
			ds.Close()
		}
		cache.CloseRedis()
	})
	srvRouter := chatserverapi.InitRouter(ds, tk)
	srv.Run(middleware.NewMiddleware(), srvRouter)
}
