/*
 * @Author: cloudyi.li
 * @Date: 2023-02-15 14:05:41
 * @LastEditTime: 2023-06-08 11:29:23
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/config/conf.go
 */
package config

type Config struct {
	Mode          string        `mapstructure:"mode"`           // gin启动模式
	Port          string        `mapstructure:"port"`           // 启动端口
	AppName       string        `mapstructure:"app-name"`       //应用名称
	Url           string        `mapstructure:"url"`            // 应用地址,用于自检 eg. http://127.0.0.1
	MaxPingCount  int           `mapstructure:"max-ping-count"` // 最大自检次数，用户健康检查
	Language      string        `mapstructure:"language"`       // 项目语言
	ExternalURL   string        `mapstructure:"externalurl"`
	JwtConfig     JwtConfig     `mapstructure:"jwt"`
	OpenAIConfig  OpenAIConfig  `mapstructure:"openai"`
	EmailCofig    EmailCofig    `mapstructure:"email"`
	DBConfig      DBConfig      `mapstructure:"database"` // 数据库信息
	RedisConfig   RedisConfig   `mapstructure:"redis"`    // redis
	LogConfig     LogConfig     `mapstructure:"log"`      // uber z
	CustomConfig  CustomConfig  `mapstructure:"custom"`
	TencentConfig TencentConfig `mapstructure:"tencent"`
	GoogelConfig  GoogelConfig  `mapstructure:"google"`
}

type JwtConfig struct {
	Secret                  string `mapstructure:"secret"`
	JwtTtl                  int64  `mapstructure:"ttl"`             // token 有效期（秒）
	JwtBlacklistGracePeriod int64  `mapstructure:"blacklistperiod"` // 黑名单宽限时间（秒）
}

type TencentConfig struct {
	SecretId  string `mapstructure:"secretid"`
	SecretKey string `mapstructure:"secretkey"`
}
type GoogelConfig struct {
	ApiKey string `mapstructure:"apikey"`
	CxId   string `mapstructure:"cxid"`
}
type OpenAIConfig struct {
	APIType    string `mapstructure:"apitype"`
	APIURL     string `mapstructure:"apiurl"`
	APIVersion string `mapstructure:"apiversion"`
	AuthToken  string `mapstructure:"authtoken"`
	OrgID      string `mapstructure:"orgid"`
	ProxyMode  string `mapstructure:"proxymode"`
	ProxyIP    string `mapstructure:"proxyip"`
	ProxyPort  string `mapstructure:"proxyport"`
}

// DBConfig is used to configure mysql database
type DBConfig struct {
	Dbname          string `mapstructure:"dbname"`
	Host            string `mapstructure:"host"`
	Port            string `mapstructure:"port"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	MaximumPoolSize int    `mapstructure:"maximum-pool-size"`
	MaximumIdleSize int    `mapstructure:"maximum-idle-size"`
	LogMode         bool   `mapstructure:"log-mode"`
}

// RedisConfig is used to configure redis
type RedisConfig struct {
	Addr         string `mapstructure:"address"`
	Password     string `mapstructure:"password"`
	Db           int    `mapstructure:"db"`
	PoolSize     int    `mapstructure:"pool-size"`
	MinIdleConns int    `mapstructure:"min-idle-conns"`
	IdleTimeout  int    `mapstructure:"idle-timeout"`
}

// LogConfig is used to configure uber zap
type LogConfig struct {
	Level      string `mapstructure:"level"`
	FileName   string `mapstructure:"file-name"`
	TimeFormat string `mapstructure:"time-format"`
	MaxSize    int    `mapstructure:"max-size"`
	MaxBackups int    `mapstructure:"max-backups"`
	MaxAge     int    `mapstructure:"max-age"`
	Compress   bool   `mapstructure:"compress"`
	LocalTime  bool   `mapstructure:"local-time"`
	Console    bool   `mapstructure:"console"`
}

type EmailCofig struct {
	Host      string `mapstructure:"smtphost"`
	Port      string `mapstructure:"smtpport"`
	Username  string `mapstructure:"smtpuser"`
	Password  string `mapstructure:"smtppassword"`
	Sender    string `mapstructure:"smtpsender"`
	ProxyMode string `mapstructure:"smtphost"`
	ProxyIP   string `mapstructure:"proxyip"`
	ProxyPort string `mapstructure:"proxyport"`
	PreCheck  bool   `mapstructure:"precheck"`
}

type CustomConfig struct {
	AzureProxy string `mapstructure:"azureproxy"`
}
