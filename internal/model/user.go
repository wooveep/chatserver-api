/*
 * @Author: cloudyi.li
 * @Date: 2023-04-04 15:44:35
 * @LastEditTime: 2023-04-04 19:41:13
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/model/user.go
 */
package model

type UserLoginReq struct {
	Mobile   string `json:"username" label:"用户名"`
	Password string `json:"password" label:"密码"`
}
