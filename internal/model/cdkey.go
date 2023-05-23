/*
 * @Author: cloudyi.li
 * @Date: 2023-05-15 16:12:29
 * @LastEditTime: 2023-05-21 19:46:51
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/model/cdkey.go
 */

package model

type CdKeyGenerateRes struct {
	CodeKey []string `json:"code_key"`
}

type CdKeyGenerateReq struct {
	GiftCardId string `json:"card_id"  validate:"required"`
	KeyNumber  int    `json:"key_number"  validate:"required"`
}

type CdkeyPayReq struct {
	CodeKey string `json:"code_key"`
}

type CdKeyAmount struct {
	CodeKey    string  `gorm:"column:CdKeys__code_key" json:"code_key"`
	CardAmount float64 `gorm:"column:card_amount" json:"card_amount"`
}

type GiftCardCreate struct {
	CardName     string  `json:"card_name" validate:"required"`
	CardComment  string  `json:"card_comment" validate:"required"`
	CardAmount   float64 `json:"card_amount" validate:"required"`
	CardDiscount float64 `json:"card_discount" validate:"required"`
	CardLink     string  `json:"card_link" validate:"required"`
}

type GiftCardUpdate struct {
	CardId       string  `json:"card_id"`
	CardName     string  `json:"card_name"`
	CardComment  string  `json:"card_comment"`
	CardAmount   float64 `json:"card_amount"`
	CardDiscount float64 `json:"card_discount"`
	CardLink     string  `json:"card_link"`
}

type GiftCardListRes struct {
	GiftCardList []GiftCardOneRes `json:"card_list"`
}

type GiftCardOneRes struct {
	CardId       string  `json:"card_id"`
	CardName     string  `json:"card_name"`
	CardComment  string  `json:"card_comment"`
	CardAmount   float64 `json:"card_amount"`
	CardDiscount float64 `json:"card_discount"`
	CardLink     string  `json:"card_link"`
}

type GiftCardOne struct {
	CardId       int64   `gorm:"column:id" json:"card_id"`
	CardName     string  `gorm:"column:card_name" json:"card_name"`
	CardComment  string  `gorm:"column:card_comment" json:"card_comment"`
	CardAmount   float64 `gorm:"column:card_amount" json:"card_amount"`
	CardDiscount float64 `gorm:"column:card_discount" json:"card_discount"`
	CardLink     string  `gorm:"column:card_link" json:"card_link"`
}
