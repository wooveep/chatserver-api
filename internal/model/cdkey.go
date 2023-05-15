/*
 * @Author: cloudyi.li
 * @Date: 2023-05-15 16:12:29
 * @LastEditTime: 2023-05-15 16:28:16
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/model/cdkey.go
 */

package model

type CdKeyGenerateRes struct {
	CodeKey []string `json:"code_key"`
}

type CdKeyGenerateReq struct {
	KeyNumber int `json:"key_number"  validate:"required"`
	KeyAmount int `json:"key_amount"  validate:"required"`
}
