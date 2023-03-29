/*
 * @Author: cloudyi.li
 * @Date: 2023-02-15 16:24:59
 * @LastEditTime: 2023-03-28 14:31:33
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/di/logger/logger.go
 */
package logger

import (
	"chatserver-api/di/config"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

func InitLogger(conf config.Log) {
	newLoggerCount := len(conf.File)
	if newLoggerCount != 0 {
		for i := 0; i < newLoggerCount; i++ {
			AddLogger(logrus.New())
		}
	}
	nextLoggerIndex := 1

	SetLevel(LevelToLogrusLevel(conf.Level))

	for _, file := range conf.File {
		currentLogger := DefaultCombinedLogger.GetLogger(nextLoggerIndex)

		switch file.Format {
		case config.LogFormatJSON:
			currentLogger.SetFormatter(&logrus.JSONFormatter{})
		case config.LogFormatText:
			currentLogger.SetFormatter(&logrus.TextFormatter{})
		}
		currentLogger.SetOutput(&lumberjack.Logger{
			Filename:   file.Path,
			MaxSize:    file.MaxSize,
			MaxAge:     file.MaxAge,
			MaxBackups: 3,
			Compress:   true,
		})
		nextLoggerIndex++
	}

}

type ProtocolLogger struct{}

const fromProtocol = "Protocol -> "

func (p ProtocolLogger) Info(format string, arg ...any) {
	Infof(fromProtocol+format, arg...)
}

func (p ProtocolLogger) Warning(format string, arg ...any) {
	Warnf(fromProtocol+format, arg...)
}

func (p ProtocolLogger) Debug(format string, arg ...any) {
	Debugf(fromProtocol+format, arg...)
}

func (p ProtocolLogger) Error(format string, arg ...any) {
	Errorf(fromProtocol+format, arg...)
}

func (p ProtocolLogger) Dump(data []byte, format string, arg ...any) {
	// if !tools.FileExist("DumpsPath") {
	// 	_ = os.MkdirAll("DumpsPath", 0o755)
	// }
	dumpFile := path.Join("DumpsPath", fmt.Sprintf("%v.dump", time.Now().Unix()))
	message := fmt.Sprintf(format, arg...)
	Errorf("出现错误 %v. 详细信息已转储至文件 %v 请连同日志提交给开发者处理", message, dumpFile)
	_ = os.WriteFile(dumpFile, data, 0o644)
}
