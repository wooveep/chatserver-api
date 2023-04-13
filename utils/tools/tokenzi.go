/*
 * @Author: cloudyi.li
 * @Date: 2023-04-12 05:38:19
 * @LastEditTime: 2023-04-13 15:28:55
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/utils/tools/tokenzi.go
 */
package tools

import "math"

func Tokenzi(str string) int {
	a := len([]byte(str))
	b := float64(a) * 1.3
	return int(math.Ceil(b))

}
