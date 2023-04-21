/*
 * @Author: cloudyi.li
 * @Date: 2023-03-29 11:34:45
 * @LastEditTime: 2023-04-21 16:52:35
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/validator/validations.go
 */
// author: maxf
// date: 2022-03-29 16:29
// version: 自定义校验器

package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

// 手机号码规则，以1开头的11位数字
var mobile, _ = regexp.Compile(`^1\d{10}$`)

func mobileValidator(fl validator.FieldLevel) bool {
	phoneNumber := fl.Field().String()
	return mobile.MatchString(phoneNumber)
}

var usernamereg = regexp.MustCompile(`(?i)admin`)

func usernameValidator(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	return !usernamereg.MatchString(username)
}
