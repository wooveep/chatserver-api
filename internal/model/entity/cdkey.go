/*
 * @Author: cloudyi.li
 * @Date: 2023-05-10 09:58:37
 * @LastEditTime: 2023-05-18 13:01:13
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/model/entity/cdkey.go
 */
package entity

import (
	"chatserver-api/pkg/jtime"

	"gorm.io/plugin/soft_delete"
)

type CdKey struct {
	Id        int64                 `gorm:"column:id;primary_key;" json:"id"`
	CodeKey   string                `gorm:"column:code_key" json:"code_key"`
	Amount    float64               `gorm:"column:amount" json:"amount"`
	CreatedAt jtime.JsonTime        `gorm:"column:created_at" json:"created_at"`
	UpdatedAt jtime.JsonTime        `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt jtime.JsonTime        `gorm:"column:deleted_at" json:"deleted_at" `
	IsDel     soft_delete.DeletedAt `gorm:"softDelete:flag,DeletedAtField:DeletedAt"`
}

func (CdKey) TableName() string {
	return "public.cdkey"
}
