/*
 * @Author: cloudyi.li
 * @Date: 2023-04-04 19:40:39
 * @LastEditTime: 2023-04-08 08:08:58
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/model/entity/user.go
 */
package entity

import "chatserver-api/pkg/jtime"

type User struct {
	Id           int64           `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	Username     string          `gorm:"column:username" json:"username"`
	Nickname     string          `gorm:"column:nickname" json:"nickname"`
	Email        string          `gorm:"column:email" json:"email"`
	Phone        string          `gorm:"column:phone" json:"phone"`
	AvatarUrl    string          `gorm:"column:avatar_url" json:"avatar_url"`
	Password     string          `gorm:"column:password" json:"password"`
	RegisteredIp string          `gorm:"column:registered_ip" json:"registered_ip"`
	IsActive     bool            `gorm:"column:is_active" json:"is_active"`
	Balance      float64         `gorm:"column:balance" json:"balance"`
	CreatedAt    *jtime.JsonTime `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    *jtime.JsonTime `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (User) TableName() string {
	return "public.user"
}
