/*
 * @Author: cloudyi.li
 * @Date: 2023-03-29 13:45:51
 * @LastEditTime: 2023-04-05 15:56:01
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/service/chat.go
 */
package service

import (
	"chatserver-api/internal/consts"
	"chatserver-api/internal/dao"
	"chatserver-api/internal/model"
	"chatserver-api/pkg/logger"
	"chatserver-api/pkg/openai"
	"chatserver-api/pkg/response"
	"errors"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

var _ ChatService = (*chatService)(nil)

type ChatService interface {
	GetChatResponse(chatMessage []openai.ChatCompletionMessage, closeWorker <-chan bool, chanStream chan<- string)
	REQMessageProcess(usermesage model.ChatChattingReq) (chatMessages []openai.ChatCompletionMessage)
	RESMessageProcess(ctx *gin.Context, chanStream <-chan string) (messages string, err error)
}

// userService 实现UserService接口
type chatService struct {
	cd dao.ChatDao
}

func NewChatService(_cd dao.ChatDao) *chatService {
	return &chatService{
		cd: _cd,
	}
}

func (cs *chatService) reqChatCompletion(chatMessage []openai.ChatCompletionMessage) (*openai.ChatCompletionStream, error) {
	client, err := openai.NewClient()
	if err != nil {
		return nil, err
	}
	req := openai.ChatCompletionRequest{
		Model:     openai.GPT3Dot5Turbo,
		MaxTokens: 200,
		Messages:  chatMessage,
		Stream:    true,
	}
	return client.CreateChatCompletionStream(req)
}

func (cs *chatService) REQMessageProcess(usermesage model.ChatChattingReq) (chatMessages []openai.ChatCompletionMessage) {
	var chatmessage openai.ChatCompletionMessage
	chatmessage.Role = openai.ChatMessageRoleUser
	chatmessage.Content = usermesage.Message
	return append(chatMessages, chatmessage)
}

func (cs *chatService) RESMessageProcess(ctx *gin.Context, chanStream <-chan string) (messages string, err error) {

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
	return
}

func (cs *chatService) GetChatResponse(chatMessage []openai.ChatCompletionMessage, closeWorker <-chan bool, chanStream chan<- string) {
	defer close(chanStream)
	stream, err := cs.reqChatCompletion(chatMessage)
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
			logger.Debugf("MessageGen:%s", response.Choices[0].Delta.Content)
		}

	}
}
