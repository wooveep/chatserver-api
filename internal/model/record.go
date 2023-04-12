/*
 * @Author: cloudyi.li
 * @Date: 2023-04-12 13:16:26
 * @LastEditTime: 2023-04-12 20:24:46
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/model/record.go
 */
package model

import "chatserver-api/pkg/jtime"

type RecordOne struct {
	Id        int64          `gorm:"column:Records__id" json:"id"`
	Sender    string         `gorm:"column:Records__sender"  json:"sender"`
	Message   string         `gorm:"column:Records__message"  json:"message" `
	CreatedAt jtime.JsonTime `gorm:"column:Records__create_at"  json:"created_at" `
}
