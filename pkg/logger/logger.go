/*
 * @Author: cloudyi.li
 * @Date: 2023-02-15 16:24:59
 * @LastEditTime: 2023-04-05 15:55:41
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/logger/logger.go
 */
package logger

import (
	"chatserver-api/internal/consts"
	"chatserver-api/pkg/config"
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InitLogger 初始化日志配置
func InitLogger(_cfg *config.LogConfig, appName string) {
	once.Do(func() {
		_logger = &logger{
			cfg: _cfg,
		}
		lumber := _logger.newLumber()
		writeSyncer := zapcore.NewMultiWriteSyncer(zapcore.AddSync(lumber))
		sugar := zap.New(_logger.newCore(writeSyncer),
			zap.ErrorOutput(writeSyncer),
			zap.AddCaller(),
			zap.AddCallerSkip(1),
			zap.Fields(zap.String("appName", appName))).
			Sugar()

		_logger.sugar = sugar
	})
}

func (l *logger) newCore(ws zapcore.WriteSyncer) zapcore.Core {
	// 默认日志级别
	atomicLevel := zap.NewAtomicLevel()
	defaultLevel := zapcore.DebugLevel
	// 会解码传递的日志级别，生成新的日志级别
	_ = (&defaultLevel).UnmarshalText([]byte(l.cfg.Level))
	atomicLevel.SetLevel(defaultLevel)
	l._level = defaultLevel

	// encoder 这部分没有放到配置文件，因为一般配置一次就不会改动
	encoder := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     l.customTimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
	var writeSyncer zapcore.WriteSyncer
	if l.cfg.Console {
		writeSyncer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout))
	} else {
		// 输出到文件时，不使用彩色日志，否则会出现乱码
		encoder.EncodeLevel = zapcore.CapitalLevelEncoder
		writeSyncer = ws
	}
	// Tips: 如果使用zapcore.NewJSONEncoder
	// encoderConfig里面就不要配置 EncodeLevel 为zapcore.CapitalColorLevelEncoder或者是
	// zapcore.LowercaseColorLevelEncoder, 不但日志级别字段不会出现颜色，而且日志级别level字段
	// 会出现乱码，因为控制颜色的字符也被JSON编码了。
	return zapcore.NewCore(zapcore.NewConsoleEncoder(encoder),
		writeSyncer,
		atomicLevel)
}

// CustomTimeEncoder 实现了 zapcore.TimeEncoder
// 实现对日期格式的自定义转换
func (l *logger) customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	format := l.cfg.TimeFormat
	if len(format) <= 0 {
		format = consts.TimeLayoutMs
	}
	enc.AppendString(t.Format(format))
}

func (l *logger) newLumber() *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   l.cfg.FileName,
		MaxSize:    l.cfg.MaxSize,
		MaxAge:     l.cfg.MaxAge,
		MaxBackups: l.cfg.MaxBackups,
		LocalTime:  l.cfg.LocalTime,
		Compress:   l.cfg.Compress,
	}
}

func (l *logger) EnabledLevel(level zapcore.Level) bool {
	return level >= l._level
}

// Sync 关闭时需要同步日志到输出
func Sync() {
	_ = _logger.sugar.Sync()
}
