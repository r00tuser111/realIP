package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

// 日志级别
const (
	LevelDebug = iota
	LevelInfo
	LevelWarn
	LevelError
)

// Logger 日志记录器
type Logger struct {
	enabled bool         // 是否启用日志
	level   int          // 日志级别
	logger  *log.Logger  // Go原生日志器
}

// 日志级别的颜色代码
var levelColors = map[int]string{
	LevelDebug: "\033[36m", // 青色
	LevelInfo:  "\033[32m", // 绿色
	LevelWarn:  "\033[33m", // 黄色
	LevelError: "\033[31m", // 红色
}

// 日志级别的前缀
var levelPrefixes = map[int]string{
	LevelDebug: "[DEBUG]",
	LevelInfo:  "[INFO] ",
	LevelWarn:  "[WARN] ",
	LevelError: "[ERROR]",
}

// New 创建新的日志记录器
func New(enabled bool) *Logger {
	return &Logger{
		enabled: enabled,
		level:   LevelDebug,
		logger:  log.New(os.Stderr, "", 0),
	}
}

// SetLevel 设置日志级别
func (l *Logger) SetLevel(level int) {
	l.level = level
}

// log 通用日志方法
func (l *Logger) log(level int, format string, args ...interface{}) {
	if !l.enabled || level < l.level {
		return
	}

	// 构建日志消息
	timestamp := time.Now().Format("15:04:05")
	prefix := levelPrefixes[level]
	color := levelColors[level]
	reset := "\033[0m"
	
	msg := fmt.Sprintf(format, args...)
	logMsg := fmt.Sprintf("%s %s%s%s %s", timestamp, color, prefix, reset, msg)
	
	l.logger.Println(logMsg)
}

// Debug 调试级别日志
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(LevelDebug, format, args...)
}

// Info 信息级别日志
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(LevelInfo, format, args...)
}

// Warn 警告级别日志
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(LevelWarn, format, args...)
}

// Error 错误级别日志
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(LevelError, format, args...)
} 