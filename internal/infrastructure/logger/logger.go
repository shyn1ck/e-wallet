package logger

import (
	"e-wallet/internal/infrastructure/config"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	Info    *log.Logger
	Error   *log.Logger
	Warning *log.Logger
	Debug   *log.Logger
	GORMLog *log.Logger
)

func Init(logCfg config.LogConfig) error {
	if _, err := os.Stat(logCfg.Directory); os.IsNotExist(err) {
		err = os.Mkdir(logCfg.Directory, 0755)
		if err != nil {
			return err
		}
	}

	lumberLogInfo := &lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/%s", logCfg.Directory, logCfg.InfoFile),
		MaxSize:    logCfg.MaxSizeMegabytes,
		MaxBackups: logCfg.MaxBackups,
		MaxAge:     logCfg.MaxAgeDays,
		Compress:   logCfg.Compress,
		LocalTime:  logCfg.LocalTime,
	}

	lumberLogError := &lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/%s", logCfg.Directory, logCfg.ErrorFile),
		MaxSize:    logCfg.MaxSizeMegabytes,
		MaxBackups: logCfg.MaxBackups,
		MaxAge:     logCfg.MaxAgeDays,
		Compress:   logCfg.Compress,
		LocalTime:  logCfg.LocalTime,
	}

	lumberLogWarning := &lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/%s", logCfg.Directory, logCfg.WarnFile),
		MaxSize:    logCfg.MaxSizeMegabytes,
		MaxBackups: logCfg.MaxBackups,
		MaxAge:     logCfg.MaxAgeDays,
		Compress:   logCfg.Compress,
		LocalTime:  logCfg.LocalTime,
	}

	lumberLogDebug := &lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/%s", logCfg.Directory, logCfg.DebugFile),
		MaxSize:    logCfg.MaxSizeMegabytes,
		MaxBackups: logCfg.MaxBackups,
		MaxAge:     logCfg.MaxAgeDays,
		Compress:   logCfg.Compress,
		LocalTime:  logCfg.LocalTime,
	}

	lumberLogGORM := &lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/%s", logCfg.Directory, logCfg.GormFile),
		MaxSize:    logCfg.MaxSizeMegabytes,
		MaxBackups: logCfg.MaxBackups,
		MaxAge:     logCfg.MaxAgeDays,
		Compress:   logCfg.Compress,
		LocalTime:  logCfg.LocalTime,
	}

	gin.DefaultWriter = io.MultiWriter(lumberLogInfo)

	Info = log.New(gin.DefaultWriter, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(lumberLogError, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(lumberLogWarning, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	Debug = log.New(lumberLogDebug, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	GORMLog = log.New(lumberLogGORM, "GORM: ", log.Ldate|log.Ltime|log.Lshortfile)

	return nil
}
