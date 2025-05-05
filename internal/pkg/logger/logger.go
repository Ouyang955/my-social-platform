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
// 类似于Java中的初始化方法，返回error类型（Go的错误处理机制，类似Java的Exception但更轻量）
func InitLogger() error {
	// 创建logs目录，0755是权限设置（读写执行权限，类似Linux chmod）
	// os.MkdirAll相当于Java中的File.mkdirs()，可以创建多级目录
	if err := os.MkdirAll("logs", 0755); err != nil {
		// Go中错误处理使用返回值而非异常，类似于Java中的return new IOException()
		return err
	}

	// 打开日志文件，使用当前日期作为文件名
	// time.Now().Format格式化日期，2006-01-02是Go特有的格式化模板（相当于Java的yyyy-MM-dd）
	filename := "logs/" + time.Now().Format("2006-01-02") + ".log"
	// os.OpenFile类似Java的new FileOutputStream()
	// os.O_APPEND|os.O_CREATE|os.O_WRONLY是位运算组合的文件打开模式：
	// - O_APPEND: 追加写入（类似Java中的FileWriter(file, true)）
	// - O_CREATE: 如果不存在则创建
	// - O_WRONLY: 只写模式
	// 0644是文件权限设置
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	// 将打开的文件赋值给全局变量logFile（类似Java中的类成员变量）
	logFile = file
	// 设置标准日志输出到文件，类似Java中的System.setOut()   原本log.println会输出到终端
	log.SetOutput(file)
	// 成功时返回nil（Go中表示空，类似Java中的null）
	return nil
}

// Log 记录日志
// 函数接收多个参数：日志级别、操作类型、用户、IP地址和消息内容
// Go中函数可以有多个返回值，这里没有返回值（类似Java中的void方法）
func Log(level LogLevel, action, user, ip, message string) {
	// 创建LogEntry结构体实例（类似Java中创建一个对象）
	// Go中使用字段名直接初始化结构体（类似Java的Builder模式）
	entry := LogEntry{
		Timestamp: time.Now(), // 当前时间，类似Java的new Date()
		Level:     level,      // 日志级别
		Action:    action,     // 操作类型
		User:      user,       // 用户标识
		IP:        ip,         // IP地址
		Message:   message,    // 日志消息
	}

	// 格式化日志输出
	// log.Printf类似Java中的String.format()和logger.info()的组合
	// 使用格式化字符串和变量创建日志内容
	// 输出到日志文件中，因为在InitLogger()中已经通过log.SetOutput(file)
	// 将标准日志输出重定向到了文件，所以这里的log.Printf会写入到文件而非终端
	log.Printf("[%s] %s | Action: %s | User: %s | IP: %s | Message: %s",
		entry.Level, // 日志级别
		entry.Timestamp.Format("2006-01-02 15:04:05"), // 格式化时间戳
		entry.Action,  // 操作类型
		entry.User,    // 用户标识
		entry.IP,      // IP地址
		entry.Message, // 日志消息
	)
}

// Close 关闭日志文件
// 类似Java中实现Closeable接口的close()方法
func Close() {
	// 检查logFile是否为nil（类似Java中的null检查）
	if logFile != nil {
		// 关闭文件，释放资源（类似Java中的file.close()）
		logFile.Close()
	}
}
