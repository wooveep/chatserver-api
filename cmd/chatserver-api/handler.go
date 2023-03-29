/*
 * @Author: cloudyi.li
 * @Date: 2023-02-15 17:24:53
 * @LastEditTime: 2023-03-29 09:33:12
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/cmd/chatserver-api/handler.go
 */
package chatserverapi

import (
	"chatserver-api/di/logger"
	"fmt"
)

func Run(args Args) {
	logger.Warn("程序启动")
	if args.Names == "openai" {
		fmt.Print(args.Names)
	} else {
		fmt.Print(args.Names)
	}
}
