/*
 * @Author: cloudyi.li
 * @Date: 2023-06-16 21:16:47
 * @LastEditTime: 2023-06-27 14:03:34
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/chatfunc/chatfunc.go
 */
package chatfunc

import (
	"chatserver-api/internal/consts"
	"chatserver-api/pkg/logger"
	"chatserver-api/pkg/openai"
	"chatserver-api/pkg/search"
	"context"
	"encoding/json"
	"time"
)

func ChatFuncProcess(ctx context.Context, funcName string, arguments string) (content string) {
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(arguments), &result); err != nil {
		logger.Errorf("%v", err)
	}
	switch funcName {
	case "GetWeather":
		content = GetWeather(ctx, result["location"].(string))
	case "EntitySearch":
		content = EntitySearch(ctx, result["query"].(string), result["etype"].(string))
	case "GoogleSearch":
		content = GoogleSearch(ctx, result["query"].(string), result["classify"].(string))
	default:
		logger.Errorf("%v", "未匹配")
	}
	return content
}

func queryExtra(message string) string {
	if message == "" {
		return ""
	}
	var req openai.CompletionRequest
	req.Stop = []string{"\n"}
	req.Temperature = 0
	req.Prompt = `Your current task is to determine whether the user input requires calling some search function to query the latest content. If there is no need to call a function, it must only return '''[NONE_FUNCTION]'''. Partial return content of the reference example,using JSON format.` + "\nCurrent Date:" + time.Now().Format(consts.TimeLayout) + "\n" + `## Functions:'''
"functions" : [   {  "name": "EntitySearch",  "description": "When the content sent by the user contains special entities, call this function to perform a knowledge graph search.",  "parameters": {  "type": "object",  "properties": {  "query": {  "type": "string",  "description": "the entiy name"  },  "type": {  "type": "string",  "enum": [  "Action",  "BioChemEntity",  "CreativeWork",  "Event",  "Intangible",  "MedicalEntity",  "Organization",  "Person",  "Place",  "Product",  "Thing"  ],  "description": "the entiy type"  }  },  "required": [  "query",  "type"  ]  }  },  {  "name": "GoogleSearch",  "description": "When the content sent by the user contains some newly occurred events, news or other things that you don't know, use this function to call the Google search engine.",  "parameters": {  "type": "object",  "properties": {  "query": {  "type": "string",  "description": "Extracted query statements from user questions that are applicable for search engines."  },  "classify": {  "type": "string",  "enum": [  "News",  "Custom"  ],  "description": "Determine what category the content to be searched belongs to."  }  },  "required": [  "query",  "classify"  ]  }  },  {  "name": "GetWeather",  "description": "Get the current weather",  "parameters": {  "type": "object",  "properties": {  "location": {  "type": "string",  "description": "The city and state, e.g. San Francisco, CA"  }  },  "required": [  "location"  ]  }  }  ]
'''
 Q: 南京第一医院介绍
 A: {  "name": "EntitySearch",  "argument": "{   \"query\": \"南京第一医院\",   \"type\": \"Organization\"  }" }

 Q: 南京明天什么天气
 A: {  "name": "GetWeather",  "argument": "{   \"location\": \"南京\"  }" }

 Q: 这个事情怎么评价呢
 A: {  "name": "[NONE_FUNCTION]",  "argument": \"\" }

 Q: RTX4070显卡性能
 A: {  "name": "GoogleSearch",  "argument": "{   \"query\": \"RTX 4070 显卡 性能\", \"classify\": \"Custom\" }" }

 Q: 今年有什么电影推荐
 A: {  "name": "GoogleSearch",  "argument": "{   \"query\": \"2023电影推荐\", \"classify\": \"News\" }" }` + "\n\n" + "Q: " + message + "\nA: "
	req.Model = "text-davinci-003"
	req.MaxTokens = 100
	client, err := openai.NewClient()
	if err != nil {
		logger.Errorf("%s", err)
		return ""
	}
	resp, err := client.CreateCompletion(req)
	if err != nil {
		logger.Errorf("%s", err)
		return ""
	}
	content := resp.Choices[0].Text
	return content
}

func CustomFuncExtension(ctx context.Context, queryoriginal string) (content string, err error) {
	var result map[string]interface{}
	ner, _ := search.NerDetec(queryoriginal)
	if ner == 0 {
		logger.Debug("[NerDetec]", logger.Pair("是否返回", false))
		return "", nil
	}
	funcCall := queryExtra(queryoriginal)
	if err := json.Unmarshal([]byte(funcCall), &result); err != nil {
		logger.Errorf("%v", err)
	}
	funcname := result["name"].(string)
	funcargument := result["argument"].(string)
	logger.Debug("[FuncExtension]", logger.Pair("调用函数", funcname), logger.Pair("调用参数", funcargument))
	if funcargument != "" {
		content = ChatFuncProcess(ctx, funcname, funcargument)
	}
	return

}
