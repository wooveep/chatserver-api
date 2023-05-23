/*
 * @Author: cloudyi.li
 * @Date: 2023-03-30 18:16:24
 * @LastEditTime: 2023-05-23 10:48:47
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

	BaseURL              string
	OrgID                string
	APIType              consts.APIType
	APIVersion           string                    // required when APIType is APITypeAzure or APITypeAzureAD
	AzureModelMapperFunc func(model string) string // replace model to azure deployment name func
	HTTPClient           *http.Client

	EmptyMessagesLimit uint
}

func DefaultConfig(config config.OpenAIConfig) ClientConfig {
	return ClientConfig{
		HTTPClient: &http.Client{},
		BaseURL:    consts.OpenaiAPIURLv1,
		OrgID:      config.OrgID,
		authToken:  config.AuthToken,
		APIType:    consts.APITypeOpenAI,

		EmptyMessagesLimit: consts.DefaultEmptyMessagesLimit,
	}
}
func DefaultAzureConfig(config config.OpenAIConfig) ClientConfig {
	return ClientConfig{
		authToken:  config.AuthToken,
		BaseURL:    config.APIURL,
		OrgID:      "",
		APIType:    consts.APITypeAzure,
		APIVersion: config.APIVersion,
		AzureModelMapperFunc: func(model string) string {
			return consts.AzureToModel[model]
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
