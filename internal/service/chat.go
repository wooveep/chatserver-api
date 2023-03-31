/*
 * @Author: cloudyi.li
 * @Date: 2023-03-29 13:45:51
 * @LastEditTime: 2023-03-31 17:23:07
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/service/chat.go
 */
package service

import (
	"chatserver-api/di/logger"
	"chatserver-api/pkg/openai"
	"errors"
	"io"
)

var _ ChatService = (*chatService)(nil)

type ChatService interface {
	GetChatResponse(chatMessage []openai.ChatCompletionMessage, closeWorker <-chan bool, chanStream chan<- string)
}

// userService 实现UserService接口
type chatService struct {
}

func NewChatService() *chatService {
	return &chatService{}
}

func (cs *chatService) reqChatCompletion(chatMessage []openai.ChatCompletionMessage) *openai.ChatCompletionStream {
	client := openai.NewClient()
	req := openai.ChatCompletionRequest{
		Model:     openai.GPT3Dot5Turbo,
		MaxTokens: 200,
		Messages:  chatMessage,
		Stream:    true,
	}
	stream, err := client.CreateChatCompletionStream(req)
	if err != nil {
		logger.Errorf("ChatCompletionStream error: %v\n", err)
	}
	return stream
}

func (cs *chatService) GetChatResponse(chatMessage []openai.ChatCompletionMessage, closeWorker <-chan bool, chanStream chan<- string) {
	defer close(chanStream)
	stream := cs.reqChatCompletion(chatMessage)
	defer stream.Close()

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
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
			logger.Infof("MessageGen:%s", response.Choices[0].Delta.Content)
		}

	}
}
