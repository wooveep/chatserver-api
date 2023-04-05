/*
 * @Author: cloudyi.li
 * @Date: 2023-03-30 18:16:24
 * @LastEditTime: 2023-04-05 15:55:51
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/openai/config.go
 */
package openai

import (
	"chatserver-api/internal/consts"
	"chatserver-api/pkg/config"
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
		BaseURL:            consts.ApiURLv1,
		OrgID:              config.AppConfig.OpenAIConfig.OrgID,
		authToken:          config.AppConfig.OpenAIConfig.AuthToken,
		EmptyMessagesLimit: consts.DefaultEmptyMessagesLimit,
	}
}
