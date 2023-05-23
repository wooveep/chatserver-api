/*
 * @Author: cloudyi.li
 * @Date: 2023-05-11 16:50:51
 * @LastEditTime: 2023-05-21 19:39:22
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/dao/cdkey.go
 */
package dao

import (
	"chatserver-api/internal/model"
	"chatserver-api/internal/model/entity"
	"context"
)

type CDkeyDao interface {
	CdKeyGenerate(ctx context.Context, Cdkeylist []entity.CdKey) (err error)
	CdKeyQuery(ctx context.Context, cdkeyId int64) (model.CdKeyAmount, error)
	CdKeyDelete(ctx context.Context, cdkeyId int64) (err error)
	GiftCardCreate(ctx context.Context, giftcard *entity.GiftCard) error
	GiftCardListGet(ctx context.Context) ([]model.GiftCardOne, error)
	GiftCardUpdate(ctx context.Context, giftcard *entity.GiftCard) error
}
