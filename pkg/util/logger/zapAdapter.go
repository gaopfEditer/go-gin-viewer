package logger

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"strings"
	"sync"
	"time"
)

type ZapGormLogger struct {
	LogLevel              gormLogger.LogLevel
	SlowThreshold         time.Duration
	SkipErrRecordNotFound bool
}

func (l *ZapGormLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

func (l *ZapGormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormLogger.Info {
		selectLogger().Sugar().Infof(msg, data...)
	}
}

func (l *ZapGormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormLogger.Warn {
		selectLogger().Sugar().Warnf(msg, data...)
	}
}

func (l *ZapGormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormLogger.Error {
		selectLogger().Sugar().Errorf(msg, data...)
	}
}

func (l *ZapGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= gormLogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= gormLogger.Error && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.SkipErrRecordNotFound):
		sql, rows := fc()
		selectLogger().Error("trace",
			zap.Error(err),
			zap.Duration("elapsed", elapsed),
			zap.Int64("rows", rows),
			zap.String("sql", cleanSQL(sql)))
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= gormLogger.Warn:
		sql, rows := fc()
		selectLogger().Warn("slow query",
			zap.Duration("elapsed", elapsed),
			zap.Int64("rows", rows),
			zap.String("sql", cleanSQL(sql)))
	case l.LogLevel >= gormLogger.Info:
		sql, rows := fc()
		selectLogger().Info("trace",
			zap.Duration("elapsed", elapsed),
			zap.Int64("rows", rows),
			zap.String("sql", cleanSQL(sql)))
	}
}

var cleanerPool = sync.Pool{
	New: func() interface{} {
		return new(strings.Builder)
	},
}

func cleanSQL(sql string) string {
	builder := cleanerPool.Get().(*strings.Builder)
	defer func() {
		builder.Reset()
		cleanerPool.Put(builder)
	}()

	builder.Grow(len(sql)) // 预分配空间以提高性能

	for _, ch := range sql {
		if ch != '\t' && ch != '\n' {
			builder.WriteRune(ch)
		}
	}

	return builder.String()
}

//
//func (z *ZapGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
//	elapsed := time.Since(begin)
//	sql, rows := fc()
//	if err != nil {
//		z.ZapLogger.Sugar().Errorf("%s [%.2fms] [rows:%v] %s", err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
//	} else {
//		z.ZapLogger.Sugar().Infof("[%.2fms] [rows:%v] %s", float64(elapsed.Nanoseconds())/1e6, rows, sql)
//	}
//}
