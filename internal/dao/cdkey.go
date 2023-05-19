/*
 * @Author: cloudyi.li
 * @Date: 2023-05-11 16:50:51
 * @LastEditTime: 2023-05-18 13:01:36
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/dao/cdkey.go
 */
package dao

import (
	"chatserver-api/internal/model/entity"
	"context"
)

type CDkeyDao interface {
	CdKeyGenerate(ctx context.Context, Cdkeylist []entity.CdKey) (err error)
	CdKeyQuery(ctx context.Context, keyId int64) (codeKey string, cdkeyAmout float64, err error)
	CdKeyDelete(ctx context.Context, cdkeyId int64) (err error)
}
