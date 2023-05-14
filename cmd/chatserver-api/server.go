/*
 * @Author: cloudyi.li
 * @Date: 2023-03-29 10:05:07
 * @LastEditTime: 2023-05-14 18:38:47
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/cmd/chatserver-api/server.go
 */
package chatserverapi

import (
	"chatserver-api/pkg/config"
	"chatserver-api/pkg/logger"
	"chatserver-api/pkg/validator"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

// HttpServer 代表当前服务端实例
type HttpServer struct {
	config *config.Config
	f      func()
}

// NewHttpServer 创建server实例
func NewHttpServer(config *config.Config) *HttpServer {
	return &HttpServer{
		config: config,
	}
}

// Router 加载路由，使用侧提供接口，实现侧需要实现该接口
type Router interface {
	Load(engine *gin.Engine)
}

// Run server的启动入口
// 加载路由, 启动服务
func (s *HttpServer) Run(rs ...Router) {
	var wg sync.WaitGroup
	wg.Add(1)
	// 设置gin启动模式，必须在创建gin实例之前
	gin.SetMode(s.config.Mode)
	g := gin.New()
	s.routerLoad(g, rs...)
	// gin validator替换
	validator.LazyInitGinValidator(s.config.Language)

	// health check
	go func() {
		if err := Ping(s.config.Port, s.config.MaxPingCount); err != nil {
			logger.Fatal("server no response")
		}
		logger.Infof("server started success! port: %s", s.config.Port)
	}()

	srv := http.Server{
		Addr:    s.config.Port,
		Handler: g,
	}
	if s.f != nil {
		srv.RegisterOnShutdown(s.f)
	}
	// graceful shutdown
	sgn := make(chan os.Signal, 1)
	signal.Notify(sgn, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sgn
		logger.Infof("server shutdown")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			logger.Errorf("server shutdown err %v \n", err)
		}
		wg.Done()
	}()

	err := srv.ListenAndServe()
	if err != nil {
		if err != http.ErrServerClosed {
			logger.Errorf("server start failed on port %s", s.config.Port)
			return
		}
	}
	wg.Wait()
	logger.Infof("server stop on port %s", s.config.Port)
}

// RouterLoad 加载自定义路由
func (s *HttpServer) routerLoad(g *gin.Engine, rs ...Router) *HttpServer {
	for _, r := range rs {
		r.Load(g)
	}
	return s
}

// RegisterOnShutdown 注册shutdown后的回调处理函数，用于清理资源
func (s *HttpServer) RegisterOnShutdown(_f func()) {
	s.f = _f
}

// Ping 用来检查是否程序正常启动
func Ping(port string, maxCount int) error {
	seconds := 1
	if len(port) == 0 {
		panic("Please specify the service port")
	}
	if !strings.HasPrefix(port, ":") {
		port += ":"
	}
	url := fmt.Sprintf("http://localhost%s/ping", port)
	for i := 0; i < maxCount; i++ {
		resp, err := http.Get(url)
		if nil == err && resp != nil && resp.StatusCode == http.StatusOK {
			return nil
		}
		logger.Infof("等待服务在线, 已等待 %d 秒，最多等待 %d 秒", seconds, maxCount)
		time.Sleep(time.Second * 1)
		seconds++
	}
	return fmt.Errorf("服务启动失败，端口 %s", port)
}
