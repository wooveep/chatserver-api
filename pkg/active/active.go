/*
 * @Author: cloudyi.li
 * @Date: 2023-05-10 14:06:46
 * @LastEditTime: 2023-05-31 18:51:25
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/active/active.go
 */
package active

import (
	"chatserver-api/pkg/cache"
	"chatserver-api/pkg/logger"
	"chatserver-api/utils/security"
	"chatserver-api/utils/uuid"
	"context"
	"strconv"

	"github.com/go-redis/redis/v8"
)

func getActiveCodeKey(code string) string {
	return "User_Active_Code_list:" + security.Md5(code)
}

func activeCodeSave(ctx context.Context, code string, userId int64) (err error) {
	// timer := 172800 * time.Second
	rc := cache.GetRedisClient()
	err = rc.SetNX(ctx, getActiveCodeKey(code), userId, 0).Err()
	return err
}

func ActiveCodeCompare(ctx context.Context, code string, userId int64) bool {
	rc := cache.GetRedisClient()
	idstr, err := rc.Get(ctx, getActiveCodeKey(code)).Result()
	if err != nil {
		if err != redis.Nil {
			logger.Errorf("Redis连接异常:%v", err.Error())
		}
		logger.Debugf("用户：%d,激活代码%s不存在", userId, code)
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
