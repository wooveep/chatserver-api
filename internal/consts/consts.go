/*
 * @Author: cloudyi.li
 * @Date: 2023-03-29 10:51:03
 * @LastEditTime: 2023-06-19 11:48:05
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/consts/consts.go
 */
package consts

type APIType string

const (
	// RequestId 请求id名称
	RequestId = "request_id"
	// UserID 用户id key
	UserID        = "user_id"
	ChatID        = "chat_id"
	RoleID        = "role_id"
	BalanceCtx    = "balance_ctx"
	CostTokenCtx  = "cost_token_ctx"
	JWTTokenCtx   = "token_ctx"
	PriceRatioCtx = "priceratio_ctx"

	InviteReward   = 3
	RegisterReward = 3
	CBCKEY         = "QLHDYGCODHTAPWEB"
	CDKEYBASE      = "XYLTNCP5H4FR8SJUVE2DWIK3MZ6B9QA7"
	InviteBase     = "zKYqDZXc8uvQA7mngtrsCPpThk53MJ2U6SR4HwfiEdeNjLB9WbIyxVaF"
	// EmbedCtx = "with_emebedding_ctx"
	// TimeLayout 时间格式
	DateLayout                     = "2006-01-02"
	TimeLayout                     = "2006-01-02 15:04:05"
	TimeLayoutMs                   = "2006-01-02 15:04:05.000"
	DefaultEmptyMessagesLimit uint = 300
	OpenaiAPIURLv1                 = "https://api.openai.com/v1"
	AzureAPIPrefix                 = "openai"
	AzureDeploymentsPrefix         = "deployments"
	AvatarSize                     = 24
	TokenPrice                     = 0.00015

	APITypeOpenAI  APIType = "OPEN_AI"
	APITypeAzure   APIType = "AZURE"
	APITypeAzureAD APIType = "AZURE_AD"

	AzureAPIKeyHeader    = "api-key"
	UserInvitePrefix     = "User_Invite_relation_list:"
	UserAvatarPrefix     = "User_Avatar_url_list:"
	UserInviteLinkPrefix = "User_Invite_Link_list:"
	UserInfoPrefix       = "User_Info_list:"
	UserBalancePrefix    = "User_Balance_list:"
	CaptchaPrefix        = "Captchat_list:"
	PresetPrefix         = "Preset_list:"
	GiftcardPrefix       = "GiftCard_list:"
	UserChatIDPrefix     = "User_ChatId_set:"
	ChatRecordIDPrefix   = "Chat_RecordId_set:"
	ChatSearchPrefix     = "Chat_Search_list:"
	// QuerySearchPrefix    = "Query_Search_list:"
	SearchCachePrefix = "Search_Cache_list:"
)

var AzureToModel = map[string]string{
	"gpt-3.5-turbo":          "gpt3",
	"text-davinci-003":       "davinci",
	"text-embedding-ada-002": "embedding",
}

var ModelMaxToken = map[string]int{
	"gpt-3.5-turbo":          4096,
	"gpt-3.5-turbo-16k":      16384,
	"gpt-3.5-turbo-16k-0613": 16384,
}

const (
	StandardUser = iota + 1
	RegularMembers
	SeniorMember
	InfiniteMember
	Administrator = 100
	Enterprise    = 301
)

var RoleToString = map[int]string{
	StandardUser:   "普通用户",
	RegularMembers: "标准会员",
	SeniorMember:   "高级会员",
	InfiniteMember: "无限会员",
	Enterprise:     "企业订阅",
	Administrator:  "管理员",
}
