/*
 * @Author: cloudyi.li
 * @Date: 2023-03-29 10:51:03
 * @LastEditTime: 2023-05-24 11:10:25
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/consts/consts.go
 */
package consts

type APIType string

const (
	// RequestId 请求id名称
	RequestId = "request_id"
	// UserID 用户id key
	UserID       = "user_id"
	ChatID       = "chat_id"
	Balance      = "balance_ctx"
	TokenCtx     = "token_ctx"
	InviteReward = 5
	CBCKEY       = "ABCDABCDABCDABCD"
	CDKEYBASE    = "E8S2DZX9WYLTN6BQA7CP5IK3MJUFR4HV"
	InviteBase   = "E8uvS2pqDZXbcde9WYfiLTNrs6BxQA7CPmn5IyzK3MwJUktFghR4HVaj"
	// EmbedCtx = "with_emebedding_ctx"
	// TimeLayout 时间格式
	TimeLayout                     = "2006-01-02 15:04:05"
	TimeLayoutMs                   = "2006-01-02 15:04:05.000"
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
	UserInvitePrefix  = "User_Invite_relation_list:"
	UserAvatarPrefix  = "User_Avatar_url_list:"
	CaptchaPrefix     = "Captchat_list:"
)

var AzureToModel = map[string]string{
	"gpt-3.5-turbo":          "gpt3",
	"text-davinci-003":       "davinci",
	"text-embedding-ada-002": "embedding",
}

const (
	StandardUser = iota + 1
	RegularMembers
	SeniorMember
	InfiniteMember
	Enterprise
	Administrator = 100
)

var RoleToString = map[int]string{
	StandardUser:   "普通用户",
	RegularMembers: "标准会员",
	SeniorMember:   "高级会员",
	InfiniteMember: "无限会员",
	Enterprise:     "企业订阅",
	Administrator:  "管理员",
}
