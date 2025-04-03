package logger

import (
	"io"
	"log"
	"os"
	"sync"
	"time"
)

const (
	// LogFileName 日志文件名
	LogFileName = "app.log"
	// CleanupInterval 清理周期（7天）
	CleanupInterval = 7 * 24 * time.Hour
)

var (
	// 确保初始化只执行一次
	once sync.Once
	// 全局日志文件
	logFile *os.File
	// 清理定时器
	cleanupTimer *time.Timer
	// 停止清理的通道
	stopCleanup chan struct{}
)

// Init 初始化日志系统
func Init() {
	once.Do(func() {
		// 创建或打开日志文件
		var err error
		logFile, err = os.OpenFile(LogFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Printf("Warning: Failed to open log file: %v", err)
			return
		}

		// 同时输出到控制台和文件
		mw := io.MultiWriter(os.Stdout, logFile)
		log.SetOutput(mw)
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

		// 启动日志清理定时器
		startLogCleanup()
	})
}

// Close 关闭日志文件和停止清理定时器
func Close() {
	// 停止清理定时器
	if stopCleanup != nil {
		close(stopCleanup)
	}
	
	// 关闭日志文件
	if logFile != nil {
		logFile.Close()
	}
}

// 启动日志清理定时器
func startLogCleanup() {
	stopCleanup = make(chan struct{})
	
	// 首次启动立即检查
	go cleanupLog()
	
	// 启动定时器
	go func() {
		ticker := time.NewTicker(CleanupInterval)
		defer ticker.Stop()
		
		for {
			select {
			case <-ticker.C:
				cleanupLog()
			case <-stopCleanup:
				return
			}
		}
	}()
}

// 清理日志文件
func cleanupLog() {
	// 记录日志清理开始
	Info("Starting scheduled log cleanup")
	
	// 检查文件大小
	fileInfo, err := os.Stat(LogFileName)
	if err != nil {
		Error("Failed to get log file info: %v", err)
		return
	}
	
	// 记录当前文件大小
	sizeMB := float64(fileInfo.Size()) / (1024 * 1024)
	Info("Current log file size: %.2f MB", sizeMB)
	
	// 如果文件存在并且大于1MB，则清空它
	if fileInfo.Size() > 1024*1024 {
		// 先关闭当前文件
		if logFile != nil {
			logFile.Close()
			logFile = nil
		}
		
		// 清空文件内容
		err = os.Truncate(LogFileName, 0)
		if err != nil {
			log.Printf("Error truncating log file: %v", err)
		}
		
		// 重新打开文件
		logFile, err = os.OpenFile(LogFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Printf("Warning: Failed to reopen log file: %v", err)
			return
		}
		
		// 重新设置输出
		mw := io.MultiWriter(os.Stdout, logFile)
		log.SetOutput(mw)
		
		Info("Log file has been cleaned up")
	} else {
		Info("Log file is still small, skipping cleanup")
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