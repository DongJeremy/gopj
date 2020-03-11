package common

import (
	"fmt"
	"io"
	"math"
	"os"
	"path"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	// DefaultFormat default time format
	DefaultFormat = "2006-01-02 15:04:05"
)

// LoggerFromConfig get logger by filename
func LoggerFromConfig(config *Config) *logrus.Logger {
	logFilePath := config.LogFilePath
	logFileName := config.LogFileName
	// 日志文件
	fileName := path.Join(logFilePath, logFileName)
	logger := logrus.New()
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		// 文件不存在,创建
		os.Create(fileName)
	}
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0666)

	if err == nil {
		//logrus.SetOutput()
		var out io.Writer
		if config.LogOutput == "both" {
			out = io.MultiWriter(file, os.Stdout)
		} else {
			out = io.MultiWriter(file)
		}
		logger.Out = out
	} else {
		logger.Info("打开 " + fileName + " 下的日志文件失败, 使用默认方式显示日志！")
	}
	return logger
}

// GinLogger 基于logrus 的log中间件
func GinLogger(log *logrus.Logger) gin.HandlerFunc {
	return func(context *gin.Context) {
		path := context.Request.URL.Path
		start := time.Now()
		context.Next()
		stop := time.Since(start)
		// 等待时间
		latency := int(math.Ceil(float64(stop.Nanoseconds()) / 1000000.0))
		statusCode := context.Writer.Status()
		clientIP := context.ClientIP()
		clientUserAgent := context.Request.UserAgent()
		referer := context.Request.Referer()
		hostname, err := os.Hostname()
		if err != nil {
			hostname = "unknow"
		}
		dataLength := context.Writer.Size()
		if dataLength < 0 {
			dataLength = 0
		}

		entry := logrus.NewEntry(log).WithFields(logrus.Fields{
			"hostname":   hostname,
			"statusCode": statusCode,
			"latency":    latency, // time to process
			"clientIP":   clientIP,
			"method":     context.Request.Method,
			"path":       path,
			"referer":    referer,
			"dataLength": dataLength,
			"userAgent":  clientUserAgent,
		})

		if len(context.Errors) > 0 {
			entry.Error(context.Errors.ByType(gin.ErrorTypePrivate).String())
		} else {
			msg := fmt.Sprintf("[%s] \"%s %s\" %d", time.Now().Format(DefaultFormat), context.Request.Method, path, statusCode)
			if statusCode > 499 {
				entry.Error(msg)
			} else if statusCode > 399 {
				entry.Warn(msg)
			} else {
				entry.Info(msg)
			}
		}
	}
}
