/*
 * @Author: cloudyi.li
 * @Date: 2023-05-10 14:06:46
 * @LastEditTime: 2023-05-24 11:21:28
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/active/active.go
 */
package active

import (
	"chatserver-api/pkg/cache"
	"chatserver-api/utils/security"
	"chatserver-api/utils/uuid"
	"context"
	"strconv"
	"time"
)

func getActiveCodeKey(code string) string {
	return "User_Active_Code_list:" + security.Md5(code)
}

func activeCodeSave(ctx context.Context, code string, userId int64) (err error) {
	timer := 172800 * time.Second
	rc := cache.GetRedisClient()
	err = rc.SetNX(ctx, getActiveCodeKey(code), userId, timer).Err()
	return err
}

func ActiveCodeCompare(ctx context.Context, code string, userId int64) bool {
	rc := cache.GetRedisClient()
	idstr, err := rc.Get(ctx, getActiveCodeKey(code)).Result()
	if idstr == "" || err != nil {
		return false
	}
	id, err := strconv.ParseInt(idstr, 10, 64)
	if id != userId || err != nil {
		return false
	}
	rc.Del(ctx, getActiveCodeKey(code))
	return true
}

func ActiveCodeGen(ctx context.Context, userId int64) (string, error) {
	uuid := uuid.GenUUID16()
	code := strconv.FormatInt(userId, 10) + uuid
	err := activeCodeSave(ctx, code, userId)
	return code, err
}
