/*
 * @Author: cloudyi.li
 * @Date: 2023-03-29 13:45:51
 * @LastEditTime: 2023-04-19 21:35:00
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/service/chat.go
 */
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
	"chatserver-api/utils/security"
	"chatserver-api/utils/tools"
	"chatserver-api/utils/uuid"
	"strconv"

	"context"
	"encoding/json"
	"errors"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

var _ ChatService = (*chatService)(nil)

type ChatService interface {
	ChatStremResGenerate(req openai.ChatCompletionRequest, closeWorker <-chan bool, chanStream chan<- string)
	ChatChattingReqProcess(ctx context.Context, chatId, userId int64, lastquestion string, memoryLevel int16) (questionId int64, req openai.ChatCompletionRequest, err error)
	ChatRegenerategReqProcess(ctx context.Context, chatId, msgid, userId int64, memoryLevel int16) (answerid int64, req openai.ChatCompletionRequest, err error)
	ChatResProcess(ctx *gin.Context, chanStream <-chan string, questionId, answerId int64) (msgid int64, messages string)
	ChatCreateNew(ctx context.Context, userId, presetId int64, chatName string) (res model.ChatCreateNewRes, err error)
	ChatGetList(ctx context.Context, userId int64) (res model.ChatListRes, err error)
	ChatMessageSave(ctx context.Context, role, message string, msgid, chatId, userId int64) (err error)
	ChatDetailGet(ctx context.Context, chatId, userId int64) (res model.ChatDetailRes, err error)
	ChatRecordGet(ctx context.Context, chatId, userId int64) (res model.RecordHistoryRes, err error)
	ChatDelete(ctx context.Context, useid, chatId int64) (err error)
	ChatUserVerify(ctx context.Context, chatId, userId int64) (err error)
	ChatBalanceVerify(ctx context.Context, userId int64) (balance float64, err error)
	ChatCostCalculate(ctx context.Context, userId int64, balance float64, promptMsgs []openai.ChatCompletionMessage, genMsg string) error
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

func (cs *chatService) ChatMessageSave(ctx context.Context, role, message string, msgid, chatId, userId int64) (err error) {
	count, err := cs.cd.ChatRecordVerify(ctx, msgid)
	if err != nil {
		return
	}
	recocd := entity.Record{}
	recocd.Id = msgid
	recocd.ChatId = chatId
	recocd.Sender = role
	recocd.Message = message
	recocd.MessageHash = security.Md5(message)
	if count == 0 {
		logger.Debugf("聊天消息记录新建")
		err = cs.cd.ChatRecordSave(ctx, &recocd)
		return
	} else {
		logger.Debugf("聊天消息记录更新")
		err = cs.cd.ChatRecordUpdate(ctx, &recocd)
		return
	}
}

func (cs *chatService) ChatDetailGet(ctx context.Context, chatId, userId int64) (res model.ChatDetailRes, err error) {
	chatdet, err := cs.cd.ChatDetailGet(ctx, userId, chatId)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	res.MaxTokens = chatdet.MaxTokens
	res.ModelName = chatdet.ModelName
	res.PresetContent = chatdet.PresetContent
	return
}

func (cs *chatService) ChatUserVerify(ctx context.Context, chatId, userId int64) (err error) {
	count, err := cs.cd.ChatUserVerify(ctx, userId, chatId)
	if count == 0 || err != nil {
		err = errors.Join(errors.New("用户会话校验失败"))
		return
	}
	return
}

func (cs *chatService) ChatRecordGet(ctx context.Context, chatId, userId int64) (res model.RecordHistoryRes, err error) {
	recordlist, err := cs.cd.ChatRecordGet(ctx, chatId, -1)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	res.Records = recordlist
	res.ChatId = chatId
	return
}

func (cs *chatService) ChatRegenerategReqProcess(ctx context.Context, chatId, msgid, userId int64, memoryLevel int16) (answerid int64, req openai.ChatCompletionRequest, err error) {
	var chatMessages []openai.ChatCompletionMessage
	var systemPreset, historyMessage openai.ChatCompletionMessage
	var logitbia map[string]int

	preset, err := cs.cd.ChatDetailGet(ctx, userId, chatId)
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
	records, answerid, err := cs.cd.ChatRegenRecordGet(ctx, chatId, msgid, memoryLevel)
	logger.Debugf("answerid:%d", answerid)
	if err != nil {
		logger.Errorf("获取会话消息记录失败: %v\n", err)
		return
	}
	systemPreset.Role = openai.ChatMessageRoleSystem
	systemPreset.Content = preset.PresetContent
	chatMessages = append(chatMessages, systemPreset)
	for _, record := range records {
		historyMessage.Role = record.Sender
		historyMessage.Content = record.Message
		logger.Debugf("ROLE:%s,Message:%s", record.Sender, record.Message)
		chatMessages = append(chatMessages, historyMessage)
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

func (cs *chatService) ChatChattingReqProcess(ctx context.Context, chatId, userId int64, lastquestion string, memoryLevel int16) (questionId int64, req openai.ChatCompletionRequest, err error) {
	var chatMessages []openai.ChatCompletionMessage
	var systemPreset, historyMessage, lastMessage openai.ChatCompletionMessage

	var logitbia map[string]int

	preset, err := cs.cd.ChatDetailGet(ctx, userId, chatId)
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
	records, err := cs.cd.ChatRecordGet(ctx, chatId, memoryLevel)
	if err != nil {
		logger.Errorf("获取会话消息记录失败: %v\n", err)
		return
	}
	systemPreset.Role = openai.ChatMessageRoleSystem
	systemPreset.Content = preset.PresetContent
	chatMessages = append(chatMessages, systemPreset)
	for _, record := range records {
		historyMessage.Role = record.Sender
		historyMessage.Content = record.Message
		logger.Debugf("ROLE:%s,Message:%s", record.Sender, record.Message)
		chatMessages = append(chatMessages, historyMessage)
	}
	lastMessage.Role = openai.ChatMessageRoleUser
	lastMessage.Content = lastquestion
	chatMessages = append(chatMessages, lastMessage)
	req.Model = preset.ModelName
	req.Stream = true
	req.MaxTokens = preset.MaxTokens
	req.LogitBias = logitbia
	req.Temperature = float32(preset.Temperature)
	req.TopP = float32(preset.TopP)
	req.FrequencyPenalty = float32(preset.Frequency)
	req.PresencePenalty = float32(preset.Presence)
	req.Messages = chatMessages
	questionId = cs.iSrv.GenSnowID()
	return
}

func (cs *chatService) ChatResProcess(ctx *gin.Context, chanStream <-chan string, questionId, answerid int64) (msgid int64, messages string) {

	msgtime := time.Now().Format(consts.TimeLayout)
	if answerid != 0 {
		msgid = answerid
	} else {
		msgid = cs.iSrv.GenSnowID()
	}
	ctx.Stream(func(w io.Writer) bool {
		if msg, ok := <-chanStream; ok {
			ctx.SSEvent("chatting", map[string]string{"question_id": strconv.FormatInt(questionId, 10), "msgid": strconv.FormatInt(msgid, 10), "time": msgtime, "delta": msg})
			messages += msg
			logger.Debugf("stream-event: ID:%s ,time:%s,msg:%s", ctx.GetString(consts.RequestId), msgtime, msg)
			return true
		}
		ctx.SSEvent("chatting", map[string]string{"question_id": strconv.FormatInt(questionId, 10), "msgid": strconv.FormatInt(msgid, 10), "time": msgtime, "text": messages})
		return false
	})
	logger.Debugf("Stream-message:%s", messages)
	return
}

func (cs *chatService) ChatStremResGenerate(req openai.ChatCompletionRequest, closeWorker <-chan bool, chanStream chan<- string) {
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
	// for _, v := range req.Messages[0].Content {
	// 	time.Sleep(100 * time.Millisecond)
	// 	select {
	// 	case <-closeWorker:
	// 		return
	// 	default:
	// 		chanStream <- string(v)
	// 		// logger.Debugf("MessageGen:%s", response.Choices[0].Delta.Content)
	// 	}
	// }
}

func (cs *chatService) ChatCreateNew(ctx context.Context, userId, presetId int64, chatName string) (res model.ChatCreateNewRes, err error) {
	chat := entity.Chat{}
	chat.Id = cs.iSrv.GenSnowID()
	chat.UserId = userId
	chat.PresetId = presetId
	chat.ChatName = chatName
	err = cs.cd.ChatCreateNew(ctx, &chat)
	if err != nil {
		return
	}
	res.ChatId = strconv.FormatInt(chat.Id, 10)
	return
}

func (cs *chatService) ChatDelete(ctx context.Context, userId, chatId int64) error {

	if chatId == -1 {
		return cs.cd.ChatDeleteAll(ctx, userId)
	} else {
		count, err := cs.cd.ChatUserVerify(ctx, userId, chatId)
		if count == 0 || err != nil {
			err = errors.Join(errors.New("用户会话校验失败"))
			return err
		}
		return cs.cd.ChatDeleteOne(ctx, userId, chatId)
	}
}

func (cs *chatService) ChatGetList(ctx context.Context, userId int64) (res model.ChatListRes, err error) {
	chatlist, err := cs.cd.ChatGetList(ctx, userId)
	if err != nil {
		return
	}
	res.ChatList = chatlist
	return
}

func (cs *chatService) ChatBalanceVerify(ctx context.Context, userId int64) (balance float64, err error) {
	userbalance, err := cs.cd.ChatBalanceGet(ctx, userId)
	return userbalance.Balance, err
}

func (cs *chatService) ChatCostCalculate(ctx context.Context, userId int64, balance float64, promptMsgs []openai.ChatCompletionMessage, genMsg string) error {
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
	logger.Debugf("User:%d promptToken:%dpromptCost:%fgenToken:%dgenCost:%ftotalCost:%fnewBalance%f", userId, promptToken, promptCost, genToken, genCost, totalCost, newBalance)
	err := cs.cd.ChatCostUpdate(ctx, userId, newBalance)
	return err
}
