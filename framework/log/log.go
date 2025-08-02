package log

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Level 日志级别
type Level int

const (
	DEBUG Level = iota
	INFO
	WARNING
	ERROR
	FATAL
)

// String 返回日志级别的字符串表示
func (l Level) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARNING:
		return "WARNING"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Logger 日志接口
type Logger interface {
	Debug(message string, context map[string]interface{})
	Info(message string, context map[string]interface{})
	Warning(message string, context map[string]interface{})
	Error(message string, context map[string]interface{})
	Fatal(message string, context map[string]interface{})
}

// FileLogger 文件日志记录器
type FileLogger struct {
	debugLogger *log.Logger
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
	fatalLogger *log.Logger
	level       Level
}

// NewFileLogger 创建文件日志记录器
func NewFileLogger(logDir string, level Level) (*FileLogger, error) {
	// 确保日志目录存在
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	// 创建日志文件
	debugFile, err := os.OpenFile(fmt.Sprintf("%s/debug.log", logDir), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	infoFile, err := os.OpenFile(fmt.Sprintf("%s/info.log", logDir), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	warnFile, err := os.OpenFile(fmt.Sprintf("%s/warning.log", logDir), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	errorFile, err := os.OpenFile(fmt.Sprintf("%s/error.log", logDir), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	fatalFile, err := os.OpenFile(fmt.Sprintf("%s/fatal.log", logDir), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &FileLogger{
		debugLogger: log.New(debugFile, "[DEBUG] ", log.LstdFlags),
		infoLogger:  log.New(infoFile, "[INFO] ", log.LstdFlags),
		warnLogger:  log.New(warnFile, "[WARNING] ", log.LstdFlags),
		errorLogger: log.New(errorFile, "[ERROR] ", log.LstdFlags),
		fatalLogger: log.New(fatalFile, "[FATAL] ", log.LstdFlags),
		level:       level,
	}, nil
}

// Debug 记录调试日志
func (l *FileLogger) Debug(message string, context map[string]interface{}) {
	if l.level <= DEBUG {
		l.debugLogger.Printf("%s %s", message, formatContext(context))
	}
}

// Info 记录信息日志
func (l *FileLogger) Info(message string, context map[string]interface{}) {
	if l.level <= INFO {
		l.infoLogger.Printf("%s %s", message, formatContext(context))
	}
}

// Warning 记录警告日志
func (l *FileLogger) Warning(message string, context map[string]interface{}) {
	if l.level <= WARNING {
		l.warnLogger.Printf("%s %s", message, formatContext(context))
	}
}

// Error 记录错误日志
func (l *FileLogger) Error(message string, context map[string]interface{}) {
	if l.level <= ERROR {
		l.errorLogger.Printf("%s %s", message, formatContext(context))
	}
}

// Fatal 记录致命错误日志
func (l *FileLogger) Fatal(message string, context map[string]interface{}) {
	if l.level <= FATAL {
		l.fatalLogger.Printf("%s %s", message, formatContext(context))
	}
}

// ConsoleLogger 控制台日志记录器
type ConsoleLogger struct {
	level Level
}

// NewConsoleLogger 创建控制台日志记录器
func NewConsoleLogger(level Level) *ConsoleLogger {
	return &ConsoleLogger{
		level: level,
	}
}

// Debug 记录调试日志
func (l *ConsoleLogger) Debug(message string, context map[string]interface{}) {
	if l.level <= DEBUG {
		fmt.Printf("[DEBUG] %s %s %s\n", time.Now().Format("2006-01-02 15:04:05"), message, formatContext(context))
	}
}

// Info 记录信息日志
func (l *ConsoleLogger) Info(message string, context map[string]interface{}) {
	if l.level <= INFO {
		fmt.Printf("[INFO] %s %s %s\n", time.Now().Format("2006-01-02 15:04:05"), message, formatContext(context))
	}
}

// Warning 记录警告日志
func (l *ConsoleLogger) Warning(message string, context map[string]interface{}) {
	if l.level <= WARNING {
		fmt.Printf("[WARNING] %s %s %s\n", time.Now().Format("2006-01-02 15:04:05"), message, formatContext(context))
	}
}

// Error 记录错误日志
func (l *ConsoleLogger) Error(message string, context map[string]interface{}) {
	if l.level <= ERROR {
		fmt.Printf("[ERROR] %s %s %s\n", time.Now().Format("2006-01-02 15:04:05"), message, formatContext(context))
	}
}

// Fatal 记录致命错误日志
func (l *ConsoleLogger) Fatal(message string, context map[string]interface{}) {
	if l.level <= FATAL {
		fmt.Printf("[FATAL] %s %s %s\n", time.Now().Format("2006-01-02 15:04:05"), message, formatContext(context))
	}
}

// formatContext 格式化上下文
func formatContext(context map[string]interface{}) string {
	if context == nil || len(context) == 0 {
		return ""
	}

	result := "{"
	first := true
	for key, value := range context {
		if !first {
			result += ", "
		}
		result += fmt.Sprintf("%s: %v", key, value)
		first = false
	}
	result += "}"
	return result
}

// 全局日志实例
var defaultLogger Logger

// 初始化默认日志记录器
func init() {
	defaultLogger = NewConsoleLogger(INFO)
}

// SetDefaultLogger 设置默认日志记录器
func SetDefaultLogger(logger Logger) {
	defaultLogger = logger
}

// Debug 全局调试日志
func Debug(message string, context map[string]interface{}) {
	if defaultLogger != nil {
		defaultLogger.Debug(message, context)
	}
}

// Info 全局信息日志
func Info(message string, context map[string]interface{}) {
	if defaultLogger != nil {
		defaultLogger.Info(message, context)
	}
}

// Warning 全局警告日志
func Warning(message string, context map[string]interface{}) {
	if defaultLogger != nil {
		defaultLogger.Warning(message, context)
	}
}

// Error 全局错误日志
func Error(message string, context map[string]interface{}) {
	if defaultLogger != nil {
		defaultLogger.Error(message, context)
	}
}

// Fatal 全局致命错误日志
func Fatal(message string, context map[string]interface{}) {
	if defaultLogger != nil {
		defaultLogger.Fatal(message, context)
	}
} 