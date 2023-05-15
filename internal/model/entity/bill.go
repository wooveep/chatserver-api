package entity

import (
	"chatserver-api/pkg/jtime"
)

type Bill struct {
	Id        int64          `gorm:"column:id;primary_key;" json:"id"`
	UserId    int64          `gorm:"column:user_id" json:"user_id"`
	Cost      int            `gorm:"column:cost" json:"cost"`
	Amount    int            `gorm:"column:amount" json:"amount"`
	Comment   string         `gorm:"column:comment" json:"comment"`
	CreatedAt jtime.JsonTime `gorm:"column:created_at" json:"created_at"`
	UpdatedAt jtime.JsonTime `gorm:"column:updated_at" json:"updated_at"`
}

func (Bill) TableName() string {
	return "public.bill"
}
