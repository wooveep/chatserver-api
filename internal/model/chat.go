/*
 * @Author: cloudyi.li
 * @Date: 2023-03-29 14:11:49
 * @LastEditTime: 2023-03-31 16:50:45
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/model/chat.go
 */
package model

type ChatChattingReq struct {
	ChatId   string `json:"chatid"`
	Qmessage string `json:"qmessage"`
}

type ChatChattingRes struct {
}
