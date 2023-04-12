/*
 * @Author: cloudyi.li
 * @Date: 2023-03-29 13:43:42
 * @LastEditTime: 2023-04-12 16:50:46
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/handler/v1/chat/chat.go
 */
package chat

import (
	"chatserver-api/internal/consts"
	"chatserver-api/internal/model"
	"chatserver-api/internal/service"
	"chatserver-api/pkg/errors"
	"chatserver-api/pkg/errors/ecode"
	"chatserver-api/pkg/logger"
	"chatserver-api/pkg/openai"
	"chatserver-api/pkg/response"
	"context"

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
		req := model.ChatChattingReq{}
		userid := ctx.GetInt64(consts.UserID)
		if err := ctx.ShouldBindJSON(&req); err != nil {
			response.JSON(ctx, errors.WithCode(ecode.ValidateErr, err.Error()), nil)
			return
		}
		if err := ch.cSrv.ChatMessageSave(context.TODO(), openai.ChatMessageRoleUser, req.Message, req.ChatId, userid); err != nil {
			response.JSON(ctx, errors.Wrap(err, ecode.Unknown, "保存失败"), nil)
			return
		}
		chanStream := make(chan string)
		openAIReq, err := ch.cSrv.ChatReqMessageProcess(context.TODO(), req.ChatId, userid)
		if err != nil {
			response.JSON(ctx, errors.Wrap(err, ecode.Unknown, "消息请求生成失败"), nil)
			return
		}
		go ch.cSrv.ChatGenStremResponse(openAIReq, ctx.Writer.CloseNotify(), chanStream)
		messages, err := ch.cSrv.ChatResMessageProcess(ctx, chanStream)
		if err != nil {
			response.JSON(ctx, errors.Wrap(err, ecode.Unknown, "消息返回失败"), nil)
			return
		}
		if err := ch.cSrv.ChatMessageSave(context.TODO(), openai.ChatMessageRoleAssistant, messages, req.ChatId, userid); err != nil {
			logger.Errorf("消息保存失败:%s", messages)
			return
		}
	}
}

func (ch *ChatHandler) CreateNewChat() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		chatCreateNewReq := model.ChatCreateNewReq{}
		userid := ctx.GetInt64(consts.UserID)

		if err := ctx.ShouldBindJSON(&chatCreateNewReq); err != nil {
			response.JSON(ctx, errors.WithCode(ecode.ValidateErr, err.Error()), nil)
			return

		}
		chatCreateNewRes, err := ch.cSrv.ChatCreateNewProcess(ctx, userid, &chatCreateNewReq)
		if err != nil {
			response.JSON(ctx, errors.WithCode(ecode.ValidateErr, err.Error()), nil)

		} else {
			response.JSON(ctx, nil, chatCreateNewRes)

		}

	}

}

func (ch *ChatHandler) ChatSessionGetList() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userid := ctx.GetInt64(consts.UserID)
		res, err := ch.cSrv.ChatGetList(context.TODO(), userid)
		if err != nil {
			response.JSON(ctx, errors.Wrap(err, ecode.Unknown, "接口调用失败"), nil)

		} else {
			response.JSON(ctx, nil, res)

		}
	}
}

func (ch *ChatHandler) ChatDetailGet() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := model.ChatDetailReq{}
		userid := ctx.GetInt64(consts.UserID)

		if err := ctx.ShouldBindJSON(&req); err != nil {
			response.JSON(ctx, errors.WithCode(ecode.ValidateErr, err.Error()), nil)
			return

		}
		res, err := ch.cSrv.ChatDetailGet(context.TODO(), req.Id, userid)
		if err != nil {
			response.JSON(ctx, errors.Wrap(err, ecode.Unknown, "接口调用失败"), nil)

		} else {
			response.JSON(ctx, nil, res)
		}
	}
}
