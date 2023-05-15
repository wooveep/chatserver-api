/*
 * @Author: cloudyi.li
 * @Date: 2023-03-29 10:30:37
 * @LastEditTime: 2023-05-15 16:10:42
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/utils/uuid/uuid.go
 */
package uuid

import (
	"chatserver-api/internal/consts"
	"fmt"
	"math/rand"
	"strings"
	"time"

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

func GenKey(uuid int64) {

}

// CB7B 8L5N TZBW NF

func IdToCode(id int64, codelen int, bitlen int64) string {
	mod := int64(0)
	res := ""
	for id != 0 {
		mod = id % bitlen
		id = id / bitlen
		res += string(consts.CodeBase[mod])
	}
	resLen := len(res)
	if resLen < codelen {
		res += consts.CodePad
		for i := 0; i < 6-resLen-1; i++ {
			rand.Seed(time.Now().UnixNano())
			res += string(consts.CodeBase[rand.Intn(int(bitlen))])
		}
	}
	crc := strings.ReplaceAll(ChecksumKey([]byte(res))[2:], "0", "V")

	return res + crc
}

// code转id
func CodeToId(code string, bitlen int64) int64 {
	res := int64(0)
	lenCode := len(code)
	//var baseArr [] byte = []byte(c.base)
	baseArr := []byte(consts.CodeBase) // 字符串进制转换为byte数组
	baseRev := make(map[byte]int)      // 进制数据键值转换为map
	for k, v := range baseArr {
		baseRev[v] = k
	}

	// 查找补位字符的位置
	isPad := strings.Index(code, consts.CodePad)
	if isPad != -1 {
		lenCode = isPad
	}

	r := 0
	for i := 0; i < lenCode; i++ {
		// 补充字符直接跳过
		if string(code[i]) == consts.CodePad {
			continue
		}
		index := baseRev[code[i]]
		b := int64(1)
		for j := 0; j < r; j++ {
			b *= bitlen
		}
		// pow 类型为 float64 , 类型转换太麻烦, 所以自己循环实现pow的功能
		//res += float64(index) * math.Pow(float64(32), float64(2))
		res += int64(index) * b
		r++
	}
	return res
}

var KEYTable = makeTable(0x80e2)

// Table is a 256-word table representing the polynomial for efficient processing.
type Table struct {
	entries  [256]uint16
	reversed bool
	noXOR    bool
}

func MakeTable(poly uint16) *Table {
	return makeTable(poly)
}

// MakeBitsReversedTable returns the Table constructed from the specified polynomial.
func MakeBitsReversedTable(poly uint16) *Table {
	return makeBitsReversedTable(poly)
}

// MakeTableNoXOR returns the Table constructed from the specified polynomial.
// Updates happen without XOR in and XOR out.
func MakeTableNoXOR(poly uint16) *Table {
	tab := makeTable(poly)
	tab.noXOR = true
	return tab
}

// makeTable returns the Table constructed from the specified polynomial.
func makeBitsReversedTable(poly uint16) *Table {
	t := &Table{
		reversed: true,
	}
	width := uint16(16)
	for i := uint16(0); i < 256; i++ {
		crc := i << (width - 8)
		for j := 0; j < 8; j++ {
			if crc&(1<<(width-1)) != 0 {
				crc = (crc << 1) ^ poly
			} else {
				crc <<= 1
			}
		}
		t.entries[i] = crc
	}
	return t
}

func makeTable(poly uint16) *Table {
	t := &Table{
		reversed: false,
	}
	for i := 0; i < 256; i++ {
		crc := uint16(i)
		for j := 0; j < 8; j++ {
			if crc&1 == 1 {
				crc = (crc >> 1) ^ poly
			} else {
				crc >>= 1
			}
		}
		t.entries[i] = crc
	}
	return t
}

func updateBitsReversed(crc uint16, tab *Table, p []byte) uint16 {
	for _, v := range p {
		crc = tab.entries[byte(crc>>8)^v] ^ (crc << 8)
	}
	return crc
}

func update(crc uint16, tab *Table, p []byte) uint16 {
	crc = ^crc

	for _, v := range p {
		crc = tab.entries[byte(crc)^v] ^ (crc >> 8)
	}

	return ^crc
}

func updateNoXOR(crc uint16, tab *Table, p []byte) uint16 {
	for _, v := range p {
		crc = tab.entries[byte(crc)^v] ^ (crc >> 8)
	}

	return crc
}

func Update(crc uint16, tab *Table, p []byte) uint16 {
	if tab.reversed {
		return updateBitsReversed(crc, tab, p)
	} else if tab.noXOR {
		return updateNoXOR(crc, tab, p)
	} else {
		return update(crc, tab, p)
	}
}

func ChecksumKey(data []byte) string {

	return fmt.Sprintf("%04X", Update(0, KEYTable, data))

}

// func Checksum(data []byte, tab *Table) uint16 { return Update(0, tab, data) }

// ParseUUIDFromStr 从指定的字符串生成uuid
// func ParseUUIDFromStr(str string) (string, error) {
// 	u, err := ParseUUIDFromStr(str)
// 	if err != nil {
// 		return "", err
// 	}
// 	return u, nil
// }
