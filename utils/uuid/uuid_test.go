/*
 * @Author: cloudyi.li
 * @Date: 2023-03-29 10:30:37
 * @LastEditTime: 2023-05-15 16:02:17
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/utils/uuid/uuid_test.go
 */
package uuid

import (
	"fmt"
	"testing"
)

func TestIdToCode(t *testing.T) {
	aSrv := *NewNode(5)

	tests := []struct {
		name string
	}{
		{
			name: "test1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// fmt.Println(tt.args.id)
			for i := 0; i < 10; i++ {
				id := aSrv.GenSnowID()
				got := IdToCode(id, 16, 32)
				// fmt.Println(got)
				code := got[:len(got)-2]
				crc := got[len(got)-2:]
				crc_cro := ChecksumKey([]byte(code))[2:]

				got2 := CodeToId(code, 32)
				fmt.Println(got)
				fmt.Println(got2)
				if got2 != id && crc != crc_cro {
					t.Errorf("%d", id)
				}
			}
		})
	}
}
