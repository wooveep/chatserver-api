/*
 * @Author: cloudyi.li
 * @Date: 2023-03-29 13:43:42
 * @LastEditTime: 2023-04-23 15:39:33
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
	"chatserver-api/pkg/tiktoken"
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
		ctx.Set(consts.ChatID, chatId)
		if err := ch.cSrv.ChatUserVerify(ctx); err != nil {
			response.JSON(ctx, errors.WithCode(ecode.NotFoundErr, "会话ID不存在"), nil)
			return
		}
		err = ch.cSrv.ChatBalanceVerify(ctx)
		if err != nil {
			response.JSON(ctx, errors.WithCode(ecode.NotFoundErr, "用户余额不足"), nil)
			return
		}
		//生成请求信息；
		answerId, openAIReq, err := ch.cSrv.ChatRegenerategReqProcess(ctx, questionId, req.MemoryLevel)
		logger.Debugf("answerid %d", answerId)
		if err != nil {
			response.JSON(ctx, errors.Wrap(err, ecode.Unknown, "消息请求生成失败"), nil)
			return
		}
		//验证请求体余额
		pre_cost := float64(tiktoken.NumTokensFromMessages(openAIReq.Messages, openAIReq.Model)+openAIReq.MaxTokens) * consts.TokenPrice
		if ctx.GetFloat64(consts.Balance) < pre_cost {
			logger.Debugf("预验证TOKEN，余额不足%f", pre_cost)
			response.JSON(ctx, errors.WithCode(ecode.NotFoundErr, "用户余额不足"), nil)
			return
		}
		//go func 请求API；
		chanStream := make(chan string)

		go ch.cSrv.ChatStremResGenerate(ctx, openAIReq, chanStream)

		//返回生成信息；
		msgId, messages := ch.cSrv.ChatResProcess(ctx, chanStream, questionId, answerId)

		//保存生成信息;
		if err := ch.cSrv.ChatMessageSave(ctx, openai.ChatMessageRoleAssistant, messages, msgId); err != nil {
			logger.Errorf("生成问题消息保存失败:%s", err)
			return
		}
		if err := ch.cSrv.ChatBalanceUpdate(ctx); err != nil {
			logger.Errorf("保存计费消息失败:%s", err)
			return
		}
	}
}

func (ch *ChatHandler) ChatChatting() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := model.ChatChattingReq{}
		// 绑定JSON请求
		if err := ctx.ShouldBindJSON(&req); err != nil {
			response.JSON(ctx, errors.WithCode(ecode.ValidateErr, err.Error()), nil)
			return
		}
		//string转换int64
		chatId, err := strconv.ParseInt(req.ChatId, 10, 64)
		if err != nil {
			response.JSON(ctx, errors.WithCode(ecode.ValidateErr, "会话ID转换错误"), nil)
			return
		}
		//验证会话ID
		ctx.Set(consts.ChatID, chatId)
		if err := ch.cSrv.ChatUserVerify(ctx); err != nil {
			response.JSON(ctx, errors.WithCode(ecode.NotFoundErr, "会话ID不存在"), nil)
			return
		}
		//获取用户余额
		err = ch.cSrv.ChatBalanceVerify(ctx)
		if err != nil {
			response.JSON(ctx, errors.WithCode(ecode.NotFoundErr, "用户余额不足"), nil)
			return
		}
		//会话请求消息处理
		questionId, openAIReq, err := ch.cSrv.ChatChattingReqProcess(ctx, req.Message, req.MemoryLevel)
		if err != nil {
			response.JSON(ctx, errors.Wrap(err, ecode.Unknown, "消息请求生成失败"), nil)
			return
		}
		//验证请求体余额
		pre_cost := float64(tiktoken.NumTokensFromMessages(openAIReq.Messages, openAIReq.Model)+openAIReq.MaxTokens) * consts.TokenPrice
		if ctx.GetFloat64(consts.Balance) < pre_cost {
			logger.Debugf("预验证TOKEN，余额不足%f", pre_cost)
			response.JSON(ctx, errors.WithCode(ecode.NotFoundErr, "用户余额不足"), nil)
			return
		}
		//会话请求消息保存
		if err := ch.cSrv.ChatMessageSave(ctx, openai.ChatMessageRoleUser, req.Message, questionId); err != nil {
			response.JSON(ctx, errors.Wrap(err, ecode.Unknown, "用户请求消息保存失败"), nil)
			return
		}

		chanStream := make(chan string)
		//开始生成回答
		go ch.cSrv.ChatStremResGenerate(ctx, openAIReq, chanStream)
		//发送回答
		msgId, messages := ch.cSrv.ChatResProcess(ctx, chanStream, questionId, 0)
		//保存回答消息
		if err := ch.cSrv.ChatMessageSave(ctx, openai.ChatMessageRoleAssistant, messages, msgId); err != nil {
			logger.Errorf("生成问题消息保存失败:%s", err)
			return
		}
		if err := ch.cSrv.ChatBalanceUpdate(ctx); err != nil {
			logger.Errorf("保存计费消息失败:%s", err)
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

func (ch *ChatHandler) ChatListGet() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		res, err := ch.cSrv.ChatListGet(ctx)
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
		if err := ctx.ShouldBindJSON(&req); err != nil {
			response.JSON(ctx, errors.WithCode(ecode.ValidateErr, err.Error()), nil)
			return
		}
		chatId, err := strconv.ParseInt(req.ChatId, 10, 64)
		if err != nil {
			response.JSON(ctx, errors.WithCode(ecode.ValidateErr, "会话ID转换错误"), nil)
			return
		}
		ctx.Set(consts.ChatID, chatId)

		if err := ch.cSrv.ChatUserVerify(ctx); err != nil {
			response.JSON(ctx, errors.WithCode(ecode.NotFoundErr, "会话ID不存在"), nil)
			return
		}
		res, err := ch.cSrv.ChatDetailGet(ctx)
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
		if err := ctx.ShouldBindQuery(&req); err != nil {
			response.JSON(ctx, errors.WithCode(ecode.ValidateErr, err.Error()), nil)
			return
		}
		chatId, err := strconv.ParseInt(req.ChatId, 10, 64)
		if err != nil {
			response.JSON(ctx, errors.WithCode(ecode.ValidateErr, "会话ID转换错误"), nil)
			return
		}
		ctx.Set(consts.ChatID, chatId)
		if chatId != -1 {
			if err := ch.cSrv.ChatUserVerify(ctx); err != nil {
				response.JSON(ctx, errors.WithCode(ecode.NotFoundErr, "会话ID不存在"), nil)
				return
			}
		}
		res, err := ch.cSrv.ChatRecordGet(ctx)
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
		if err := ctx.ShouldBindQuery(&req); err != nil {
			response.JSON(ctx, errors.WithCode(ecode.ValidateErr, err.Error()), nil)
			return
		}
		chatId, err := strconv.ParseInt(req.ChatId, 10, 64)
		if err != nil {
			response.JSON(ctx, errors.WithCode(ecode.ValidateErr, "会话ID转换错误"), nil)
			return
		}
		ctx.Set(consts.ChatID, chatId)
		if chatId != -1 {
			if err := ch.cSrv.ChatUserVerify(ctx); err != nil {
				response.JSON(ctx, errors.WithCode(ecode.NotFoundErr, "会话ID不存在"), nil)
				return
			}
		}
		err = ch.cSrv.ChatDelete(ctx)
		if err != nil {
			response.JSON(ctx, errors.Wrap(err, ecode.Unknown, "接口调用失败"), nil)

		} else {
			response.JSON(ctx, nil, nil)
		}
	}
}

func (ch *ChatHandler) ChatUpdate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req model.ChatUpdateReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			response.JSON(ctx, errors.WithCode(ecode.ValidateErr, err.Error()), nil)
			return
		}
		chatId, err := strconv.ParseInt(req.ChatId, 10, 64)
		if err != nil {
			response.JSON(ctx, errors.WithCode(ecode.ValidateErr, "会话ID转换错误"), nil)
			return
		}
		ctx.Set(consts.ChatID, chatId)
		if err := ch.cSrv.ChatUserVerify(ctx); err != nil {
			response.JSON(ctx, errors.WithCode(ecode.NotFoundErr, "会话ID不存在"), nil)
			return
		}
		err = ch.cSrv.ChatUpdate(ctx, req.ChatName)
		if err != nil {
			response.JSON(ctx, errors.Wrap(err, ecode.Unknown, "会话更新失败"), nil)

		} else {
			response.JSON(ctx, nil, nil)
		}
	}
}
