/*
 * @Author: cloudyi.li
 * @Date: 2023-03-29 10:30:37
 * @LastEditTime: 2023-05-19 11:14:47
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/utils/uuid/uuid_test.go
 */
package uuid

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

// YE5E F49T 458B 8AE2

func TestBase34(t *testing.T) {
	aSrv := *NewNode(5)
	for i := 0; i < 10; i++ {
		t.Run(strconv.FormatInt(int64(i), 10), func(t *testing.T) {
			id := aSrv.GenSnowID()
			got := Base34(uint64(id))
			got2 := Base34ToNum(got)
			fmt.Printf("ID:=%v, code:=  %v , want:= %v\n", id, string(got), got2)
			if !reflect.DeepEqual(uint64(id), got2) {
				t.Errorf("ID:=%v, code:=  %v , want:= %v", id, string(got), got2)
			}
		})
	}
}

func TestCodeToID(t *testing.T) {
	aSrv := *NewNode(5)
	for i := 0; i < 1000; i++ {
		t.Run(strconv.FormatInt(int64(i), 10), func(t *testing.T) {
			id := aSrv.GenSnowID()
			got := IdToCode(id)
			got2 := CodeToId(got)
			fmt.Printf("ID:=%v, code:=  %v , want:= %v\n", id, got, got2)
			if !reflect.DeepEqual(id, got2) {
				t.Errorf("ID:=%v, code:=  %v , want:= %v", id, got, got2)
			}
		})
	}
}

// VIVT 7RAS AHAH V
