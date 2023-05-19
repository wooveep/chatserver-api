/*
 * @Author: cloudyi.li
 * @Date: 2023-05-11 16:51:24
 * @LastEditTime: 2023-05-18 13:01:24
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/dao/query/cdkey.go
 */
package query

import (
	"chatserver-api/internal/dao"
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

func (cd *cdkeyDao) CdKeyQuery(ctx context.Context, keyId int64) (codeKey string, cdkeyAmout float64, err error) {
	var cdkey entity.CdKey
	err = cd.ds.Master().Where("id = ?", keyId).Select("code_key,amount").Find(&cdkey).Error
	if err != nil {
		return
	}
	codeKey = cdkey.CodeKey
	cdkeyAmout = cdkey.Amount
	return
}

func (cd *cdkeyDao) CdKeyGenerate(ctx context.Context, Cdkeylist []entity.CdKey) (err error) {
	return cd.ds.Master().CreateInBatches(&Cdkeylist, 100).Error
}
