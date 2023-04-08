/*
 * @Author: cloudyi.li
 * @Date: 2023-03-29 10:30:37
 * @LastEditTime: 2023-04-08 12:34:28
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/utils/uuid/uuid.go
 */
package uuid

import (
	"strings"

	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
)

// GenUUID 生成一个随机的唯一ID
func GenUUID() string {
	return uuid.NewString()
}

// GenUUID16 截取uuid前16位
func GenUUID16() string {
	uuidStr := uuid.NewString()
	uuidStr = strings.ReplaceAll(uuidStr, "-", "")
	return uuidStr[0:16]
}

func GenID() (id int64, err error) {
	node, err := snowflake.NewNode(1)
	id = node.Generate().Int64()
	return
}

// ParseUUIDFromStr 从指定的字符串生成uuid
// func ParseUUIDFromStr(str string) (string, error) {
// 	u, err := ParseUUIDFromStr(str)
// 	if err != nil {
// 		return "", err
// 	}
// 	return u, nil
// }
