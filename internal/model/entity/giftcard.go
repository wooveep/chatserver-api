/*
 * @Author: cloudyi.li
 * @Date: 2023-05-20 20:49:15
 * @LastEditTime: 2023-05-21 19:06:14
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/model/entity/giftcard.go
 */
package entity

import (
	"chatserver-api/pkg/jtime"

	"gorm.io/plugin/soft_delete"
)

type GiftCard struct {
	Id           int64                 `gorm:"column:id;primary_key;" json:"id"`
	CardName     string                `gorm:"column:card_name" json:"card_name"`
	CardComment  string                `gorm:"column:card_comment" json:"card_comment"`
	CardAmount   float64               `gorm:"column:card_amount" json:"card_amount"`
	CardDiscount float64               `gorm:"column:card_discount" json:"card_discount"`
	CardBuyLink  string                `gorm:"column:card_link" json:"card_link"`
	CreatedAt    jtime.JsonTime        `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    jtime.JsonTime        `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt    jtime.JsonTime        `gorm:"column:deleted_at" json:"deleted_at" `
	IsDel        soft_delete.DeletedAt `gorm:"softDelete:flag,DeletedAtField:DeletedAt"`
	CdKeys       []CdKey               `gorm:"foreignKey:giftcard_id;references:id"`
}

func (GiftCard) TableName() string {
	return "public.giftcard"
}
