/*
 * @Author: cloudyi.li
 * @Date: 2023-02-15 11:44:21
 * @LastEditTime: 2023-04-05 15:54:43
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/logger/log.go
 */
package logger

import (
	"chatserver-api/internal/consts"
	"chatserver-api/pkg/config"
	"context"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	_logger *logger
	once    sync.Once
)

type logger struct {
	cfg    *config.LogConfig
	sugar  *zap.SugaredLogger
	_level zapcore.Level
}

// DefaultPair 表示接收打印的键值对参数
type DefaultPair struct {
	key   string
	value interface{}
}

func Pair(key string, v interface{}) DefaultPair {
	return DefaultPair{
		key:   key,
		value: v,
	}
}

func spread(kvs ...DefaultPair) []interface{} {
	s := make([]interface{}, 0, len(kvs))
	for _, v := range kvs {
		s = append(s, v.key, v.value)
	}
	return s
}

// Debug 打印debug级别信息
func Debug(message string, kvs ...DefaultPair) {
	if !_logger.EnabledLevel(zapcore.DebugLevel) {
		return
	}
	args := spread(kvs...)
	_logger.sugar.Debugw(message, args...)
}

// Info 打印info级别信息
func Info(message string, kvs ...DefaultPair) {
	if !_logger.EnabledLevel(zapcore.InfoLevel) {
		return
	}
	args := spread(kvs...)
	_logger.sugar.Infow(message, args...)
}

// Warn 打印warn级别信息
func Warn(message string, kvs ...DefaultPair) {
	if !_logger.EnabledLevel(zapcore.WarnLevel) {
		return
	}
	args := spread(kvs...)
	_logger.sugar.Warnw(message, args...)
}

// Error 打印error级别信息
func Error(message string, kvs ...DefaultPair) {
	if !_logger.EnabledLevel(zapcore.ErrorLevel) {
		return
	}
	args := spread(kvs...)
	_logger.sugar.Errorw(message, args...)
}

// Panic 打印错误信息，然后panic
func Panic(message string, kvs ...DefaultPair) {
	if !_logger.EnabledLevel(zapcore.PanicLevel) {
		return
	}
	args := spread(kvs...)
	_logger.sugar.Panicw(message, args...)
}

// Fatal 打印错误信息，然后退出
func Fatal(message string, kvs ...DefaultPair) {
	if !_logger.EnabledLevel(zapcore.FatalLevel) {
		return
	}
	args := spread(kvs...)
	_logger.sugar.Fatalw(message, args...)
}

// Debugf 格式化输出debug级别日志
func Debugf(template string, args ...interface{}) {
	_logger.sugar.Debugf(template, args...)
}

// Infof 格式化输出info级别日志
func Infof(template string, args ...interface{}) {
	_logger.sugar.Infof(template, args...)
}

// Warnf 格式化输出warn级别日志
func Warnf(template string, args ...interface{}) {
	_logger.sugar.Warnf(template, args...)
}

// Errorf 格式化输出error级别日志
func Errorf(template string, args ...interface{}) {
	_logger.sugar.Errorf(template, args...)
}

// Panicf 格式化输出日志，并panic
func Panicf(template string, args ...interface{}) {
	_logger.sugar.Panicf(template, args...)
}

// Fatalf 格式化输出日志，并退出
func Fatalf(template string, args ...interface{}) {
	_logger.sugar.Fatalf(template, args...)
}

// tempLogger 临时的logger
type tempLogger struct {
	extra []DefaultPair
}

// getPrefix 根据extra生成日志前缀，比如 "requestId:%s name:%s "
func (tl *tempLogger) getPrefix(template string, args []interface{}) ([]interface{}, string) {

	if len(tl.extra) > 0 {
		var prefix string
		tmp := make([]interface{}, 0, len(args)+len(tl.extra))
		for _, pair := range tl.extra {
			prefix += pair.key + ":%s,"
			tmp = append(tmp, pair.value)
		}
		args = append(tmp, args...)
		template = prefix + template
	}
	return args, template
}

