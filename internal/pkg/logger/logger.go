package logger

import (
	"log"
	"os"
	"time"
)

// LogLevel 定义日志级别
type LogLevel string

const (
	INFO    LogLevel = "INFO"
	WARNING LogLevel = "WARNING"
	ERROR   LogLevel = "ERROR"
)

// LogEntry 定义日志条目结构
type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     LogLevel  `json:"level"`
	Action    string    `json:"action"`
	User      string    `json:"user,omitempty"`
	IP        string    `json:"ip,omitempty"`
	Message   string    `json:"message"`
}

var (
	// 日志文件
	logFile *os.File
)

// InitLogger 初始化日志系统
func InitLogger() error {
	// 创建logs目录
	if err := os.MkdirAll("logs", 0755); err != nil {
		return err
	}

	// 打开日志文件，使用当前日期作为文件名
	filename := "logs/" + time.Now().Format("2006-01-02") + ".log"
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	logFile = file
	log.SetOutput(file)
	return nil
}

// Log 记录日志
func Log(level LogLevel, action, user, ip, message string) {
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Action:    action,
		User:      user,
		IP:        ip,
		Message:   message,
	}

	// 格式化日志输出
	log.Printf("[%s] %s | Action: %s | User: %s | IP: %s | Message: %s",
		entry.Level,
		entry.Timestamp.Format("2006-01-02 15:04:05"),
		entry.Action,
		entry.User,
		entry.IP,
		entry.Message,
	)
}

// Close 关闭日志文件
func Close() {
	if logFile != nil {
		logFile.Close()
	}
}
