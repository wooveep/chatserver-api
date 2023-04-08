/*
 * @Author: cloudyi.li
 * @Date: 2023-03-29 13:43:42
 * @LastEditTime: 2023-04-08 12:59:06
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/handler/v1/chat/chat.go
 */
package chat

import (
	"chatserver-api/internal/model"
	"chatserver-api/internal/service"
	"chatserver-api/pkg/errors"
	"chatserver-api/pkg/errors/ecode"
	"chatserver-api/pkg/logger"
	"chatserver-api/pkg/response"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	cSrv service.ChatService
}

func NewChatHandler(_cSrv service.ChatService) *ChatHandler {

	ch := &ChatHandler{
		cSrv: _cSrv,
	}
	return ch
}

func (ch *ChatHandler) ChattingStreamSend() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		chatChattingReq := model.ChatChattingReq{}
		if err := ctx.ShouldBindJSON(&chatChattingReq); err != nil {
			response.JSON(ctx, errors.WithCode(ecode.ValidateErr, err.Error()), nil)
		} else {
			chanStream := make(chan string)
			genMessage := ch.cSrv.REQMessageProcess(chatChattingReq)
			go ch.cSrv.GetChatResponse(genMessage, ctx.Writer.CloseNotify(), chanStream)
			messages, err := ch.cSrv.RESMessageProcess(ctx, chanStream)
			if err != nil {
				response.JSON(ctx, err, nil)
			}
			logger.Info(messages)
		}
	}
}

func (ch *ChatHandler) CreateNewChat() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		chatCreateNewReq := model.ChatCreateNewReq{}
		if err := ctx.ShouldBindJSON(&chatCreateNewReq); err != nil {
			response.JSON(ctx, errors.WithCode(ecode.ValidateErr, err.Error()), nil)
		} else {
			chatCreateNewRes, err := ch.cSrv.ChatCreatNewProcess(ctx, &chatCreateNewReq)
			if err != nil {
				response.JSON(ctx, errors.WithCode(ecode.ValidateErr, err.Error()), nil)

			} else {
				response.JSON(ctx, nil, chatCreateNewRes)

			}
		}
	}

}
