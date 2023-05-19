/*
 * @Author: cloudyi.li
 * @Date: 2023-03-29 10:30:37
 * @LastEditTime: 2023-05-19 11:42:49
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/utils/uuid/uuid.go
 */
package uuid

import (
	"chatserver-api/internal/consts"
	"container/list"
	"fmt"
	"regexp"
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

// GetInvCodeByUIDUniqueNew 获取指定长度的邀请码
func GetInvCodeByUID(uid int64) string {
	// 放大 + 加盐
	uid = uid*3 + 22521358345
	AlphanumericSet := []rune(consts.InviteBase)
	var code []rune
	slIdx := make([]byte, 8)
	// 扩散
	for i := 0; i < 8; i++ {
		slIdx[i] = byte(uid % int64(len(AlphanumericSet)))                    // 获取 62 进制的每一位值
		slIdx[i] = (slIdx[i] + byte(i)*slIdx[0]) % byte(len(AlphanumericSet)) // 其他位与个位加和再取余（让个位的变化影响到所有位）
		uid = uid / int64(len(AlphanumericSet))                               // 相当于右移一位（62进制）
	}
	// 混淆
	for i := 0; i < 8; i++ {
		idx := (byte(i) * 5) % byte(8)
		code = append(code, AlphanumericSet[slIdx[idx]])
	}
	return string(code)
}

func IdToCode(uid int64) string {
	byteCode := Base34(uint64(uid))
	crc := strings.ReplaceAll(ChecksumKey(byteCode), "0", "V")
	pattern := `([A-Z0-9]{4})`
	re := regexp.MustCompile(pattern)
	codekey := string(byteCode) + crc
	shufflekey := shuffleString(codekey)
	formatted := re.ReplaceAllString(shufflekey, "$1-")
	// formatted = formatted[:len(formatted)-1] // 去除末尾的破折号
	return formatted
}

func CodeToId(code string) int64 {
	formatted := strings.ReplaceAll(code, "-", "")
	unshuffled := unshuffleString(formatted)
	codelen := len(unshuffled)
	codekey := unshuffled[:codelen-4]
	byteCode := []byte(codekey)
	crccode := unshuffled[codelen-4:]
	crc := strings.ReplaceAll(ChecksumKey(byteCode), "0", "V")
	if crc != crccode {
		return 0
	}
	uid := Base34ToNum(byteCode)
	return int64(uid)
}

func Base34(n uint64) []byte {
	var base []byte = []byte(consts.CDKEYBASE)
	quotient := n
	mod := uint64(0)
	l := list.New()
	for quotient != 0 {
		mod = quotient % 32
		quotient = quotient / 32
		l.PushFront(base[int(mod)])
	}
	listLen := l.Len()
	if listLen >= 6 {
		res := make([]byte, 0, listLen)
		for i := l.Front(); i != nil; i = i.Next() {
			res = append(res, i.Value.(byte))
		}
		return res
	} else {
		res := make([]byte, 0, 6)
		for i := 0; i < 6; i++ {
			if i < 6-listLen {
				res = append(res, base[0])
			} else {
				res = append(res, l.Front().Value.(byte))
				l.Remove(l.Front())
			}

		}
		return res
	}
}

func Base34ToNum(str []byte) uint64 {
	baseMap := make(map[byte]int)
	var base []byte = []byte(consts.CDKEYBASE)
	for i, v := range base {
		baseMap[v] = i
	}
	var res uint64 = 0
	var r uint64 = 0
	for i := len(str) - 1; i >= 0; i-- {
		v, ok := baseMap[str[i]]
		if !ok {
			fmt.Printf("")
			return 0
		}
		var b uint64 = 1
		for j := uint64(0); j < r; j++ {
			b *= 32
		}
		res += b * uint64(v)
		r++
	}
	return res
}
