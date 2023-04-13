/*
 * @Author: cloudyi.li
 * @Date: 2023-03-29 10:30:37
 * @LastEditTime: 2023-04-12 21:34:40
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/utils/uuid/uuid.go
 */
package uuid

import (
	"strings"

	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
)

type SnowNode struct {
	Node *snowflake.Node
}

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
func NewNode(i int64) *SnowNode {
	node, err := snowflake.NewNode(i)
	if err != nil {
		panic(err)
	}
	return &SnowNode{
		Node: node,
	}
}

func (sn *SnowNode) GenSnowID() int64 {
	id := sn.Node.Generate().Int64()
	return id
}
func (sn *SnowNode) GenSnowStr() string {
	id := sn.Node.Generate().String()
	return id
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
