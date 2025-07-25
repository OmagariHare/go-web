package utils

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger

func InitLogger() {
	// 配置日志轮转
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./logs/app.log",
		MaxSize:    10,    // 每个日志文件最大10MB
		MaxBackups: 5,     // 最多保留5个备份
		MaxAge:     30,    // 最多保留30天
		Compress:   false, // 不压缩
	}

	// 创建zap配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 创建核心
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(lumberJackLogger), zapcore.AddSync(os.Stdout)),
		zapcore.DebugLevel,
	)

	// 创建logger
	Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}

func SyncLogger() {
	if Logger != nil {
		if err := Logger.Sync(); err != nil {
			Logger.Error("failed to sync logger", zap.Error(err))
		}
	}
}