/*
 * @Author: cloudyi.li
 * @Date: 2023-04-12 13:16:26
 * @LastEditTime: 2023-04-20 16:41:57
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/model/record.go
 */
package model

import "chatserver-api/pkg/jtime"

type RecordOne struct {
	Id        int64          `gorm:"column:id" json:"record_id"`
	Sender    string         `gorm:"column:sender"  json:"sender"`
	Message   string         `gorm:"column:message"  json:"message" `
	CreatedAt jtime.JsonTime `gorm:"column:created_at"  json:"created_at" `
}

type RecordHistoryReq struct {
	ChatId string `form:"chat_id"  validate:"required"`
}

type RecordOneRes struct {
	Id        string         `json:"record_id"`
	Sender    string         `json:"sender"`
	Message   string         `json:"message" `
	CreatedAt jtime.JsonTime `json:"created_at" `
}

type RecordHistoryRes struct {
	ChatId  string         `json:"chat_id"`
	Records []RecordOneRes ` json:"record_list"`
}
