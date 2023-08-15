package logger

import (
	"context"
	"fmt"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/helpers"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type GormLogger struct {
	ctx              context.Context `json:"ctx,omitempty"`
	logger.Interface `json:"logger.Interface,omitempty"`
}

func NewGormLogger(ctx context.Context) *GormLogger {
	return &GormLogger{
		ctx,
		logger.New(
			log.New(os.Stdout, "\n", log.LstdFlags),
			logger.Config{
				LogLevel:                  logger.Info,
				SlowThreshold:             200 * time.Millisecond,
				IgnoreRecordNotFoundError: false,
				Colorful:                  true,
				ParameterizedQueries:      false,
			},
		),
	}
}

func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	gLogger := *l
	gLogger.Interface = l.Interface.LogMode(level)
	return &gLogger
}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	msg = fmt.Sprintf("[%s] %s", helpers.GetRequestIdFromContext(ctx), msg)
	l.Interface.Info(ctx, msg, data...)
}

func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	msg = fmt.Sprintf("[%s] %s", helpers.GetRequestIdFromContext(ctx), msg)
	l.Interface.Warn(ctx, msg, data...)
}

func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	msg = fmt.Sprintf("[%s] %s", helpers.GetRequestIdFromContext(ctx), msg)
	l.Interface.Error(ctx, msg, data...)
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) { /**/
	sql, rows := fc()
	nfc := func() (string, int64) {
		msg := fmt.Sprintf("[%s] %s", helpers.GetRequestIdFromContext(ctx), sql)
		return msg, rows
	}
	l.Interface.Trace(ctx, begin, nfc, err)
}
