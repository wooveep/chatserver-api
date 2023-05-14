/*
 * @Author: cloudyi.li
 * @Date: 2023-03-30 18:16:24
 * @LastEditTime: 2023-05-12 23:16:17
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/openai/config.go
 */
package openai

import (
	"chatserver-api/internal/consts"
	"chatserver-api/pkg/config"
	"net/http"
	"regexp"
)

// ClientConfig is a configuration of a client.
type ClientConfig struct {
	authToken string

	BaseURL              string
	OrgID                string
	APIType              consts.APIType
	APIVersion           string                    // required when APIType is APITypeAzure or APITypeAzureAD
	AzureModelMapperFunc func(model string) string // replace model to azure deployment name func
	HTTPClient           *http.Client

	EmptyMessagesLimit uint
}

func DefaultConfig() ClientConfig {
	return ClientConfig{
		HTTPClient: &http.Client{},
		BaseURL:    consts.OpenaiAPIURLv1,
		OrgID:      config.AppConfig.OpenAIConfig.OrgID,
		authToken:  config.AppConfig.OpenAIConfig.AuthToken,
		APIType:    consts.APITypeOpenAI,

		EmptyMessagesLimit: consts.DefaultEmptyMessagesLimit,
	}
}
func DefaultAzureConfig(apiKey, baseURL string) ClientConfig {
	return ClientConfig{
		authToken:  config.AppConfig.OpenAIConfig.AuthToken,
		BaseURL:    baseURL,
		OrgID:      "",
		APIType:    consts.APITypeAzure,
		APIVersion: "2023-03-15-preview",
		AzureModelMapperFunc: func(model string) string {
			return regexp.MustCompile(`[.:]`).ReplaceAllString(model, "")
		},

		HTTPClient: &http.Client{},

		EmptyMessagesLimit: consts.DefaultEmptyMessagesLimit,
	}
}

func (ClientConfig) String() string {
	return "<OpenAI API ClientConfig>"
}

func (c ClientConfig) GetAzureDeploymentByModel(model string) string {
	if c.AzureModelMapperFunc != nil {
		return c.AzureModelMapperFunc(model)
	}

	return model
}
