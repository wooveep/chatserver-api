/*
 * @Author: cloudyi.li
 * @Date: 2023-04-11 11:21:30
 * @LastEditTime: 2023-04-11 11:23:15
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/utils/tools/value.go
 */
package tools

func DefaultValue(a interface{}, b interface{}) interface{} {
	switch a.(type) {
	case bool:
		if a == false {
			return b.(bool)
		} else {
			return a.(bool)
		}
	case uint16:
		if a == 0 {
			return b.(uint16)
		} else {
			return a.(uint16)
		}
	case int:
		if a == 0 {
			return b.(int)
		} else {
			return a.(int)
		}
	case float64:
		if a == 0.0 {
			return b.(float64)
		} else {
			return a.(float64)
		}
	case string:
		if a == "" {
			return b.(string)
		} else {
			return a.(string)
		}
	default:
		return nil
	}
}
