/*
 * @Author: cloudyi.li
 * @Date: 2023-06-16 21:16:47
 * @LastEditTime: 2023-06-16 21:30:42
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/chatfunc/chatfunc.go
 */
package chatfunc

import (
	"chatserver-api/pkg/logger"
	"context"
	"encoding/json"
)

func ChatFuncProcess(funcName string, arguments string) (content string) {
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(arguments), &result); err != nil {
		logger.Errorf("%v", err)
	}
	switch funcName {
	case "GetWeather":
		content = GetWeather(context.Background(), result["location"].(string))
	default:
		logger.Errorf("%v", "未匹配")

	}
	return content
}
