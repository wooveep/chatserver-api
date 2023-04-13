/*
 * @Author: cloudyi.li
 * @Date: 2023-03-29 13:45:51
 * @LastEditTime: 2023-04-13 16:11:17
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/service/chat.go
 */
package service

import (
	"chatserver-api/internal/consts"
	"chatserver-api/internal/dao"
	"chatserver-api/internal/model"
	"chatserver-api/internal/model/entity"
	"chatserver-api/pkg/logger"
	"chatserver-api/pkg/openai"
	"chatserver-api/pkg/response"
	"chatserver-api/utils/security"
	"chatserver-api/utils/tools"
	"chatserver-api/utils/uuid"

	"context"
	"encoding/json"
	"errors"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

var _ ChatService = (*chatService)(nil)

type ChatService interface {
	ChatGenStremResponse(req openai.ChatCompletionRequest, closeWorker <-chan bool, chanStream chan<- string)
	ChatReqMessageProcess(ctx context.Context, chatid, userid int64) (req openai.ChatCompletionRequest, err error)
	ChatResMessageProcess(ctx *gin.Context, chanStream <-chan string) (messages string, err error)
	ChatCreateNewProcess(ctx context.Context, userid int64, newchatreq *model.ChatCreateNewReq) (newchatres model.ChatCreateNewRes, err error)
	ChatGetList(ctx context.Context, userid int64) (res model.ChatListRes, err error)
	ChatMessageSave(ctx context.Context, role, message string, chatid, userid int64) (err error)
	ChatDetailGet(ctx context.Context, chatid, userid int64) (res model.ChatDetailRes, err error)
	ChatRecordGet(ctx context.Context, chatid, userid int64) (res model.RecordHistoryRes, err error)
	ChatDelete(ctx context.Context, useid, chatid int64) (err error)
	ChatBalanceVerify(ctx context.Context, userid int64) (balance float64, err error)
	ChatCostCalculate(ctx context.Context, userid int64, balance float64, promptMsgs []openai.ChatCompletionMessage, genMsg string) error
}

// userService 实现UserService接口
type chatService struct {
	cd   dao.ChatDao
	iSrv uuid.SnowNode
}

func NewChatService(_cd dao.ChatDao) *chatService {
	return &chatService{
		cd:   _cd,
		iSrv: *uuid.NewNode(1),
	}
}

func (cs *chatService) ChatMessageSave(ctx context.Context, role, message string, chatid, userid int64) (err error) {
	recocd := entity.Record{}
	recocd.Id = cs.iSrv.GenSnowID()
	recocd.ChatId = chatid
	recocd.Sender = role
	recocd.Message = message
	recocd.MessageHash = security.Md5(message)
	err = cs.cd.ChatRecordSave(ctx, &recocd)
	return
}

func (cs *chatService) ChatDetailGet(ctx context.Context, chatid, userid int64) (res model.ChatDetailRes, err error) {
	chatdet, err := cs.cd.ChatDetailGet(ctx, userid, chatid)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	res.MaxTokens = chatdet.MaxTokens
	res.ModelName = chatdet.ModelName
	res.PresetContent = chatdet.PresetContent
	res.MemoryLevel = chatdet.MemoryLevel
	return
}

func (cs *chatService) ChatRecordGet(ctx context.Context, chatid, userid int64) (res model.RecordHistoryRes, err error) {
	count, err := cs.cd.ChatUserVerify(ctx, userid, chatid)
	if count == 0 || err != nil {
		err = errors.Join(errors.New("用户会话校验失败"))
		return
	}
	recordlist, err := cs.cd.ChatRecordGet(ctx, chatid, -1)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	res.Records = recordlist
	res.Id = chatid
	return
}

func (cs *chatService) ChatReqMessageProcess(ctx context.Context, chatid, userid int64) (req openai.ChatCompletionRequest, err error) {
	var chatMessages []openai.ChatCompletionMessage
	var chatmessage openai.ChatCompletionMessage
	var logitbia map[string]int

	preset, err := cs.cd.ChatDetailGet(ctx, userid, chatid)
	if err != nil {
		logger.Errorf("获取会话详情失败: %v\n", err)
		return
	}
	data, err := preset.LogitBias.MarshalJSON()
	if err != nil {
		logger.Errorf("序列化LogitBias失败: %v\n", err)
		return
	}
	err = json.Unmarshal(data, &logitbia)
	if err != nil {
		logger.Errorf("序列化LogitBias失败: %v\n", err)
		return
	}
	records, err := cs.cd.ChatRecordGet(ctx, chatid, preset.MemoryLevel)
	if err != nil {
		logger.Errorf("获取会话消息记录失败: %v\n", err)
		return
	}
	chatmessage.Role = openai.ChatMessageRoleSystem
	chatmessage.Content = preset.PresetContent
	chatMessages = append(chatMessages, chatmessage)
	for _, record := range records {
		chatmessage.Role = record.Sender
		chatmessage.Content = record.Message
		logger.Debugf("ROLE:%s,Message:%s", record.Sender, record.Message)
		chatMessages = append(chatMessages, chatmessage)
	}
	req.Model = preset.ModelName
	req.Stream = true
	req.MaxTokens = preset.MaxTokens
	req.LogitBias = logitbia
	req.Temperature = float32(preset.Temperature)
	req.TopP = float32(preset.TopP)
	req.FrequencyPenalty = float32(preset.Frequency)
	req.PresencePenalty = float32(preset.Presence)
	req.Messages = chatMessages
	return
}

func (cs *chatService) ChatResMessageProcess(ctx *gin.Context, chanStream <-chan string) (messages string, err error) {

	msgtime := time.Now().Format(consts.TimeLayout)
	ctx.Stream(func(w io.Writer) bool {
		if msg, ok := <-chanStream; ok {
			ctx.SSEvent("chatting", response.UnifyRes(ctx, nil, map[string]string{"time": msgtime, "msg": msg}))
			messages += msg
			logger.Debugf("stream-event: ID:%s ,time:%s,msg:%s", ctx.GetString(consts.RequestId), msgtime, msg)
			return true
		}
		ctx.SSEvent("chatting", response.UnifyRes(ctx, nil, map[string]string{"time": msgtime, "msg": "[DONE]"}))
		return false
	})
	logger.Debugf("Stream-message:%s", messages)
	return
}

func (cs *chatService) ChatGenStremResponse(req openai.ChatCompletionRequest, closeWorker <-chan bool, chanStream chan<- string) {
	defer close(chanStream)
	client, err := openai.NewClient()
	if err != nil {
		return
	}
	stream, err := client.CreateChatCompletionStream(req)
	if err != nil {
		logger.Errorf("ChatCompletionStream error: %v\n", err)
		return
	}
	defer stream.Close()

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			// chanStream <- "[DONE]"
			logger.Info("Stream finished")
			return
		}
		if err != nil {
			logger.Errorf("Stream error: %v\n", err)
			return
		}
		select {
		case <-closeWorker:
			return
		default:
			chanStream <- response.Choices[0].Delta.Content
			// logger.Debugf("MessageGen:%s", response.Choices[0].Delta.Content)
		}
	}
}

