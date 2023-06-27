/*
 * @Author: cloudyi.li
 * @Date: 2023-04-11 11:21:30
 * @LastEditTime: 2023-06-15 14:09:01
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/utils/tools/value.go
 */
package tools

import "time"

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

func TimeConvert(str string) (string, string, string, error) {
	t, err := time.Parse("2006-01-02T15:04:05-0700", str)
	if err != nil {
		return "", "", "", err
	}

	// 将时间格式化为指定格式的字符串
	return t.Format("2006-01-02"), t.Format("15:04:05"), t.Format("Monday"), nil
}
