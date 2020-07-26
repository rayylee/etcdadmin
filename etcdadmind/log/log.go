package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"strings"
)

type Config struct {
	Level string
	File  string
}

var (
	levelMap = map[string]zapcore.Level{
		"debug": zap.DebugLevel,
		"info":  zap.InfoLevel,
		"warn":  zap.WarnLevel,
		"error": zap.ErrorLevel,
	}
	globalLogger *zap.Logger
)

func Init(cfg Config) error {

	lvl := strings.ToLower(cfg.Level)

	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   cfg.File,
		MaxSize:    100, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
	})

	encoder := zap.NewProductionEncoderConfig()
	encoder.EncodeTime = zapcore.ISO8601TimeEncoder
	// encoder.EncodeCaller = zapcore.FullCallerEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoder),
		w,
		levelMap[lvl],
	)

	caller := zap.AddCaller()
	development := zap.Development()

	globalLogger = zap.New(core, caller, development)

	return nil
}

func GetLogger() *zap.Logger {
	return globalLogger
}
