/*
 * @Author: cloudyi.li
 * @Date: 2023-03-29 10:51:03
 * @LastEditTime: 2023-04-23 14:00:33
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/consts/consts.go
 */
package consts

const (
	// RequestId 请求id名称
	RequestId = "request_id"
	// UserID 用户id key
	UserID   = "user_id"
	ChatID   = "chat_id"
	Balance  = "balance_ctx"
	TokenCtx = "token_ctx"
	// TimeLayout 时间格式
	TimeLayout                     = "2006-01-02 15:04:05"
	TimeLayoutMs                   = "2006-01-02 15:04:05.000"
	ApiURLv1                       = "https://api.openai.com/v1"
	DefaultEmptyMessagesLimit uint = 300
	AvatarSize                     = 24
	TokenPrice                     = 0.00007
)
