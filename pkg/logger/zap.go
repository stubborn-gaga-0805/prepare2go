package logger

import (
	"context"
	"errors"
	"fmt"
	"github.com/natefinch/lumberjack"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/consts"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/wework"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

var (
	requestId string
	helper    *zap.SugaredLogger
)

func Helper() *zap.SugaredLogger {
	return helper.With(zap.String("request_id", requestId))
}

func HelperWithContext(ctx context.Context) *zap.SugaredLogger {
	traceId := fmt.Sprintf("%s", ctx.Value(consts.ContextRequestIdKey))
	return helper.With(zap.String("request_id", traceId))
}

func WithRequestId(rid string) {
	requestId = rid
}

type Logger struct {
	zl *zap.Logger
	zs *zap.SugaredLogger

	requestId string
}

func InitLog(c Config, isLocal bool) *Logger {
	var (
		defaultLevel = zapcore.InfoLevel
		syncer       = []zapcore.WriteSyncer{
			getLogWriter(c.Zap),
		}
		options = []zap.Option{
			zap.AddStacktrace(zap.NewAtomicLevelAt(zap.ErrorLevel)),
			zap.AddCaller(),
			zap.AddCallerSkip(0),
		}
	)
	if isLocal {
		syncer = append(syncer, zapcore.AddSync(os.Stdout))
		options = append(options, zap.Development())
		defaultLevel = zapcore.DebugLevel
	}
	core := zapcore.NewCore(
		getEncoder(isLocal),
		zapcore.NewMultiWriteSyncer(syncer...),
		zap.NewAtomicLevelAt(defaultLevel),
	)
	core = zapcore.RegisterHooks(core, func(entry zapcore.Entry) error {
		if entry.Level >= zapcore.FatalLevel {
			err := errors.New(entry.Message)
			wework.ErrorBroadcast(err, requestId, entry.Stack)
		}
		return nil
	})
	zap.ReplaceGlobals(zap.New(core, options...))
	helper = zap.S()

	return &Logger{
		zl: zap.L(),
		zs: zap.S(),
	}
}

func getEncoder(isLocal bool) zapcore.Encoder {
	// 定义zap的配置
	zap.NewDevelopmentEncoderConfig()
	encoder := zapcore.EncoderConfig{
		TimeKey:        "t",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeTime:     getCustomTimeEncoder,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	if isLocal {
		return zapcore.NewConsoleEncoder(encoder)
	}
	return zapcore.NewJSONEncoder(encoder)
}

func getLogWriter(c Zap) zapcore.WriteSyncer {
	var (
		filePath = fmt.Sprintf("./%s/%s", c.FilePath, c.FileName)
	)
	fmt.Println(filePath)
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    c.MaxSize,
		MaxAge:     c.MaxAge,
		MaxBackups: c.MaxBackups,
		Compress:   false,
		LocalTime:  true,
	}

	return zapcore.AddSync(lumberJackLogger)
}

func getCustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000000"))
}

func (l *Logger) WithRequestId(rid string) *Logger {
	l.requestId = rid
	return l
}

// S 返回Zap包 SugaredLog
func (l *Logger) S() *zap.SugaredLogger {
	return l.zs.With(zap.String("request_id", l.requestId))
}

// L 返回Zap包 Logger
func (l *Logger) L() *zap.Logger {
	return l.zl.With(zap.String("request_id", l.requestId))
}
