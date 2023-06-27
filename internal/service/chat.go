/*
 * @Author: cloudyi.li
 * @Date: 2023-03-29 13:45:51
 * @LastEditTime: 2023-06-27 15:07:08
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
	"chatserver-api/pkg/chatfunc"
	"chatserver-api/pkg/logger"
	"chatserver-api/pkg/openai"
	"chatserver-api/pkg/pgvector"
	"chatserver-api/pkg/tiktoken"
	"chatserver-api/utils/security"
	"chatserver-api/utils/uuid"
	"fmt"
	"net/url"
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
	ChatRecordDelete(ctx *gin.Context, msgId int64) (err error)
	ChatBalanceVerify(ctx *gin.Context) (err error)
	ChatBalanceUpdate(ctx *gin.Context) (err error)
	ChatMessageSave(ctx *gin.Context, role, message string, msgid int64, mtype int) error
	ChatDetailGet(ctx *gin.Context) (res model.ChatDetailRes, err error)
	ChatUserVerify(ctx *gin.Context) (err error)
	ChatRegenerategReqProcess(ctx *gin.Context, msgid int64, memoryLevel int16) (answerid int64, req openai.ChatCompletionRequest, err error)
	ChatChattingReqProcess(ctx *gin.Context, lastquestion string, memoryLevel int16) (questionId int64, req openai.ChatCompletionRequest, err error)
	ChatReqUnit(ctx *gin.Context, lastquestion string, records []model.RecordOne) (req openai.ChatCompletionRequest, err error)
	ChatStremResGenerate(ctx *gin.Context, retry int, req openai.ChatCompletionRequest, chanStream chan<- string, closeNotify <-chan bool)
	ChatStreamResProcess(ctx *gin.Context, chanStream <-chan string, questionId, answerid int64) (msgid int64, messages string, flag int)
	ChatEmbeddingSave(ctx context.Context, title, body, classify string, embeddata openai.Embedding) error
	ChatEmbeddingGenerate(str []string) (embedVectors []openai.Embedding, err error)
	ChatEmbeddingCompare(ctx context.Context, question, classify string) (contextStr string, err error)
	ChatCostCalculate(ctx *gin.Context, promptMsgs []openai.ChatCompletionMessage, model string)
	ChatSearchExtension(ctx *gin.Context, question string) (result string)
	ChatFuncCallSave(ctx *gin.Context, role, message string, mtype int) error
}

// userService 实现UserService接口
type chatService struct {
	cd   dao.ChatDao
	uSrv UserService
	rc   *redis.Client
	// jieba tokenize.Tokenizer
	iSrv uuid.SnowNode
}

// func NewChatService(_cd dao.ChatDao, _uSrv UserService, _jieba tokenize.Tokenizer) *chatService {
func NewChatService(_cd dao.ChatDao, _uSrv UserService) *chatService {

	return &chatService{
		cd:   _cd,
		uSrv: _uSrv,
		iSrv: *uuid.NewNode(1),
		rc:   cache.GetRedisClient(),
		// jieba: _jieba,
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
	cs.rc.SAdd(ctx, consts.UserChatIDPrefix+strconv.FormatInt(chat.Id, 10))
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
		val3, _ := cs.rc.SMembers(ctx, consts.UserChatIDPrefix+strconv.FormatInt(userId, 10)).Result()
		for _, v := range val3 {
			cs.rc.Del(ctx, consts.ChatRecordIDPrefix+v)
			cs.rc.Del(ctx, consts.ChatSearchPrefix+v)
		}
		cs.rc.Del(ctx, consts.UserChatIDPrefix+strconv.FormatInt(userId, 10))
		return cs.cd.ChatDeleteAll(ctx, userId)
	} else {
		cs.rc.SRem(ctx, consts.UserChatIDPrefix+strconv.FormatInt(userId, 10), chatId)
		cs.rc.Del(ctx, consts.ChatRecordIDPrefix+strconv.FormatInt(chatId, 10))
		cs.rc.Del(ctx, consts.ChatSearchPrefix+strconv.FormatInt(chatId, 10))
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
	cs.rc.Del(ctx, consts.ChatSearchPrefix+strconv.FormatInt(chatId, 10))
	return
}

func (cs *chatService) ChatRecordDelete(ctx *gin.Context, msgId int64) (err error) {
	chatId := ctx.GetInt64(consts.ChatID)
	err = cs.cd.ChatRecordDelete(ctx, msgId)
	cs.rc.SRem(ctx, consts.ChatRecordIDPrefix+strconv.FormatInt(chatId, 10), msgId)
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
	priceratio := ctx.GetInt(consts.PriceRatioCtx)
	cost := float64(token) * consts.TokenPrice * float64(priceratio)
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

func (cs *chatService) ChatMessageSave(ctx *gin.Context, role, message string, msgid int64, mtype int) (err error) {
	chatId := ctx.GetInt64(consts.ChatID)
	var exist bool
	n, err := cs.rc.Exists(ctx, consts.ChatRecordIDPrefix+strconv.FormatInt(chatId, 10)).Result()
	if n > 0 {
		exist, err = cs.rc.SIsMember(ctx, consts.ChatRecordIDPrefix+strconv.FormatInt(chatId, 10), msgid).Result()
	}
	if err != nil || n <= 0 {
		recordlist, err := cs.cd.ChatRecordIdGet(ctx, chatId)
		if err != nil {
			return err
		}
		if len(recordlist) > 0 {
			for _, v := range recordlist {
				cs.rc.SAdd(ctx, consts.ChatRecordIDPrefix+strconv.FormatInt(chatId, 10), v)
				if msgid == v {
					exist = true
				}
			}
		} else {
			exist = false
		}
	}
	record := entity.Record{}
	record.Id = msgid
	record.ChatId = chatId
	record.Sender = role
	record.Message = message
	if message != "" {
		record.MessageHash = security.Md5(message)
		record.MessageToken = tiktoken.NumTokensSingleString(message)
	}
	switch mtype {
	case 1:
		record.IsCall = true
	case 2:
		record.IsFunc = true
	case 3:
		record.HasCall = true
	default:
	}
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
	chatId := ctx.GetInt64(consts.ChatID)
	var lastquestion string
	var new_record []model.RecordOne
	records, answerid, err := cs.cd.ChatRegenRecordGet(ctx, chatId, msgid, memoryLevel)
	logger.Debugf("answerid:%d", answerid)
	if err != nil {
		logger.Errorf("获取会话消息记录失败: %v\n", err)
		return
	}
	reclen := len(records)

	switch reclen {
	case 0:
		lastquestion = " "
	case 1:
		lastquestion = records[0].Message
		// records = []model.RecordOne{}
	default:
		lastquestion = records[reclen-1].Message
		newlen := reclen - 1
		for i := 0; i < newlen; i++ {
			new_record = append(new_record, records[i])
			var func_content []model.RecordOne
			var end int64
			// var start int64
			if records[i].HasCall {
				end = -1
				for j := i + 1; j < newlen; j++ {
					if records[j].Sender == openai.ChatMessageRoleUser {
						end = records[j].Id
						break
					}
				}
				func_content, _ = cs.cd.ChatFuncHisGet(ctx, chatId, records[i].Id, end)
				new_record = append(new_record, func_content...)
			}
		}

	}
	req, err = cs.ChatReqUnit(ctx, lastquestion, new_record)
	if err != nil {
		logger.Errorf("%v\n", err)
		return
	}
	return
}

func (cs *chatService) ChatChattingReqProcess(ctx *gin.Context, lastquestion string, memoryLevel int16) (questionId int64, req openai.ChatCompletionRequest, err error) {
	chatId := ctx.GetInt64(consts.ChatID)
	var new_record []model.RecordOne
	records, err := cs.cd.ChatRecordGet(ctx, chatId, memoryLevel)
	if err != nil {
		logger.Errorf("获取会话消息记录失败: %v\n", err)
		return
	}
	reclen := len(records)
	for i := 0; i < reclen; i++ {
		new_record = append(new_record, records[i])
		var func_content []model.RecordOne
		var end int64
		// var start int64
		if records[i].HasCall {
			end = -1
			for j := i + 1; j < reclen; j++ {
				if records[j].Sender == openai.ChatMessageRoleUser {
					end = records[j].Id
					break
				}
			}
			func_content, _ = cs.cd.ChatFuncHisGet(ctx, chatId, records[i].Id, end)

			new_record = append(new_record, func_content...)
		}
	}

	req, err = cs.ChatReqUnit(ctx, lastquestion, new_record)

	if err != nil {
		logger.Errorf("%v\n", err)
		return
	}
	questionId = cs.iSrv.GenSnowID()
	return
}

func (cs *chatService) ChatReqUnit(ctx *gin.Context, lastquestion string, records []model.RecordOne) (req openai.ChatCompletionRequest, err error) {
	var chatMessages []openai.ChatCompletionMessage
	var systemPreset, historyMessage, lastMessage openai.ChatCompletionMessage
	var logitbia map[string]int
	var emquestion, embedcontexts, messagePrefix string
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
	req.Model = preset.ModelName
	req.Stream = true
	req.MaxTokens = preset.MaxTokens
	req.LogitBias = logitbia
	req.Temperature = float32(preset.Temperature)
	req.TopP = float32(preset.TopP)
	req.FrequencyPenalty = float32(preset.Frequency)
	req.PresencePenalty = float32(preset.Presence)
	req.User = strconv.FormatInt(ctx.GetInt64(consts.UserID), 10)
	ctx.Set(consts.PriceRatioCtx, 1)
	systemPreset.Role = openai.ChatMessageRoleSystem
	switch preset.Extension {
	// 默认智能助手
	case 1:
		{

			systemPreset.Content = strings.Replace(preset.PresetContent, "{{ current_date }}", time.Now().Local().Format(consts.TimeLayout), -1)
		}
	// 联网搜索
	case 2:
		{
			// lastquestion 查询拼接
			ctx.Set(consts.PriceRatioCtx, 5)
			if req.Model == openai.GPT3Dot5Turbo0613 || req.Model == openai.GPT3Dot5Turbo16K0613 {
				req.Functions = []*openai.FunctionDefine{&chatfunc.FuncGetWeather, &chatfunc.FuncEntitySearch, &chatfunc.FuncGoogleSearch}
				systemPreset.Content = strings.Replace(preset.PresetContent, "{{ current_date }}", time.Now().Format(consts.TimeLayout), -1)
			} else {
				searchContext := cs.ChatSearchExtension(ctx, lastquestion)
				systemPreset.Content = strings.Replace(preset.PresetContent, "{{ current_date }}", time.Now().Format(consts.DateLayout), -1) + "\n## Contexts:\n" + searchContext
			}
		}
	//embedding 数据
	case 3:
		{
			if len(records) != 0 {
				for _, v := range records {
					if v.Sender == "user" {
						emquestion += v.Message
					}
				}
			}
			emquestion += lastquestion
			//将用户问题进行关键词提取
			// emkeyword := cs.jieba.GetKeyword(emquestion) + lastquestion
			//通过用户问题 lastquestion + records（User历史）获取Context信息
			// embedcontexts, err = cs.ChatEmbeddingCompare(ctx, emkeyword, preset.Classify)
			embedcontexts, err = cs.ChatEmbeddingCompare(ctx, lastquestion, preset.Classify)
			if err != nil {
				logger.Errorf("获取embedding上下文失败: %v\n", err)
				return
			}
			//替换拼接PresetContent
			systemPreset.Content = strings.Replace(preset.PresetContent, "{{ context }}", embedcontexts, -1)
		}
	// 翻译助手
	case 4:
		{
			records = []model.RecordOne{}
			messagePrefix = "Translate:  "
			systemPreset.Content = preset.PresetContent
		}
	default:

		systemPreset.Content = preset.PresetContent
	}
	chatMessages = append(chatMessages, systemPreset)
	for _, record := range records {
		historyMessage = openai.ChatCompletionMessage{}
		if record.IsFunc {
			historyMessage.Role = openai.ChatMessageRoleFunction
			historyMessage.Name = record.Sender
			historyMessage.Content = record.Message
		} else if record.IsCall {
			historyMessage.Role = openai.ChatMessageRoleAssistant
			historyMessage.Content = ""
			funccall := &openai.FunctionCall{Name: record.Sender, Arguments: record.Message}
			historyMessage.FunctionCall = funccall
		} else {
			historyMessage.Role = record.Sender
			historyMessage.Content = record.Message
		}
		logger.Debugf("ROLE:%s,Message:%s", record.Sender, record.Message)
		chatMessages = append(chatMessages, historyMessage)
	}
	lastMessage.Role = openai.ChatMessageRoleUser
	lastMessage.Content = messagePrefix + lastquestion
	chatMessages = append(chatMessages, lastMessage)
	req.Messages = chatMessages
	cs.ChatCostCalculate(ctx, chatMessages[1:], preset.ModelName)
	return
}
func (cs *chatService) ChatStreamResProcess(ctx *gin.Context, chanStream <-chan string, questionId, answerid int64) (msgid int64, messages string, flag int) {

	msgtime := time.Now().Format(consts.TimeLayout)
	if answerid != 0 {
		msgid = answerid
	} else {
		msgid = cs.iSrv.GenSnowID()
	}
	ctx.Stream(func(w io.Writer) bool {
		if msg, ok := <-chanStream; ok {
			switch msg {
			case "[CONTENT_FILTER_PRE]":
				{
					err_messages := "尊敬的客户，非常感谢您使用我们的服务。我们注意到您最近提交的问题被我们的内容过滤器拦截了。我们深表歉意，因为我们的过滤器是为了保护我们的用户免受不良内容的侵害而设置的。但是，我们也理解您的问题对您来说非常重要。如果您有任何疑问或需要进一步的帮助，请随时联系网站管理员。再次感谢您的支持和理解。"
					flag = -1
					ctx.SSEvent("chatting", map[string]string{"question_id": strconv.FormatInt(questionId, 10), "msgid": strconv.FormatInt(msgid, 10), "time": msgtime, "text": err_messages})
					return false
				}
			case "[CONTENT_FILTER_BACK]":
				{
					err_messages := "尊敬的客户，非常感谢您使用我们的服务。我们注意到您最近提交的问题被我们的内容过滤器拦截了。我们深表歉意，因为我们的过滤器是为了保护我们的用户免受不良内容的侵害而设置的。但是，我们也理解您的问题对您来说非常重要。如果您有任何疑问或需要进一步的帮助，请随时联系网站管理员。再次感谢您的支持和理解。"
					flag = 2
					ctx.SSEvent("chatting", map[string]string{"question_id": strconv.FormatInt(questionId, 10), "msgid": strconv.FormatInt(msgid, 10), "time": msgtime, "text": err_messages})
					return false
				}
			case "[FUNCTION_CALL]":
				{
					flag = 1
					// msgid = cs.iSrv.GenSnowID()
					return true
				}
			case "[REQ_ERROR]":
				{
					err_messages := messages + "尊敬的客户，非常感谢您使用我们的服务。由于API暂时异常,我们深表歉意,请随时联系网站管理员。再次感谢您的支持和理解。"

					ctx.SSEvent("chatting", map[string]string{"question_id": strconv.FormatInt(questionId, 10), "msgid": strconv.FormatInt(msgid, 10), "time": msgtime, "text": err_messages})
					return false
				}
			default:
				{
					messages += msg
					ctx.SSEvent("chatting", map[string]string{"question_id": strconv.FormatInt(questionId, 10), "msgid": strconv.FormatInt(msgid, 10), "time": msgtime, "delta": messages})
					// logger.Debugf("stream-event: ID:%s ,time:%s,msg:%s", ctx.GetString(consts.RequestId), msgtime, msg)
					return true
				}
			}
		}

		ctx.SSEvent("chatting", map[string]string{"question_id": strconv.FormatInt(questionId, 10), "msgid": strconv.FormatInt(msgid, 10), "time": msgtime, "text": messages})
		return false
	})
	logger.Debugf("Stream-message:%s", messages)
	return
}

func (cs *chatService) ChatStremResGenerate(ctx *gin.Context, retry int, req openai.ChatCompletionRequest, chanStream chan<- string, closeNotify <-chan bool) {
	var chatMessages []openai.ChatCompletionMessage
	var lastMessage, blankMessage, funcMessage openai.ChatCompletionMessage
	var resmessage, funcname, funcargument string
	blankMessage.Content = "[cmd:continue]"
	// blankMessage.Content = " "
	blankMessage.Role = openai.ChatMessageRoleSystem
	chatMessages = req.Messages
	// cs.ChatCostCalculate(ctx, chatMessages[1:], req.Model)
	for _, v := range chatMessages {
		logger.Debugf("role: %s  , content: %s ", v.Role, v.Content)
	}
	client, err := openai.NewClient()
	if err != nil {
		close(chanStream)
		return
	}
	stream, err := client.CreateChatCompletionStream(req)
	if err != nil {
		if apierror, ok := err.(*openai.APIError); ok {
			switch apierror.HTTPStatusCode {
			case 429, 504, 500, 503:
				{
					if retry < 3 {
						logger.Warnf("ChatCompletionStream error重试: %v\n", err)
						retry += 1
						cs.ChatStremResGenerate(ctx, retry, req, chanStream, closeNotify)
					}
				}
			case 400:
				{
					if strings.Contains(apierror.Message, "content filtering") {
						chanStream <- "[CONTENT_FILTER_PRE]"
					} else {
						chanStream <- "[REQ_ERROR]"
					}
				}
			default:
				{
					chanStream <- "[REQ_ERROR]"
				}
			}
		} else if urlErr, ok := err.(*url.Error); ok {
			if urlErr.Timeout() && retry < 3 {
				logger.Warnf("ChatCompletionStream error重试: %v\n", err)
				retry += 1
				cs.ChatStremResGenerate(ctx, retry, req, chanStream, closeNotify)
			}
		}
		logger.Errorf("ChatCompletionStream error: %v\n", err)
		time.Sleep(10 * time.Millisecond)
		close(chanStream)
		return
	}
	for {
		if stream == nil {
			close(chanStream)
			return
		}
		response, err := stream.Recv()
		if len(response.Choices) > 0 {
			switch response.Choices[0].FinishReason {
			case openai.FinishReasonLength:
				{
					logger.Debugf("chat请求ID：%stoken截断", response.ID)
					stream.Close()
					lastMessage.Content = resmessage
					lastMessage.Role = openai.ChatMessageRoleAssistant
					chatMessages = append(chatMessages, lastMessage, blankMessage)
					req.Messages = chatMessages
					req.FunctionCall = "none"
					if consts.ModelMaxToken[req.Model] < tiktoken.NumTokensFromMessages(chatMessages, req.Model)+req.MaxTokens {
						close(chanStream)
						return
					}
					cs.ChatStremResGenerate(ctx, retry, req, chanStream, closeNotify)
					return
				}
			case openai.FinishReasonStop:
				{
					logger.Debugf("chat请求ID：%s正常结束", response.ID)
					stream.Close()
					close(chanStream)
					return
				}
			case openai.FinishReasonContentFilter:
				{
					logger.Debugf("chat请求ID：%s内容过滤", response.ID)
					chanStream <- "[CONTENT_FILTER_BACK]"
					stream.Close()
					time.Sleep(1 * time.Millisecond)
					close(chanStream)
					return
				}
			case openai.FinishReasonFunctionCall:
				{
					//
					chanStream <- "[FUNCTION_CALL]"
					stream.Close()
					logger.Debug("chat调用函数", logger.Pair("OpenAI_ID", response.ID), logger.Pair("funcname", funcname), logger.Pair("funccontent", funcargument))
					if err := cs.ChatFuncCallSave(ctx, funcname, funcargument, 1); err != nil {
						logger.Error("[消息保存]", logger.Pair("错误原因", err))
					}
					lastMessage.Role = openai.ChatMessageRoleAssistant
					lastMessage.Content = ""
					funcCall := &openai.FunctionCall{Name: funcname, Arguments: funcargument}
					lastMessage.FunctionCall = funcCall

					funcContent := chatfunc.ChatFuncProcess(ctx, funcname, funcargument)
					funcMessage.Content = funcContent
					funcMessage.Name = funcname
					funcMessage.Role = openai.ChatMessageRoleFunction
					if err := cs.ChatFuncCallSave(ctx, funcname, funcContent, 2); err != nil {
						logger.Error("[消息保存]", logger.Pair("错误原因", err))
					}
					chatMessages = append(chatMessages, lastMessage, funcMessage)
					req.Messages = chatMessages

					cs.ChatStremResGenerate(ctx, retry, req, chanStream, closeNotify)
					return
				}
			default:
				{
					if response.Choices[0].Delta.FunctionCall == nil {
						resmessage += response.Choices[0].Delta.Content
						logger.Debugf(response.Choices[0].Delta.Content)
						select {
						case <-closeNotify:
							stream.Close()
							close(chanStream)
							return
						default:
							chanStream <- response.Choices[0].Delta.Content
						}
					} else {
						if response.Choices[0].Delta.FunctionCall.Name != "" {
							funcname = response.Choices[0].Delta.FunctionCall.Name
						}
						funcargument += response.Choices[0].Delta.FunctionCall.Arguments
					}
				}
			}
		}
		if errors.Is(err, io.EOF) {
			logger.Info("Stream finished")
			stream.Close()
			close(chanStream)
			return
		}
		if err != nil {
			logger.Errorf("Stream error: %v\n", err)
			stream.Close()
			close(chanStream)
			return
		}
	}
}

func (cs *chatService) ChatFuncCallSave(ctx *gin.Context, role, message string, mtype int) error {
	chatId := ctx.GetInt64(consts.ChatID)
	record := entity.Record{}
	record.Id = cs.iSrv.GenSnowID()
	record.ChatId = chatId
	record.Sender = role
	record.Message = message
	record.MessageHash = security.Md5(message)
	record.MessageToken = tiktoken.NumTokensSingleString(message)
	switch mtype {
	case 1:
		record.IsCall = true
	case 2:
		record.IsFunc = true
	case 3:
		record.HasCall = true
	default:
	}
	return cs.cd.ChatRecordSave(ctx, &record)

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
			logger.Debugf(v.Body)
		}
		return
	}
	return
}

func (cs *chatService) ChatSearchExtension(ctx *gin.Context, question string) (result string) {
	chatId := ctx.GetInt64(consts.ChatID)

	result, err := chatfunc.CustomFuncExtension(ctx, question)
	if err != nil {
		logger.Warnf("搜索异常:%v", err.Error())
		return
	}
	if result == "" {
		result, err = cs.rc.Get(ctx, consts.ChatSearchPrefix+strconv.FormatInt(chatId, 10)).Result()
		if err == nil {
			return
		} else {
			if err != redis.Nil {
				logger.Errorf("Redis连接异常:%v", err.Error())
				return
			}
			logger.Debugf(" 缓存不存在:%v", err.Error())
			return
		}

	}
	// err = cs.rc.Set(ctx, consts.ChatSearchPrefix+security.Md5(question), result, 30*time.Minute).Err()
	err = cs.rc.Set(ctx, consts.ChatSearchPrefix+strconv.FormatInt(chatId, 10), result, 0).Err()
	if err != nil {
		logger.Errorf("Redis连接异常:%v", err.Error())
	}

	return
}
