/*
 * @Author: cloudyi.li
 * @Date: 2023-02-15 14:05:41
 * @LastEditTime: 2023-03-29 09:21:49
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/di/config/main.go
 */
package config

type Config struct {
	Log Log `toml:"log"`
	// Account   OfficialAccountConfig `toml:"account"`
	ApiServer ApiServerConfig `toml:"apiserver"`
	// Redis     RedisConfig           `toml:"redis"`
}

type Log struct {
	Level string `toml:"level"`
	File  []File `toml:"file"`
}

type File struct {
	Format  LogFormat `toml:"format"`
	Path    string    `toml:"path"`
	MaxSize int       `toml:"max_size"`
	MaxAge  int       `toml:"max_age"`
}

type LogFormat string

const (
	LogFormatJSON LogFormat = "json"
	LogFormatText LogFormat = "text"
)

type ApiServerConfig struct {
	Listen string `toml:"listen"`
}
