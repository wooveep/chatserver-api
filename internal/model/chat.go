/*
 * @Author: cloudyi.li
 * @Date: 2023-03-29 14:11:49
 * @LastEditTime: 2023-04-05 15:24:17
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/model/chat.go
 */
package model

type ChatCreateNewReq struct {
	ChatName    string `json:"chatname" validate:"required" label:"会话名称"`
	PresetId    string `json:"presetid" validate:"required" label:"预设ID"`
	MemoryLevel string `json:"memorylevel" validate:"required" label:"消息记忆"`
}
type ChatCreateNewRes struct {
	ChatId string `json:"chatid"`
}

type ChatChattingReq struct {
	ChatId  string `json:"chatid" validate:"required" label:"会话ID"`
	Message string `json:"message" validate:"required" label:"消息"`
}

type ChatChattingRes struct {
}
