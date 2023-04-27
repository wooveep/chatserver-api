/*
 * @Author: cloudyi.li
 * @Date: 2023-04-27 10:45:28
 * @LastEditTime: 2023-04-27 10:56:11
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/model/entity/documents.go
 */
package entity

import (
	"chatserver-api/pkg/jtime"
	"chatserver-api/pkg/pgvector"
)

type Documents struct {
	Id         int64           `gorm:"column:id;primary_key;" json:"id"`
	Title      string          `gorm:"column:title" json:"title"`
	Subsection string          `gorm:"column:subsection" json:"subsection"`
	Body       string          `gorm:"column:body" json:"body"`
	Tokens     int             `gorm:"column:tokens" json:"tokens"`
	Embedding  pgvector.Vector `gorm:"column:embedding" json:"embedding"`
	CreatedAt  jtime.JsonTime  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  jtime.JsonTime  `gorm:"column:updated_at" json:"updated_at"`
}

func (Documents) TableName() string {
	return "embed.documents"
}
