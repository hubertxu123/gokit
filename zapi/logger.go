package zapi

import (
	"context"
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"time"
)

// zap improved
const defaultTraceIdKey = "traceId"

type Logger struct {
	origin     *zap.Logger
	ctx        context.Context
	traceIdKey string
	traceId    string
}

func New(origin *zap.Logger, opts ...Option) *Logger {
	l := &Logger{
		origin:     origin,
		traceIdKey: defaultTraceIdKey,
	}
	for _, opt := range opts {
		opt(l)
	}
	return l
}

func (l *Logger) clone() *Logger {
	cp := *l
	return &cp
}

func (l *Logger) WithContext(ctx context.Context) *Logger {
	zl := l.clone()
	if traceId, ok := ctx.Value(zl.traceIdKey).(string); ok {
		zl.traceId = traceId
		zl.origin.With(zap.String(zl.traceIdKey, zl.traceId))
	}
	return zl
}

func (l *Logger) WithTraceId(traceId string) *Logger {
	zl := l.clone()
	zl.traceId = traceId
	zl.origin.With(zap.String(zl.traceIdKey, zl.traceId))
	return zl
}

func (l *Logger) GenTraceLogger() *Logger {
	return l.WithTraceId(fmt.Sprintf("%s", uuid.NewV4()))
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.origin.Debug(msg, fields...)
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.origin.Info(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.origin.Warn(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.origin.Error(msg, fields...)
}

func (l *Logger) DPanic(msg string, fields ...zap.Field) {
	l.origin.DPanic(msg, fields...)
}

func (l *Logger) Panic(msg string, fields ...zap.Field) {
	l.origin.Panic(msg, fields...)
}

func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.origin.Fatal(msg, fields...)
}

var DefaultLogger = Default()

func Default() *Logger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	origin := zap.New(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.NewMultiWriteSyncer(
				zapcore.AddSync(zapcore.AddSync(os.Stdout)),
				//zapcore.AddSync(getWriter(filename)),
			),
			zap.DebugLevel,
		),
		zap.AddCaller(),
		zap.AddCallerSkip(2),
		zap.AddStacktrace(zapcore.WarnLevel),
	)
	return New(origin)
}

func getWriter(filename string) io.Writer {
	hook, err := rotatelogs.New(
		filename+".%Y%m%d",
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Hour*24*7),
		rotatelogs.WithRotationTime(time.Hour*24),
	)
	if err != nil {
		panic(err)
	}
	return hook
}

func Debug(msg string, fields ...zap.Field) {
	DefaultLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	DefaultLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	DefaultLogger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	DefaultLogger.Error(msg, fields...)
}

func DPanic(msg string, fields ...zap.Field) {
	DefaultLogger.DPanic(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	DefaultLogger.Panic(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	DefaultLogger.Fatal(msg, fields...)
}
