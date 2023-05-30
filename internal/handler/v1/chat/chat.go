/*
 * @Author: cloudyi.li
 * @Date: 2023-03-29 13:43:42
 * @LastEditTime: 2023-05-29 16:57:26
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
	"chatserver-api/pkg/tika"
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
		pre_token := tiktoken.NumTokensFromMessages(openAIReq.Messages, openAIReq.Model) + openAIReq.MaxTokens
		if pre_token >= consts.ModelMaxToken[openAIReq.Model] {
			logger.Debugf("预验证TOKEN，超出模型内存%d", pre_token)
			response.JSON(ctx, errors.WithCode(ecode.OversizeErr, "问题过长超出模型内存"), nil)
			return
		}
		pre_cost := float64(pre_token) * consts.TokenPrice
		if ctx.GetFloat64(consts.BalanceCtx) < pre_cost {
			logger.Debugf("预验证TOKEN，余额不足%f", pre_cost)
			response.JSON(ctx, errors.WithCode(ecode.NotFoundErr, "用户余额不足"), nil)
			return
		}
		//go func 请求API；
		chanStream := make(chan string)
		go ch.cSrv.ChatStremResGenerate(ctx, openAIReq, chanStream)

		//返回生成信息；
		msgId, messages := ch.cSrv.ChatStreamResProcess(ctx, chanStream, questionId, answerId)

		ch.cSrv.ChatCostCalculate(ctx, []openai.ChatCompletionMessage{{Role: openai.ChatMessageRoleAssistant, Content: messages}}, openAIReq.Model)

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
		pre_token := tiktoken.NumTokensFromMessages(openAIReq.Messages, openAIReq.Model) + openAIReq.MaxTokens
		if pre_token >= consts.ModelMaxToken[openAIReq.Model] {
			logger.Debugf("预验证TOKEN，超出模型内存%d", pre_token)
			response.JSON(ctx, errors.WithCode(ecode.OversizeErr, "问题过长超出模型内存"), nil)
			return
		}
		pre_cost := float64(pre_token) * consts.TokenPrice
		if ctx.GetFloat64(consts.BalanceCtx) < pre_cost {
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
		msgId, messages := ch.cSrv.ChatStreamResProcess(ctx, chanStream, questionId, 0)
		ch.cSrv.ChatCostCalculate(ctx, []openai.ChatCompletionMessage{{Role: openai.ChatMessageRoleAssistant, Content: messages}}, openAIReq.Model)
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

func (ch *ChatHandler) ChatRecordClear() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req model.RecordClearReq
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
		if err := ch.cSrv.ChatUserVerify(ctx); err != nil {
			response.JSON(ctx, errors.WithCode(ecode.NotFoundErr, "会话ID不存在"), nil)
			return
		}
		err = ch.cSrv.ChatRecordClear(ctx)
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
		var presetId int64
		if len(req.PresetId) > 0 {
			presetId, err = strconv.ParseInt(req.PresetId, 10, 64)
			if err != nil {
				response.JSON(ctx, errors.WithCode(ecode.ValidateErr, "会话ID转换错误"), nil)
				return
			}

		} else {
			presetId = 0
		}
		if err != nil {
			response.JSON(ctx, errors.WithCode(ecode.ValidateErr, "会话ID转换错误"), nil)
			return
		}
		ctx.Set(consts.ChatID, chatId)
		if err := ch.cSrv.ChatUserVerify(ctx); err != nil {
			response.JSON(ctx, errors.WithCode(ecode.NotFoundErr, "会话ID不存在"), nil)
			return
		}
		err = ch.cSrv.ChatUpdate(ctx, req.ChatName, presetId)
		if err != nil {
			response.JSON(ctx, errors.Wrap(err, ecode.Unknown, "会话更新失败"), nil)

		} else {
			response.JSON(ctx, nil, nil)
		}
	}
}

// 通过接口传入字符串构建embedding存储
func (ch *ChatHandler) ChatEmbeddingString() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req model.DocsBatchList
		if err := ctx.ShouldBindJSON(&req); err != nil {
			response.JSON(ctx, errors.WithCode(ecode.ValidateErr, err.Error()), nil)
			return
		}
		textlist := req.BatchList
		textlen := len(textlist)
		for i := 0; i < textlen; i += 10 {
			end := i + 10
			if end > textlen {
				end = textlen
			}
			batchlist := textlist[i:end]
			embeddinglist, err := ch.cSrv.ChatEmbeddingGenerate(batchlist)
			if err != nil {
				response.JSON(ctx, nil, nil)
				return
			}
			for j := 0; j < len(batchlist); j++ {
				err = ch.cSrv.ChatEmbeddingSave(ctx, req.BatchTitle, batchlist[j], req.Classify, embeddinglist[j])
				if err != nil {
					response.JSON(ctx, nil, nil)
					return
				}
			}
		}
		response.JSON(ctx, nil, nil)

	}
}

// 上传文件并构建Embedding存储
func (ch *ChatHandler) ChatEmbeddingFile() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		file, err := ctx.FormFile("file")
		if err != nil {
			response.JSON(ctx, errors.WithCode(ecode.Unknown, err.Error()), nil)
			return
		}
		// 检查文件类型是否为PDF
		// if fileHeader := file.Header.Get("Content-Type"); fileHeader != "application/pdf" ||  {
		// 	// ctx.JSON(http.StatusBadRequest, gin.H{"error": "上传文件必须是PDF格式"})
		// 	response.JSON(ctx, errors.WithCode(ecode.ValidateErr, "上传文件必须是PDF格式"), nil)

		// 	return
		// }
		// 获取文件名
		fileHeader := file.Header.Get("Content-Type")

		title := ctx.PostForm("title")
		if title == "" {
			title = file.Filename
		}

		// 获取类别参数
		// classify := ctx.PostForm("classify")
		// 保存文件到本地
		err = ctx.SaveUploadedFile(file, "uploadfile/"+file.Filename)
		if err != nil {
			// ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			response.JSON(ctx, errors.WithCode(ecode.Unknown, err.Error()), nil)
			return
		}
		var textlist []string
		if fileHeader == "application/pdf" {
			textlist, err = tika.ReadPd3f(title, file.Filename)
			if err != nil {
				// ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				response.JSON(ctx, errors.WithCode(ecode.Unknown, err.Error()), nil)
				return
			}
		}
		if fileHeader == "text/markdown" {
			textlist, err = tika.ProcessMarkDown(file.Filename)
			if err != nil {
				// ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				response.JSON(ctx, errors.WithCode(ecode.Unknown, err.Error()), nil)
				return
			}
		}
		// textlen := len(textlist)
		// for i := 0; i < textlen; i += 10 {
		// 	end := i + 10
		// 	if end > textlen {
		// 		end = textlen
		// 	}
		// 	batchlist := textlist[i:end]
		// 	embeddinglist, err := ch.cSrv.ChatEmbeddingGenerate(batchlist)
		// 	if err != nil {
		// 		// ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// 		response.JSON(ctx, err, nil)
		// 		return
		// 	}
		// 	for j := 0; j < len(batchlist); j++ {
		// 		err = ch.cSrv.ChatEmbeddingSave(ctx, title, batchlist[j], classify, embeddinglist[j])
		// 		if err != nil {
		// 			// ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// 			response.JSON(ctx, err, nil)
		// 			return
		// 		}
		// 	}
		// }
		response.JSON(ctx, nil, map[string]interface{}{"list": textlist})
	}
}
