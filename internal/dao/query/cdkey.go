/*
 * @Author: cloudyi.li
 * @Date: 2023-05-11 16:51:24
 * @LastEditTime: 2023-05-21 19:35:24
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/dao/query/cdkey.go
 */
package query

import (
	"chatserver-api/internal/dao"
	"chatserver-api/internal/model"
	"chatserver-api/internal/model/entity"
	"chatserver-api/pkg/db"
	"context"
)

var _ dao.CDkeyDao = (*cdkeyDao)(nil)

type cdkeyDao struct {
	ds db.IDataSource
}

func NewCDkeyDao(_ds db.IDataSource) *cdkeyDao {
	return &cdkeyDao{
		ds: _ds,
	}
}
func (cd *cdkeyDao) CdKeyDelete(ctx context.Context, cdkeyId int64) (err error) {
	return cd.ds.Master().Delete(&entity.CdKey{Id: cdkeyId}).Error
}

func (cd *cdkeyDao) CdKeyQuery(ctx context.Context, cdkeyId int64) (model.CdKeyAmount, error) {
	var cdkey model.CdKeyAmount
	err := cd.ds.Master().InnerJoins("CdKeys", cd.ds.Master().Where(&entity.CdKey{Id: cdkeyId})).Model(&entity.GiftCard{}).Find(&cdkey).Error
	return cdkey, err
}

func (cd *cdkeyDao) CdKeyGenerate(ctx context.Context, Cdkeylist []entity.CdKey) error {
	return cd.ds.Master().CreateInBatches(&Cdkeylist, 100).Error
}

func (cd *cdkeyDao) GiftCardCreate(ctx context.Context, giftcard *entity.GiftCard) error {
	return cd.ds.Master().Create(giftcard).Error
}

func (cd *cdkeyDao) GiftCardUpdate(ctx context.Context, giftcard *entity.GiftCard) error {
	return cd.ds.Master().Updates(giftcard).Error
}

func (cd *cdkeyDao) GiftCardListGet(ctx context.Context) ([]model.GiftCardOne, error) {
	var cardlist []model.GiftCardOne
	err := cd.ds.Master().Model(&entity.GiftCard{}).Find(&cardlist).Error
	return cardlist, err
}
