/*
 * @Author: cloudyi.li
 * @Date: 2023-03-30 18:16:24
 * @LastEditTime: 2023-03-30 19:04:47
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/openai/config.go
 */
package openai

import (
	"chatserver-api/di/config"
	"chatserver-api/internal/constant"
	"net/http"
)

// ClientConfig is a configuration of a client.
type ClientConfig struct {
	authToken string

	HTTPClient *http.Client

	BaseURL string
	OrgID   string

	EmptyMessagesLimit uint
}

func DefaultConfig() ClientConfig {
	return ClientConfig{
		HTTPClient:         &http.Client{},
		BaseURL:            constant.ApiURLv1,
		OrgID:              config.AppConfig.OpenAIConfig.OrgID,
		authToken:          config.AppConfig.OpenAIConfig.AuthToken,
		EmptyMessagesLimit: constant.DefaultEmptyMessagesLimit,
	}
}
