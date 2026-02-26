package logger

import (
	"context"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type contextKey string

const RequestIDKey contextKey = "request_id"

var (
	Log  *zap.Logger
	once sync.Once
)

type Config struct {
	Env         string
	LogFilePath string
}

func Init(cfg *Config) {
	once.Do(func() {
		Log = newLogger(cfg)
	})
}

func newLogger(cfg *Config) *zap.Logger {
	var level zapcore.Level

	if cfg.Env == "production" {
		level = zapcore.InfoLevel
	} else {
		level = zapcore.DebugLevel
	}

	outputPaths := []string{"stdout"}
	if cfg.LogFilePath != "" {
		outputPaths = append(outputPaths, cfg.LogFilePath)
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(level),
		Development:      cfg.Env != "production",
		Encoding:         "json",
		EncoderConfig:    encoderConfig,
		OutputPaths:      outputPaths,
		ErrorOutputPaths: []string{"stderr"},
	}

	log, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}

	return log
}

func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}

func WithRequestID(ctx context.Context) *zap.Logger {
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok && requestID != "" {
		return Log.With(zap.String("request_id", requestID))
	}
	return Log
}

func Debug(msg string, fields ...zap.Field) {
	Log.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	Log.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	Log.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	Log.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	Log.Fatal(msg, fields...)
	os.Exit(1)
}

func DebugCtx(ctx context.Context, msg string, fields ...zap.Field) {
	WithRequestID(ctx).Debug(msg, fields...)
}

func InfoCtx(ctx context.Context, msg string, fields ...zap.Field) {
	WithRequestID(ctx).Info(msg, fields...)
}

func WarnCtx(ctx context.Context, msg string, fields ...zap.Field) {
	WithRequestID(ctx).Warn(msg, fields...)
}

func ErrorCtx(ctx context.Context, msg string, fields ...zap.Field) {
	WithRequestID(ctx).Error(msg, fields...)
}

func With(fields ...zap.Field) *zap.Logger {
	return Log.With(fields...)
}
