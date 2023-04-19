/*
 * @Author: cloudyi.li
 * @Date: 2023-03-29 13:43:42
 * @LastEditTime: 2023-04-19 21:36:11
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
	"strconv"

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

func (ch *ChatHandler) ChatRegenerateg() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := model.ChatRegenerategReq{}
		userId := ctx.GetInt64(consts.UserID)
		if err := ctx.ShouldBindJSON(&req); err != nil {
			response.JSON(ctx, errors.WithCode(ecode.ValidateErr, err.Error()), nil)
			return
		}
		logger.Debugf("chatid%s,queid%s", req.ChatId, req.QuestionId)
		chatId, err := strconv.ParseInt(req.ChatId, 10, 64)
		questionId, err := strconv.ParseInt(req.QuestionId, 10, 64)
		if err != nil {
			response.JSON(ctx, errors.WithCode(ecode.ValidateErr, "ID转换错误"), nil)
			return
		}
		if err := ch.cSrv.ChatUserVerify(ctx, chatId, userId); err != nil {
			response.JSON(ctx, errors.WithCode(ecode.NotFoundErr, "会话ID不存在"), nil)
			return
		}
		balance, err := ch.cSrv.ChatBalanceVerify(ctx, userId)
		if balance < 0 || err != nil {
			response.JSON(ctx, errors.WithCode(ecode.NotFoundErr, "用户余额不足"), nil)
			return
		}
		//生成请求信息；
		answerId, openAIReq, err := ch.cSrv.ChatRegenerategReqProcess(ctx, chatId, questionId, userId, req.MemoryLevel)
		logger.Debugf("answerid %d", answerId)
		if err != nil {
			response.JSON(ctx, errors.Wrap(err, ecode.Unknown, "消息请求生成失败"), nil)
			return
		}
		//go func 请求API；
		chanStream := make(chan string)

		go ch.cSrv.ChatStremResGenerate(openAIReq, ctx.Writer.CloseNotify(), chanStream)

		//返回生成信息；
		msgId, messages := ch.cSrv.ChatResProcess(ctx, chanStream, questionId, answerId)

		//保存生成信息;
		if err := ch.cSrv.ChatMessageSave(ctx, openai.ChatMessageRoleAssistant, messages, msgId, chatId, userId); err != nil {
			logger.Errorf("生成问题消息保存失败:%s", err)
			return
		}

		if err := ch.cSrv.ChatCostCalculate(ctx, userId, balance, openAIReq.Messages, messages); err != nil {
			logger.Errorf("计费信息处理失败:%s", err)
			return
		}
	}
}

func (ch *ChatHandler) ChatChatting() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := model.ChatChattingReq{}
		userId := ctx.GetInt64(consts.UserID)
		if err := ctx.ShouldBindJSON(&req); err != nil {
			response.JSON(ctx, errors.WithCode(ecode.ValidateErr, err.Error()), nil)
			return
		}
		chatId, err := strconv.ParseInt(req.ChatId, 10, 64)
		if err != nil {
			response.JSON(ctx, errors.WithCode(ecode.ValidateErr, "会话ID转换错误"), nil)
			return
		}
		if err := ch.cSrv.ChatUserVerify(ctx, chatId, userId); err != nil {
			response.JSON(ctx, errors.WithCode(ecode.NotFoundErr, "会话ID不存在"), nil)
			return
		}

		balance, err := ch.cSrv.ChatBalanceVerify(ctx, userId)
		if balance < 0.00025 || err != nil {
			response.JSON(ctx, errors.WithCode(ecode.NotFoundErr, "用户余额不足"), nil)
			return
		}

		questionId, openAIReq, err := ch.cSrv.ChatChattingReqProcess(ctx, chatId, userId, req.Message, req.MemoryLevel)
		if err != nil {
			response.JSON(ctx, errors.Wrap(err, ecode.Unknown, "消息请求生成失败"), nil)
			return
		}

		if err := ch.cSrv.ChatMessageSave(ctx, openai.ChatMessageRoleUser, req.Message, questionId, chatId, userId); err != nil {
			response.JSON(ctx, errors.Wrap(err, ecode.Unknown, "用户请求消息保存失败"), nil)
			return
		}

		chanStream := make(chan string)
		go ch.cSrv.ChatStremResGenerate(openAIReq, ctx.Writer.CloseNotify(), chanStream)
		msgId, messages := ch.cSrv.ChatResProcess(ctx, chanStream, questionId, 0)

		if err := ch.cSrv.ChatMessageSave(ctx, openai.ChatMessageRoleAssistant, messages, msgId, chatId, userId); err != nil {
			logger.Errorf("生成问题消息保存失败:%s", err)
			return
		}
		if err := ch.cSrv.ChatCostCalculate(ctx, userId, balance, openAIReq.Messages, messages); err != nil {
			logger.Errorf("计费信息处理失败:%s", err)
			return
		}
	}
}

func (ch *ChatHandler) ChatCreateNew() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := model.ChatCreateNewReq{}
		userId := ctx.GetInt64(consts.UserID)
		if err := ctx.ShouldBindJSON(&req); err != nil {
			response.JSON(ctx, errors.WithCode(ecode.ValidateErr, err.Error()), nil)
			return

		}
		PresetId, err := strconv.ParseInt(req.PresetId, 10, 64)
		if err != nil {
			response.JSON(ctx, errors.WithCode(ecode.ValidateErr, "会话预设ID转换错误"), nil)
			return
		}
		chatCreateNewRes, err := ch.cSrv.ChatCreateNew(ctx, userId, PresetId, req.ChatName)
		if err != nil {
			response.JSON(ctx, errors.WithCode(ecode.ValidateErr, err.Error()), nil)

		} else {
			response.JSON(ctx, nil, chatCreateNewRes)

		}

	}
}

func (ch *ChatHandler) ChatSessionGetList() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.GetInt64(consts.UserID)
		res, err := ch.cSrv.ChatGetList(ctx, userId)
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
		userId := ctx.GetInt64(consts.UserID)

		if err := ctx.ShouldBindJSON(&req); err != nil {
			response.JSON(ctx, errors.WithCode(ecode.ValidateErr, err.Error()), nil)
			return
		}
		if err := ch.cSrv.ChatUserVerify(ctx, req.ChatId, userId); err != nil {
			response.JSON(ctx, errors.WithCode(ecode.NotFoundErr, "会话ID不存在"), nil)
			return
		}
		res, err := ch.cSrv.ChatDetailGet(ctx, req.ChatId, userId)
		if err != nil {
			response.JSON(ctx, errors.Wrap(err, ecode.Unknown, "接口调用失败"), nil)

		} else {
			response.JSON(ctx, nil, res)
		}
	}
}

func (ch *ChatHandler) ChatRecordHistory() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req model.RecordHistoryReq
		userId := ctx.GetInt64(consts.UserID)
		if err := ctx.ShouldBindQuery(&req); err != nil {
			response.JSON(ctx, errors.WithCode(ecode.ValidateErr, err.Error()), nil)
			return
		}
		if err := ch.cSrv.ChatUserVerify(ctx, req.ChatId, userId); err != nil {
			response.JSON(ctx, errors.WithCode(ecode.NotFoundErr, "会话ID不存在"), nil)
			return
		}
		res, err := ch.cSrv.ChatRecordGet(ctx, req.ChatId, userId)
		if err != nil {
			response.JSON(ctx, errors.Wrap(err, ecode.Unknown, "接口调用失败"), nil)

		} else {
			response.JSON(ctx, nil, res)
		}
	}
}

func (ch *ChatHandler) ChatDelete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req model.ChatDeleteReq
		userId := ctx.GetInt64(consts.UserID)
		if err := ctx.ShouldBindQuery(&req); err != nil {
			response.JSON(ctx, errors.WithCode(ecode.ValidateErr, err.Error()), nil)
			return
		}
		err := ch.cSrv.ChatDelete(ctx, userId, req.ChatId)
		if err != nil {
			response.JSON(ctx, errors.Wrap(err, ecode.Unknown, "接口调用失败"), nil)

		} else {
			response.JSON(ctx, nil, nil)
		}
	}
}
