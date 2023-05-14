/*
 * @Author: cloudyi.li
 * @Date: 2023-03-29 10:51:03
 * @LastEditTime: 2023-05-12 23:01:21
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/consts/consts.go
 */
package consts

type APIType string

const (
	// RequestId 请求id名称
	RequestId = "request_id"
	// UserID 用户id key
	UserID   = "user_id"
	ChatID   = "chat_id"
	Balance  = "balance_ctx"
	TokenCtx = "token_ctx"
	CBCKEY   = "ABCDABCDABCDABCD"
	// EmbedCtx = "with_emebedding_ctx"
	// TimeLayout 时间格式
	TimeLayout                     = "2006-01-02 15:04:05"
	TimeLayoutMs                   = "2006-01-02 15:04:05.000"
	ApiURLv1                       = "https://api.openai.com/v1"
	DefaultEmptyMessagesLimit uint = 300
	OpenaiAPIURLv1                 = "https://api.openai.com/v1"
	AzureAPIPrefix                 = "openai"
	AzureDeploymentsPrefix         = "deployments"
	AvatarSize                     = 24
	TokenPrice                     = 0.00007

	APITypeOpenAI  APIType = "OPEN_AI"
	APITypeAzure   APIType = "AZURE"
	APITypeAzureAD APIType = "AZURE_AD"

	AzureAPIKeyHeader = "api-key"
)
