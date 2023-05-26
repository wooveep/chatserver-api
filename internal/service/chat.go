/*
 * @Author: cloudyi.li
 * @Date: 2023-03-29 13:45:51
 * @LastEditTime: 2023-05-26 16:28:43
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
	"chatserver-api/pkg/cache"
	"chatserver-api/pkg/logger"
	"chatserver-api/pkg/openai"
	"chatserver-api/pkg/pgvector"
	"chatserver-api/pkg/tiktoken"
	"chatserver-api/utils/security"
	"chatserver-api/utils/uuid"
	"fmt"
	"strconv"
	"strings"

	"context"
	"encoding/json"
	"errors"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var _ ChatService = (*chatService)(nil)

type ChatService interface {
	ChatCreateNew(ctx context.Context, userId, presetId int64, chatName string) (res model.ChatCreateNewRes, err error)
	ChatDelete(ctx *gin.Context) error
	ChatUpdate(ctx *gin.Context, chatName string, presetId int64) error
	ChatListGet(ctx *gin.Context) (res model.ChatListRes, err error)
	ChatRecordGet(ctx *gin.Context) (res model.RecordHistoryRes, err error)
	ChatRecordClear(ctx *gin.Context) (err error)
	ChatBalanceVerify(ctx *gin.Context) (err error)
	ChatBalanceUpdate(ctx *gin.Context) (err error)
	ChatMessageSave(ctx *gin.Context, role, message string, msgid int64) (err error)
	ChatDetailGet(ctx *gin.Context) (res model.ChatDetailRes, err error)
	ChatUserVerify(ctx *gin.Context) (err error)
	ChatRegenerategReqProcess(ctx *gin.Context, msgid int64, memoryLevel int16) (answerid int64, req openai.ChatCompletionRequest, err error)
	ChatChattingReqProcess(ctx *gin.Context, lastquestion string, memoryLevel int16) (questionId int64, req openai.ChatCompletionRequest, err error)
	ChatStremResGenerate(ctx *gin.Context, req openai.ChatCompletionRequest, chanStream chan<- string)
	ChatStreamResProcess(ctx *gin.Context, chanStream <-chan string, questionId, answerid int64) (msgid int64, messages string)
	ChatEmbeddingSave(ctx context.Context, title, body, classify string, embeddata openai.Embedding) error
	ChatEmbeddingGenerate(str []string) (embedVectors []openai.Embedding, err error)
	ChatEmbeddingCompare(ctx context.Context, question, classify string) (contextStr string, err error)
	ChatCostCalculate(ctx *gin.Context, promptMsgs []openai.ChatCompletionMessage, model string)
	// ChatTest(ctx context.Context, text string) (keyword string)
}

// userService 实现UserService接口
type chatService struct {
	cd   dao.ChatDao
	uSrv UserService
	rc   *redis.Client

	iSrv uuid.SnowNode
}

func NewChatService(_cd dao.ChatDao, _uSrv UserService) *chatService {
	return &chatService{
		cd:   _cd,
		uSrv: _uSrv,
		iSrv: *uuid.NewNode(1),
		rc:   cache.GetRedisClient(),
	}
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
	cs.rc.SAdd(ctx, consts.UserChatIDPrefix+strconv.FormatInt(userId, 10), chat.Id)
	return
}

func (cs *chatService) ChatUpdate(ctx *gin.Context, chatName string, presetId int64) error {
	chatId := ctx.GetInt64(consts.ChatID)
	chat := entity.Chat{}
	chat.Id = chatId
	if presetId != 0 {
		chat.PresetId = presetId
	}
	if len(chatName) != 0 {
		chat.ChatName = chatName
	}
	return cs.cd.ChatUpdate(ctx, &chat)
}

func (cs *chatService) ChatDelete(ctx *gin.Context) error {
	userId := ctx.GetInt64(consts.UserID)
	chatId := ctx.GetInt64(consts.ChatID)
	if chatId == -1 {
		cs.rc.Del(ctx, consts.UserChatIDPrefix+strconv.FormatInt(userId, 10))
		return cs.cd.ChatDeleteAll(ctx, userId)
	} else {
		cs.rc.SRem(ctx, consts.UserChatIDPrefix+strconv.FormatInt(userId, 10), chatId)
		return cs.cd.ChatDeleteOne(ctx, userId, chatId)
	}
}

// 验证会话用户归属
func (cs *chatService) ChatUserVerify(ctx *gin.Context) (err error) {
	userId := ctx.GetInt64(consts.UserID)
	chatId := ctx.GetInt64(consts.ChatID)
	n, err := cs.rc.Exists(ctx, consts.UserChatIDPrefix+strconv.FormatInt(userId, 10)).Result()
	if n > 0 {
		exist, err := cs.rc.SIsMember(ctx, consts.UserChatIDPrefix+strconv.FormatInt(userId, 10), chatId).Result()
		if err != nil {
			logger.Errorf("Redis连接异常:%v", err.Error())
		}
		if exist {
			return nil
		}
	}
	count, err := cs.cd.ChatUserVerify(ctx, userId, chatId)
	if err != nil {
		return err
	}
	if count != 0 {
		return nil
	} else {
		return errors.New("用户会话校验失败")
	}
}

func (cs *chatService) ChatListGet(ctx *gin.Context) (res model.ChatListRes, err error) {
	userId := ctx.GetInt64(consts.UserID)
	var chatOne model.ChatOneRes
	var chatListRes []model.ChatOneRes
	chatlist, err := cs.cd.ChatListGet(ctx, userId)
	if err != nil {
		return
	}
	for _, v := range chatlist {
		cs.rc.SAdd(ctx, consts.UserChatIDPrefix+strconv.FormatInt(userId, 10), v.ChatId)
		chatOne.ChatId = strconv.FormatInt(v.ChatId, 10)
		chatOne.PresetId = strconv.FormatInt(v.PresetId, 10)
		chatOne.ChatName = v.ChatName
		chatOne.CreatedAt = v.CreatedAt
		chatListRes = append(chatListRes, chatOne)

	}
	res.ChatList = chatListRes
	return
}

func (cs *chatService) ChatDetailGet(ctx *gin.Context) (res model.ChatDetailRes, err error) {
	userId := ctx.GetInt64(consts.UserID)
	chatId := ctx.GetInt64(consts.ChatID)
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

func (cs *chatService) ChatRecordGet(ctx *gin.Context) (res model.RecordHistoryRes, err error) {
	chatId := ctx.GetInt64(consts.ChatID)
	var recordOne model.RecordOneRes
	var recordListRes []model.RecordOneRes
	recordlist, err := cs.cd.ChatRecordGet(ctx, chatId, -1)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	for i := 0; i < len(recordlist); i++ {
		recordOne.Id = strconv.FormatInt(recordlist[i].Id, 10)
		recordOne.Message = recordlist[i].Message
		recordOne.CreatedAt = recordlist[i].CreatedAt
		recordOne.Sender = recordlist[i].Sender
		cs.rc.SAdd(ctx, consts.ChatRecordIDPrefix+strconv.FormatInt(chatId, 10), recordlist[i].Id)
		recordListRes = append(recordListRes, recordOne)
	}
	res.Records = recordListRes
	res.ChatId = strconv.FormatInt(chatId, 10)
	return
}

func (cs *chatService) ChatRecordClear(ctx *gin.Context) (err error) {
	chatId := ctx.GetInt64(consts.ChatID)
	err = cs.cd.ChatRecordClear(ctx, chatId)
	cs.rc.Del(ctx, consts.ChatRecordIDPrefix+strconv.FormatInt(chatId, 10))
	return
}

func (cs *chatService) ChatBalanceVerify(ctx *gin.Context) (err error) {
	userId := ctx.GetInt64(consts.UserID)
	userbalance, err := cs.uSrv.UserGetBalance(ctx, userId)
	ctx.Set(consts.BalanceCtx, userbalance)
	if userbalance < float64(0) {
		err = errors.New("Insufficient account balance")
	}
	return err
}

func (cs *chatService) ChatBalanceUpdate(ctx *gin.Context) (err error) {
	userId := ctx.GetInt64(consts.UserID)
	balance := ctx.GetFloat64(consts.BalanceCtx)
	token := ctx.GetInt(consts.CostTokenCtx)
	cost := float64(token) * consts.TokenPrice
	comment := fmt.Sprintf("消费-会话消耗令牌数:%d", token)
	return cs.uSrv.UserBalanceChange(ctx, userId, balance, -cost, comment)
}

func (cs *chatService) ChatCostCalculate(ctx *gin.Context, promptMsgs []openai.ChatCompletionMessage, model string) {
	old_token := ctx.GetInt(consts.CostTokenCtx)
	token := tiktoken.NumTokensFromMessages(promptMsgs, model)
	ctx.Set(consts.CostTokenCtx, token+old_token)
	logger.Debugf("本次消耗TOKEN：%d", token)
	return
}

func (cs *chatService) ChatMessageSave(ctx *gin.Context, role, message string, msgid int64) (err error) {
	chatId := ctx.GetInt64(consts.ChatID)
	var exist bool
	n, err := cs.rc.Exists(ctx, consts.ChatRecordIDPrefix+strconv.FormatInt(chatId, 10)).Result()
	if n > 0 {
		exist, err = cs.rc.SIsMember(ctx, consts.ChatRecordIDPrefix+strconv.FormatInt(chatId, 10), msgid).Result()
	}
	if err != nil || n <= 0 {
		count, err := cs.cd.ChatRecordVerify(ctx, msgid)
		if err != nil {
			return err
		}
		if count != 0 {
			exist = true
		} else {
			exist = false
		}
	}
	record := entity.Record{}
	record.Id = msgid
	record.ChatId = chatId
	record.Sender = role
	record.Message = message
	record.MessageHash = security.Md5(message)
	record.MessageToken = tiktoken.NumTokensSingleString(message)
	if !exist {
		logger.Debugf("聊天消息记录新建")
		cs.rc.SAdd(ctx, consts.ChatRecordIDPrefix+strconv.FormatInt(chatId, 10), msgid)
		err = cs.cd.ChatRecordSave(ctx, &record)
		return
	} else {
		logger.Debugf("聊天消息记录更新")
		err = cs.cd.ChatRecordUpdate(ctx, &record)
		return
	}
}

func (cs *chatService) ChatRegenerategReqProcess(ctx *gin.Context, msgid int64, memoryLevel int16) (answerid int64, req openai.ChatCompletionRequest, err error) {
	var chatMessages []openai.ChatCompletionMessage
	var systemPreset, historyMessage openai.ChatCompletionMessage
	var logitbia map[string]int
	var emquestion, embedcontexts string
	userId := ctx.GetInt64(consts.UserID)
	chatId := ctx.GetInt64(consts.ChatID)
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
	if preset.WithEmbedding {
		if len(records) != 0 {
			for _, v := range records {
				if v.Sender == "user" {
					emquestion += v.Message
				}
			}
		}
		//将用户问题进行关键词提取
		// emkeyword := cs.jieba.GetKeyword(emquestion)
		//通过用户问题 lastquestion + records（User历史）获取Context信息
		embedcontexts, err = cs.ChatEmbeddingCompare(ctx, emquestion, preset.Classify)
		if err != nil {
			logger.Errorf("获取embedding上下文失败: %v\n", err)
			return
		}
		//替换拼接PresetContent
		systemPreset.Content = strings.Replace(preset.PresetContent, "{{ context }}", embedcontexts, -1)

	} else {
		systemPreset.Content = preset.PresetContent
	}

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
	cs.ChatCostCalculate(ctx, chatMessages[1:], preset.ModelName)
	return
}

func (cs *chatService) ChatChattingReqProcess(ctx *gin.Context, lastquestion string, memoryLevel int16) (questionId int64, req openai.ChatCompletionRequest, err error) {
	var chatMessages []openai.ChatCompletionMessage
	var systemPreset, historyMessage, lastMessage openai.ChatCompletionMessage
	var logitbia map[string]int
	var emquestion, embedcontexts string
	userId := ctx.GetInt64(consts.UserID)
	chatId := ctx.GetInt64(consts.ChatID)
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
	if preset.WithEmbedding {
		if len(records) != 0 {
			for _, v := range records {
				if v.Sender == "user" {
					emquestion += v.Message
				}
			}
		}

		emquestion += lastquestion
		//将用户问题进行关键词提取
		// emkeyword := cs.jieba.GetKeyword(emquestion)
		//通过用户问题 lastquestion + records（User历史）获取Context信息
		embedcontexts, err = cs.ChatEmbeddingCompare(ctx, emquestion, preset.Classify)
		if err != nil {
			logger.Errorf("获取embedding上下文失败: %v\n", err)
			return
		}
		//替换拼接PresetContent
		systemPreset.Content = strings.Replace(preset.PresetContent, "{{ context }}", embedcontexts, -1)

	} else {
		systemPreset.Content = preset.PresetContent
	}

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
	cs.ChatCostCalculate(ctx, chatMessages[1:], preset.ModelName)
	return
}

func (cs *chatService) ChatStreamResProcess(ctx *gin.Context, chanStream <-chan string, questionId, answerid int64) (msgid int64, messages string) {

	msgtime := time.Now().Format(consts.TimeLayout)
	if answerid != 0 {
		msgid = answerid
	} else {
		msgid = cs.iSrv.GenSnowID()
	}
	ctx.Stream(func(w io.Writer) bool {
		if msg, ok := <-chanStream; ok {
			if msg == "[content_filter]" {
				messages = "尊敬的客户，非常感谢您使用我们的服务。我们注意到您最近提交的问题被我们的内容过滤器拦截了。我们深表歉意，因为我们的过滤器是为了保护我们的用户免受不良内容的侵害而设置的。但是，我们也理解您的问题对您来说非常重要。如果您有任何疑问或需要进一步的帮助，请随时联系网站管理员。再次感谢您的支持和理解。"
				ctx.SSEvent("chatting", map[string]string{"question_id": strconv.FormatInt(questionId, 10), "msgid": strconv.FormatInt(msgid, 10), "time": msgtime, "text": messages})
				return false
			}
			if msg == "[REQ_ERROR]" {
				messages += "尊敬的客户，非常感谢您使用我们的服务。由于API暂时异常,我们深表歉意,请随时联系网站管理员。再次感谢您的支持和理解。"
				ctx.SSEvent("chatting", map[string]string{"question_id": strconv.FormatInt(questionId, 10), "msgid": strconv.FormatInt(msgid, 10), "time": msgtime, "text": messages})
				return false
			}
			messages += msg
			ctx.SSEvent("chatting", map[string]string{"question_id": strconv.FormatInt(questionId, 10), "msgid": strconv.FormatInt(msgid, 10), "time": msgtime, "delta": messages})
			// logger.Debugf("stream-event: ID:%s ,time:%s,msg:%s", ctx.GetString(consts.RequestId), msgtime, msg)
			return true
		}

		ctx.SSEvent("chatting", map[string]string{"question_id": strconv.FormatInt(questionId, 10), "msgid": strconv.FormatInt(msgid, 10), "time": msgtime, "text": messages})
		return false
	})
	logger.Debugf("Stream-message:%s", messages)
	return
}

func (cs *chatService) ChatStremResGenerate(ctx *gin.Context, req openai.ChatCompletionRequest, chanStream chan<- string) {
	var chatMessages []openai.ChatCompletionMessage
	var lastMessage, blankMessage openai.ChatCompletionMessage
	var resmessage string
	var reqnew openai.ChatCompletionRequest
	blankMessage.Content = "[cmd:continue]"
	blankMessage.Role = openai.ChatMessageRoleUser
	chatMessages = req.Messages
	// cs.ChatCostCalculate(ctx, chatMessages[1:], req.Model)
	for _, v := range chatMessages {
		logger.Debugf("role: %s ;content: %s ", v.Role, v.Content)
	}
	client, err := openai.NewClient()
	if err != nil {
		close(chanStream)
		return
	}
	stream, err := client.CreateChatCompletionStream(req)
	if err != nil {
		logger.Errorf("ChatCompletionStream error: %v\n", err)
		chanStream <- "[REQ_ERROR]"
		close(chanStream)
		return
	}
	for {
		if stream == nil {
			close(chanStream)
			return
		}
		response, err := stream.Recv()
		if len(response.Choices) == 1 {
			if response.Choices[0].FinishReason == "length" {
				logger.Debugf("chat请求ID：%s", response.ID)
				stream.Close()
				lastMessage.Content = resmessage
				lastMessage.Role = openai.ChatMessageRoleAssistant
				// cs.ChatCostCalculate(ctx, []openai.ChatCompletionMessage{lastMessage}, req.Model)
				chatMessages = append(chatMessages, lastMessage, blankMessage)
				reqnew = req
				reqnew.Messages = chatMessages
				if consts.ModelMaxToken[req.Model] < tiktoken.NumTokensFromMessages(chatMessages, req.Model)+req.MaxTokens {
					close(chanStream)
					return
				}
				cs.ChatStremResGenerate(ctx, reqnew, chanStream)
				return
			}
			if response.Choices[0].FinishReason == "stop" {
				logger.Debugf("chat请求ID：%s", response.ID)
				// lastMessage.Content = resmessage
				// lastMessage.Role = openai.ChatMessageRoleAssistant
				// cs.ChatCostCalculate(ctx, []openai.ChatCompletionMessage{lastMessage}, req.Model)
				stream.Close()
				close(chanStream)
				return
			}
			if response.Choices[0].FinishReason == "content_filter" {
				logger.Debugf("chat请求ID：%s", response.ID)
				// lastMessage.Content = resmessage
				// lastMessage.Role = openai.ChatMessageRoleAssistant
				// cs.ChatCostCalculate(ctx, []openai.ChatCompletionMessage{lastMessage}, req.Model)
				chanStream <- "[content_filter]"
				stream.Close()
				close(chanStream)
				return
			}
		}
		if errors.Is(err, io.EOF) {
			// lastMessage.Content = resmessage
			// lastMessage.Role = openai.ChatMessageRoleAssistant
			// cs.ChatCostCalculate(ctx, []openai.ChatCompletionMessage{lastMessage}, req.Model)
			logger.Info("Stream finished")
			stream.Close()
			close(chanStream)
			return
		}
		if err != nil {
			logger.Errorf("Stream error: %v\n", err)
			// lastMessage.Content = resmessage
			// lastMessage.Role = openai.ChatMessageRoleAssistant
			// cs.ChatCostCalculate(ctx, []openai.ChatCompletionMessage{lastMessage}, req.Model)
			stream.Close()
			close(chanStream)
			return
		}
		select {
		case <-ctx.Writer.CloseNotify():
			// lastMessage.Content = resmessage
			// lastMessage.Role = openai.ChatMessageRoleAssistant
			// cs.ChatCostCalculate(ctx, []openai.ChatCompletionMessage{lastMessage}, req.Model)
			stream.Close()
			close(chanStream)
			return
		default:
			chanStream <- response.Choices[0].Delta.Content
			logger.Debugf(response.Choices[0].Delta.Content)
			resmessage += response.Choices[0].Delta.Content
		}
	}
}

func (cs *chatService) ChatEmbeddingGenerate(str []string) (embedVectors []openai.Embedding, err error) {
	var req openai.EmbeddingRequest
	req.Model = openai.AdaEmbeddingV2
	req.Input = str
	client, err := openai.NewClient()
	if err != nil {
		return
	}
	resp, err := client.CreateEmbeddings(req)
	if err != nil {
		logger.Errorf("Embeddings error: %v\n", err)
		return
	}
	// tokens = resp.Usage.TotalTokens
	embedVectors = resp.Data
	return
}

func (cs *chatService) ChatEmbeddingSave(ctx context.Context, title, body, classify string, embeddata openai.Embedding) error {
	docs := entity.Documents{}
	docs.Id = cs.iSrv.GenSnowID()
	docs.Title = title
	docs.Classify = classify
	// docs.Subsection = sub
	docs.Body = body
	docs.Tokens = tiktoken.NumTokensSingleString(body)
	docs.Embedding = pgvector.NewVector(embeddata.Embedding)
	return cs.cd.DocEmbeddingSave(ctx, &docs)
}

func (cs *chatService) ChatEmbeddingCompare(ctx context.Context, question, classify string) (contextStr string, err error) {
	//获取question Embedding信息
	embedvectors, err := cs.ChatEmbeddingGenerate([]string{question})
	if err != nil {
		return
	}
	textbody, err := cs.cd.ChatEmbeddingCompare(ctx, pgvector.NewVector(embedvectors[0].Embedding), classify)
	if len(textbody) != 0 {
		for _, v := range textbody {
			contextStr += v.Body
		}
		return
	}
	return
}

// func (cs *chatService) ChatTest(ctx context.Context, text string) (keyword string) {
// 	return cs.jieba.GetKeyword(text)
// }
