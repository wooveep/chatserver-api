/*
 * @Author: cloudyi.li
 * @Date: 2023-05-11 16:51:24
 * @LastEditTime: 2023-05-11 16:53:59
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/dao/query/cdkey.go
 */
package query

import (
	"chatserver-api/internal/dao"
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

func (cd *cdkeyDao) CdKeyVerify(ctx context.Context) {

}
