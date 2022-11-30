package log

import (
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Field = zap.Field

var (
	Logger  *zap.Logger
	String  = zap.String
	Any     = zap.Any
	Int     = zap.Int
	Float32 = zap.Float32
)

func InitLogger(logpath string) {
	hook := &lumberjack.Logger{
		Filename:   logpath,
		MaxSize:    100,
		MaxAge:     7,
		MaxBackups: 30,
		Compress:   true,
	}
	writeSyncer := zapcore.AddSync(hook)
	encoder := zapcore.NewJSONEncoder(
		zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			MessageKey:     "msg",
			FunctionKey:    "Function",
			StacktraceKey:  "stacktrace",
			CallerKey:      "Caller",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	)

	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zapcore.DebugLevel)

	var writes = []zapcore.WriteSyncer{writeSyncer}
	writes = append(writes, zapcore.AddSync(os.Stdout))
	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(writes...),
		zapcore.DebugLevel,
	)
	caller := zap.AddCaller()
	development := zap.Development()
	file := zap.Fields(zap.String("log", "go-chat"))
	Logger = zap.New(core, caller, development, file)
	Logger.Info("Start Logging")

}
