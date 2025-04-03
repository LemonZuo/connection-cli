package logger

import (
	"io"
	"log"
	"os"
	"sync"
)

var (
	// 确保初始化只执行一次
	once sync.Once
	// 全局日志文件
	logFile *os.File
)

// Init 初始化日志系统
func Init() {
	once.Do(func() {
		// 创建或打开日志文件
		var err error
		logFile, err = os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Printf("Warning: Failed to open log file: %v", err)
			return
		}

		// 同时输出到控制台和文件
		mw := io.MultiWriter(os.Stdout, logFile)
		log.SetOutput(mw)
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	})
}

// Close 关闭日志文件
func Close() {
	if logFile != nil {
		logFile.Close()
	}
}

// Info 记录信息日志
func Info(format string, v ...interface{}) {
	log.Printf(format, v...)
}

// Error 记录错误日志
func Error(format string, v ...interface{}) {
	log.Printf("ERROR: "+format, v...)
}

// Debug 记录调试日志
func Debug(format string, v ...interface{}) {
	log.Printf("DEBUG: "+format, v...)
} 