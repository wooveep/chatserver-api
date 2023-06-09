/*
 * @Author: cloudyi.li
 * @Date: 2023-04-04 19:40:39
 * @LastEditTime: 2023-05-18 13:06:52
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/model/entity/user.go
 */
package entity

import (
	"chatserver-api/pkg/jtime"

	"gorm.io/plugin/soft_delete"
)

type User struct {
	Id           int64                 `gorm:"column:id;primary_key;" json:"id"`
	Username     string                `gorm:"column:username" json:"username"`
	Nickname     string                `gorm:"column:nickname" json:"nickname"`
	Email        string                `gorm:"column:email" json:"email"`
	Phone        string                `gorm:"column:phone" json:"phone"`
	AvatarUrl    string                `gorm:"column:avatar_url" json:"avatar_url"`
	Password     string                `gorm:"column:password" json:"password"`
	RegisteredIp string                `gorm:"column:registered_ip" json:"registered_ip"`
	IsActive     bool                  `gorm:"column:is_active" json:"is_active"`
	Balance      float64               `gorm:"column:balance" json:"balance"`
	Role         int                   `gorm:"column:role" json:"role"`
	CreatedAt    jtime.JsonTime        `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    jtime.JsonTime        `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt    jtime.JsonTime        `gorm:"column:deleted_at" json:"deleted_at" `
	IsDel        soft_delete.DeletedAt `gorm:"softDelete:flag,DeletedAtField:DeletedAt"`
	Chats        []Chat                `gorm:"foreignKey:user_id;references:id"`
	Invites      []Invite              `gorm:"foreignKey:user_id;references:id"`
	Bills        []Bill                `gorm:"foreignKey:user_id;references:id"`
	UserLogs     []UserLog             `gorm:"foreignKey:user_id;references:id"`
}

func (User) TableName() string {
	return "public.user"
}
