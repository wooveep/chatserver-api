/*
 * @Author: cloudyi.li
 * @Date: 2023-04-12 13:16:26
 * @LastEditTime: 2023-04-13 13:47:08
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
	Id int64 `form:"chat_id"  validate:"required"`
}

type RecordHistoryRes struct {
	Id      int64       `json:"chat_id"`
	Records []RecordOne ` json:"records"`
}
