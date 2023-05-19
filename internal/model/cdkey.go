/*
 * @Author: cloudyi.li
 * @Date: 2023-05-15 16:12:29
 * @LastEditTime: 2023-05-19 09:12:22
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/model/cdkey.go
 */

package model

type CdKeyGenerateRes struct {
	CodeKey []string `json:"code_key"`
}

type CdKeyGenerateReq struct {
	KeyNumber int     `json:"key_number"  validate:"required"`
	KeyAmount float64 `json:"key_amount"  validate:"required"`
}

type CdkeyPayReq struct {
	CodKey string `json:"code_key"`
}
