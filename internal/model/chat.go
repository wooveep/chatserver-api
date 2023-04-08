/*
 * @Author: cloudyi.li
 * @Date: 2023-03-29 14:11:49
 * @LastEditTime: 2023-04-06 12:16:55
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/model/chat.go
 */
package model

type ChatCreateNewReq struct {
	ChatName    string `json:"chatname" validate:"required" label:"会话名称"`
	PresetId    int64  `json:"presetid" validate:"required" label:"预设ID"`
	MemoryLevel int16  `json:"memorylevel" validate:"required" label:"消息记忆"`
}
type ChatCreateNewRes struct {
	ChatId int64 `json:"chatid"`
}

type ChatChattingReq struct {
	ChatId  string `json:"chatid" validate:"required" label:"会话ID"`
	Message string `json:"message" validate:"required" label:"消息"`
}

type ChatChattingRes struct {
}
