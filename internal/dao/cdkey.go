/*
 * @Author: cloudyi.li
 * @Date: 2023-05-11 16:50:51
 * @LastEditTime: 2023-05-15 13:48:45
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
	CdKeyVerify(ctx context.Context, codekey string) (cdkeyId int64, cdkeyAmout int, err error)
	CdKeyDelete(ctx context.Context, cdkeyId int64) (err error)
}