func (tl *tempLogger) getArgs(kvs []DefaultPair) []interface{} {
	var args []interface{}
	if len(tl.extra) > 0 {
		tl.extra = append(tl.extra, kvs...)
		args = spread(tl.extra...)
	} else {
		args = spread(kvs...)
	}
	return args
}

// RID 实现rid(RequestID打印) 使用格式 log.RID(ctx).Debug(), 可以继续拓展 比如Log.RID(ctx).AppName(ctx).Debug()
func RID(ctx context.Context) *tempLogger {
	tl := &tempLogger{extra: make([]DefaultPair, 0)}
	if ctx == nil {
		return tl
	}
	if v := ctx.Value(consts.RequestId); v != nil && v != "" {
		tl.extra = append(tl.extra, Pair(consts.RequestId, v))
	}
	return tl
}

func (tl *tempLogger) Debug(message string, kvs ...DefaultPair) {
	// 这里重复写的原因是zap的log设置的SKIP是1，
	//并且使用的全局只有一个logger，不能修改SKIP，否则打印的位置不正确，后续都是重复代码
	// Debug(message, tl.extra...) 这种写法要修改SKIP
	if !_logger.EnabledLevel(zapcore.DebugLevel) {
		return
	}
	args := tl.getArgs(kvs)
	_logger.sugar.Debugw(message, args...)
}

func (tl *tempLogger) Info(message string, kvs ...DefaultPair) {
	if !_logger.EnabledLevel(zapcore.InfoLevel) {
		return
	}
	args := tl.getArgs(kvs)
	_logger.sugar.Infow(message, args...)
}

// Warn 打印warn级别信息
func (tl *tempLogger) Warn(message string, kvs ...DefaultPair) {
	if !_logger.EnabledLevel(zapcore.WarnLevel) {
		return
	}
	args := tl.getArgs(kvs)
	_logger.sugar.Warnw(message, args...)
}

// Error 打印error级别信息
func (tl *tempLogger) Error(message string, kvs ...DefaultPair) {
	if !_logger.EnabledLevel(zapcore.ErrorLevel) {
		return
	}
	args := tl.getArgs(kvs)
	_logger.sugar.Errorw(message, args...)
}

// Panic 打印错误信息，然后panic
func (tl *tempLogger) Panic(message string, kvs ...DefaultPair) {
	if !_logger.EnabledLevel(zapcore.PanicLevel) {
		return
	}
	args := tl.getArgs(kvs)
	_logger.sugar.Panicw(message, args...)
}

// Fatal 打印错误信息，然后退出
func (tl *tempLogger) Fatal(message string, kvs ...DefaultPair) {
	if !_logger.EnabledLevel(zapcore.FatalLevel) {
		return
	}
	args := tl.getArgs(kvs)
	_logger.sugar.Fatalw(message, args...)
}

// Debugf 格式化输出debug级别日志
func (tl *tempLogger) Debugf(template string, args ...interface{}) {
	args, template = tl.getPrefix(template, args)
	_logger.sugar.Debugf(template, args...)
}

// Infof 格式化输出info级别日志
func (tl *tempLogger) Infof(template string, args ...interface{}) {
	args, template = tl.getPrefix(template, args)
	_logger.sugar.Infof(template, args...)
}

// Warnf 格式化输出warn级别日志
func (tl *tempLogger) Warnf(template string, args ...interface{}) {
	args, template = tl.getPrefix(template, args)
	_logger.sugar.Warnf(template, args...)
}

// Errorf 格式化输出error级别日志
func (tl *tempLogger) Errorf(template string, args ...interface{}) {
	args, template = tl.getPrefix(template, args)
	_logger.sugar.Errorf(template, args...)
}

// Panicf 格式化输出日志，并panic
func (tl *tempLogger) Panicf(template string, args ...interface{}) {
	args, template = tl.getPrefix(template, args)
	_logger.sugar.Panicf(template, args...)
}

// Fatalf 格式化输出日志，并退出
func (tl *tempLogger) Fatalf(template string, args ...interface{}) {
	args, template = tl.getPrefix(template, args)
	_logger.sugar.Fatalf(template, args...)
}
