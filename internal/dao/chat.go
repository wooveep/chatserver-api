/*
 * @Author: cloudyi.li
 * @Date: 2023-04-05 15:36:51
 * @LastEditTime: 2023-04-05 15:47:44
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/dao/chat.go
 */
package dao

import (
	"chatserver-api/internal/model/entity"
	"context"
)

type ChatDao interface {
	ChatCreateNew(ctx context.Context, chat *entity.ChatSession) error
}
