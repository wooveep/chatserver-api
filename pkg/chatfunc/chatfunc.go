/*
 * @Author: cloudyi.li
 * @Date: 2023-06-16 21:16:47
 * @LastEditTime: 2023-06-21 21:54:57
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/chatfunc/chatfunc.go
 */
package chatfunc

import (
	"chatserver-api/pkg/logger"
	"context"
	"encoding/json"
)

func ChatFuncProcess(ctx context.Context, funcName string, arguments string) (content string) {
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(arguments), &result); err != nil {
		logger.Errorf("%v", err)
	}
	switch funcName {
	case "GetWeather":
		content = GetWeather(ctx, result["location"].(string))
	case "EntitySearch":
		content = EntitySearch(ctx, result["query"].(string), result["etype"].(string))
	case "GoogleSearch":
		content = GoogleSearch(ctx, result["query"].(string))
	default:
		logger.Errorf("%v", "未匹配")
	}
	return content
}
