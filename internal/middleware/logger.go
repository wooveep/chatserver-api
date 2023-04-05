/*
 * @Author: cloudyi.li
 * @Date: 2023-03-29 11:55:27
 * @LastEditTime: 2023-04-05 15:56:13
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/middleware/logger.go
 */
package middleware

import (
	"bytes"
	"chatserver-api/internal/consts"
	"chatserver-api/pkg/logger"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger 记录每次请求的请求信息和响应信息
func Logger(c *gin.Context) {
	// 请求前
	t := time.Now()
	reqPath := c.Request.URL.Path
	reqId := c.GetString(consts.RequestId)
	method := c.Request.Method
	ip := c.ClientIP()
	requestBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		requestBody = []byte{}
	}
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))

	logger.Info("New request start",
		logger.Pair(consts.RequestId, reqId),
		logger.Pair("host", ip),
		logger.Pair("host", ip),
		logger.Pair("path", reqPath),
		logger.Pair("method", method),
		logger.Pair("body", string(requestBody)))

	c.Next()
	// 请求后
	latency := time.Since(t)
	logger.Info("New request end",
		logger.Pair(consts.RequestId, reqId),
		logger.Pair("host", ip),
		logger.Pair("path", reqPath),
		logger.Pair("cost", latency))
}
