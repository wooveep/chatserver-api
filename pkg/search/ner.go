/*
 * @Author: cloudyi.li
 * @Date: 2023-06-06 11:23:44
 * @LastEditTime: 2023-06-08 22:37:28
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/search/ner.go
 */
package search

// 使用腾讯云API检测进行NER实体检测
import (
	"chatserver-api/pkg/config"
	"chatserver-api/pkg/logger"
	"encoding/json"
	"strings"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	nlp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/nlp/v20190408"
)

type nlpResponse struct {
	Response struct {
		NormalText       string `json:"NormalText"`
		BasicParticiples []struct {
			Word        string `json:"Word"`
			BeginOffset int    `json:"BeginOffset"`
			Length      int    `json:"Length"`
			Pos         string `json:"Pos"`
		} `json:"BasicParticiples,omitempty"`
		CompoundParticiples []struct {
			Word        string `json:"Word"`
			BeginOffset int    `json:"BeginOffset"`
			Length      int    `json:"Length"`
			Pos         string `json:"Pos"`
		} `json:"CompoundParticiples,omitempty"`
		Entities []struct {
			Word        string `json:"Word"`
			BeginOffset int    `json:"BeginOffset"`
			Length      int    `json:"Length"`
			Type        string `json:"Type"`
			Name        string `json:"Name"`
		} `json:"Entities,omitempty"`
		RequestID string `json:"RequestId"`
	} `json:"Response"`
}

func genin(target string, str_array []string) (count int) {

	for _, v := range str_array {
		if v == target {
			count += 1
		}
	}
	return
}

func nerDetec(query string) (int, string) {
	var nlpres nlpResponse
	var participles []string
	tencentcfg := config.AppConfig.TencentConfig

	credential := common.NewCredential(
		tencentcfg.SecretId,
		tencentcfg.SecretKey,
	)
	// 实例化一个client选项，可选的，没有特殊需求可以跳过
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "nlp.tencentcloudapi.com"
	// 实例化要请求产品的client对象,clientProfile是可选的
	client, _ := nlp.NewClient(credential, "", cpf)

	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := nlp.NewParseWordsRequest()

	request.Text = common.StringPtr(query)

	// 返回的resp是一个ParseWordsResponse的实例，与请求对象对应
	response, err := client.ParseWords(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		logger.Errorf("An API error has returned: %s", err)
		return 0, ""
	}
	if err != nil {
		logger.Errorf("%s", err)
		return 0, ""
	}
	// 输出json格式的字符串回包
	json.Unmarshal([]byte(response.ToJsonString()), &nlpres)
	for _, v := range nlpres.Response.CompoundParticiples {

		logger.Debug(v.Pos + "||" + v.Word)
	}
	var generic []string
	if len(nlpres.Response.Entities) > 0 {
		for _, v := range nlpres.Response.Entities {
			logger.Debug(v.Type + "||" + v.Word)
			generic = append(generic, v.Type)
			if v.Type == "loc.generic" || v.Type == "org.generic" || v.Type == "medicine" || v.Type == "event.generic" || v.Type == "product.generic" || v.Type == "person.generic" {
				participles = append(participles, v.Word)
			}
		}
	} else {
		return 0, ""
	}

	if len(nlpres.Response.Entities) == genin("quantity.generic", generic) {
		return 0, ""
	}
	if genin("time.generic", generic) > 0 {
		return 2, strings.Join(participles, " ")
	}
	return 1, strings.Join(participles, " ")
}
