/*
 * @Author: cloudyi.li
 * @Date: 2023-03-29 13:43:42
 * @LastEditTime: 2023-03-31 22:27:50
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/handler/v1/chat/chat.go
 */
package chat

import (
	"chatserver-api/di/logger"
	"chatserver-api/internal/model"
	"chatserver-api/internal/service"
	"chatserver-api/pkg/response"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	chatSrv service.ChatService
}

func NewChatHandler(_chatSrv service.ChatService) *ChatHandler {

	chah := &ChatHandler{
		chatSrv: _chatSrv,
	}
	return chah
}

func (chah *ChatHandler) ChattingStreamSend() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var chatChattingReq model.ChatChattingReq
		chanStream := make(chan string)
		if err := ctx.Bind(&chatChattingReq); err != nil {
			response.JSON(ctx, err, nil)
		}

		go chah.chatSrv.GetChatResponse(ctx.Writer.CloseNotify(), chanStream)
		ctx.Stream(func(w io.Writer) bool {
			if msg, ok := <-chanStream; ok {
				ctx.SSEvent("chatting", response.UnifyRes(ctx, nil, map[string]string{"a": "b", "msg": msg}))
				logger.Infof("stream-event:%d", time.Now().UnixMilli())
				return true
			}
			return false
		})
	}
}