func (cs *chatService) ChatCreateNewProcess(ctx context.Context, userid int64, req *model.ChatCreateNewReq) (res model.ChatCreateNewRes, err error) {
	chat := entity.Chat{}
	chat.Id = cs.iSrv.GenSnowID()
	res.ChatId = chat.Id
	chat.UserId = userid
	chat.PresetId = req.PresetId
	chat.ChatName = req.ChatName
	chat.MemoryLevel = req.MemoryLevel
	err = cs.cd.ChatCreateNew(ctx, &chat)
	if err != nil {
		return
	}
	return
}

func (cs *chatService) ChatDelete(ctx context.Context, userid, chatid int64) error {

	if chatid == -1 {
		return cs.cd.ChatDeleteAll(ctx, userid)
	} else {
		count, err := cs.cd.ChatUserVerify(ctx, userid, chatid)
		if count == 0 || err != nil {
			err = errors.Join(errors.New("用户会话校验失败"))
			return err
		}
		return cs.cd.ChatDeleteOne(ctx, userid, chatid)
	}
}

func (cs *chatService) ChatGetList(ctx context.Context, userid int64) (res model.ChatListRes, err error) {
	chatlist, err := cs.cd.ChatGetList(ctx, userid)
	if err != nil {
		return
	}
	res.ChatList = chatlist
	return
}

func (cs *chatService) ChatBalanceVerify(ctx context.Context, userid int64) (balance float64, err error) {
	userbalance, err := cs.cd.ChatBalanceGet(ctx, userid)
	return userbalance.Balance, err
}
func (cs *chatService) ChatCostCalculate(ctx context.Context, userid int64, balance float64, promptMsgs []openai.ChatCompletionMessage, genMsg string) error {
	var promptMsg string
	for _, value := range promptMsgs {
		promptMsg += value.Content
	}
	promptToken := tools.Tokenzi(promptMsg)
	promptCost := float64(promptToken) * 0.00025
	genToken := tools.Tokenzi(genMsg)
	genCost := float64(genToken) * 0.00025
	totalCost := genCost + promptCost
	newBalance := balance - totalCost
	logger.Debugf("User:%d promptToken:%dpromptCost:%fgenToken:%dgenCost:%ftotalCost:%fnewBalance%f", userid, promptToken, promptCost, genToken, genCost, totalCost, newBalance)
	err := cs.cd.ChatCostUpdate(ctx, userid, newBalance)
	return err
}
