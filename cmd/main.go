/*
 * @Author: cloudyi.li
 * @Date: 2023-02-15 11:26:59
 * @LastEditTime: 2023-03-28 14:34:25
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/cmd/main.go
 */
package main

import (
	chatserverapi "chatserver-api/cmd/chatserver-api"
	"chatserver-api/di/config"
	"chatserver-api/di/logger"
)

func main() {
	var args chatserverapi.Args
	args = chatserverapi.LoadArgsValid()
	config.InitConfig(args.Config)
	logger.InitLogger(config.AppConfig.Log)
	chatserverapi.Run(args)
}
