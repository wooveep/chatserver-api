/*
 * @Author: cloudyi.li
 * @Date: 2023-02-15 16:22:00
 * @LastEditTime: 2023-04-05 15:55:12
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/config/config.go
 */
package config

import (
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var AppConfig *Config

// Load is a loader to load config file.
func Load(configFilePath string) *Config {
	resolveRealPath(configFilePath)
	// 初始化配置文件
	if err := initConfig(); err != nil {
		panic(err)
	}
	// 监控配置文件，并热加载
	watchConfig()

	return AppConfig
}

func initConfig() error {
	viper.SetConfigType("yml")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("APPLICATION")
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	// 解析配置文件
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	// 解析到struct
	AppConfig = &Config{}
	if err := viper.Unmarshal(AppConfig); err != nil {
		panic(err)
	}
	return nil
}

// 监控配置文件变动
// 注意：有些配置修改后，及时重新加载也要重新启动应用才行，比如端口
func watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		// 忽略错误
		Load(in.Name)
	})
}

// 如果未传递配置文件路径将使用约定的环境配置文件
func resolveRealPath(path string) {
	if path != "" {
		viper.SetConfigFile(path)
	} else {
		// 设置默认的config
		viper.AddConfigPath("configs")
		viper.SetConfigName("config")
	}
}
