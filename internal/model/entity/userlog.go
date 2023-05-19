/*
 * @Author: cloudyi.li
 * @Date: 2023-05-18 13:04:50
 * @LastEditTime: 2023-05-18 13:06:15
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/model/entity/userlog.go
 */
package entity

import (
	"chatserver-api/pkg/jtime"
)

type UserLog struct {
	Id        int64          `gorm:"column:id;primary_key;" json:"id"`
	UserId    int64          `gorm:"column:user_id" json:"user_id"`
	UserIP    string         `gorm:"column:user_ip" json:"user_ip"`
	Business  string         `gorm:"column:business" json:"business"`
	Operation string         `gorm:"column:operation" json:"operation"`
	CreatedAt jtime.JsonTime `gorm:"column:created_at" json:"created_at"`
	UpdatedAt jtime.JsonTime `gorm:"column:updated_at" json:"updated_at"`
}

func (UserLog) TableName() string {
	return "public.userlog"
}
