/*
 * @Author: cloudyi.li
 * @Date: 2023-03-29 10:51:03
 * @LastEditTime: 2023-04-08 15:35:37
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/consts/consts.go
 */
package consts

const (
	// RequestId 请求id名称
	RequestId = "request_id"
	// TimeLayout 时间格式
	TimeLayout   = "2006-01-02 15:04:05"
	TimeLayoutMs = "2006-01-02 15:04:05.000"
	// UserID 用户id key
	UserID                         = "user_id"
	TokenCtx                       = "token_ctx"
	ApiURLv1                       = "https://api.openai.com/v1"
	DefaultEmptyMessagesLimit uint = 300
	AvatarSize                     = 24
)
