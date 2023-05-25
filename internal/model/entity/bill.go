/*
 * @Author: cloudyi.li
 * @Date: 2023-05-15 12:00:04
 * @LastEditTime: 2023-05-25 17:12:11
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/model/entity/bill.go
 */
package entity

import (
	"chatserver-api/pkg/jtime"
)

type Bill struct {
	Id          int64          `gorm:"column:id;primary_key;" json:"id"`
	UserId      int64          `gorm:"column:user_id" json:"user_id"`
	CostChange  float64        `gorm:"column:cost_change" json:"cost_change"`
	Balance     float64        `gorm:"column:balance" json:"balance"`
	CostComment string         `gorm:"column:cost_comment" json:"cost_comment"`
	CreatedAt   jtime.JsonTime `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   jtime.JsonTime `gorm:"column:updated_at" json:"updated_at"`
}

func (Bill) TableName() string {
	return "public.bill"
}
