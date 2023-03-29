/*
 * @Author: cloudyi.li
 * @Date: 2023-02-15 16:22:00
 * @LastEditTime: 2023-03-28 14:34:33
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/di/config/config.go
 */
package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var AppConfig Config

func InitConfig(configFile string) {
	if configFile == "" {
		configFile = fmt.Sprintf("configs/config.toml")
	}
	viper.SetConfigFile(configFile)
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s %s\n", err, configFile))
	}
	err = viper.Unmarshal(&AppConfig)
	if err != nil {
		panic(fmt.Errorf("Unable to decode into struct, %v", err))
	}
}
