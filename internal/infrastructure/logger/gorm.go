package logger

import (
	"context"
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewGORMLogger(logWriter *log.Logger) logger.Interface {
	return &GORMLogger{
		LogWriter: logWriter,
	}
}

type GORMLogger struct {
	LogWriter *log.Logger
}

func (l *GORMLogger) LogMode(level logger.LogLevel) logger.Interface {
	return &GORMLogger{LogWriter: l.LogWriter}
}

func (l *GORMLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	l.LogWriter.Printf("INFO: "+msg, args...)
}

func (l *GORMLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	l.LogWriter.Printf("WARN: "+msg, args...)
}

func (l *GORMLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	l.LogWriter.Printf("ERROR: "+msg, args...)
}

func (l *GORMLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	sql, rows := fc()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return
		}
		l.LogWriter.Printf("ERROR: %v | SQL: %s | Rows: %d", err, sql, rows)
		return
	}
	l.LogWriter.Printf("SQL: %s | Rows: %d | Duration: %v", sql, rows, time.Since(begin))
}
