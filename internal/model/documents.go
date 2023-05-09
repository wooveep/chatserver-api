/*
 * @Author: cloudyi.li
 * @Date: 2023-05-07 11:20:06
 * @LastEditTime: 2023-05-08 15:08:53
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/model/documents.go
 */
package model

type DocsCompare struct {
	Body string `gorm:"column:body" json:"body"`
}

type DocsBatchList struct {
	BatchTitle string   `json:"batch_title" validate:"required"`
	BatchList  []string `json:"batch_list" validate:"required"`
}
