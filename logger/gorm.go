package logger

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/zhanglp0129/lpdrive-server/config"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

// GormLogger gorm自定义日志
type GormLogger struct{}

func (g *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return &GormLogger{}
}

func (g *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	L.WithContext(ctx).Infof(msg, data)
}

func (g *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	L.WithContext(ctx).Warnf(msg, data)
}

func (g *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	L.WithContext(ctx).Errorf(msg, data)
}

func (g *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rows int64), err error) {
	elapsed := time.Since(begin)
	slowThreshold := time.Duration(config.C.Database.SlowThreshold) * time.Millisecond
	sql, rows := fc()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		L.WithContext(ctx).WithError(err).WithFields(logrus.Fields{
			"sql":     sql,
			"rows":    rows,
			"elapsed": elapsed,
		}).Error()
	} else if elapsed > slowThreshold {
		L.WithContext(ctx).WithFields(logrus.Fields{
			"sql":     sql,
			"rows":    rows,
			"elapsed": elapsed,
		}).Warn()
	} else {
		L.WithContext(ctx).WithFields(logrus.Fields{
			"sql":     sql,
			"rows":    rows,
			"elapsed": elapsed,
		}).Info()
	}
}
