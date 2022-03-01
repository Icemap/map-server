package logger

import (
	"fmt"
	"io"
	"log"
	"map-server/config"
	"os"
	"sync"
	"time"

	rotate "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger    *zap.Logger
	logConfig config.Log
	once      sync.Once
)

// initLog logger init
func initLog() {
	logConfig = config.ReadConfig().Log

	zapConfig := zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		EncodeLevel: zapcore.CapitalLevelEncoder,
		TimeKey:     "ts",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		CallerKey:    "file",
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	}
	encoder := zapcore.NewConsoleEncoder(zapConfig)
	logLevel := getLevel()

	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.WarnLevel && lvl >= logLevel
	})

	warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel && lvl >= logLevel
	})

	infoWriter := getWriter(logConfig.LogName)
	warnWriter := getWriter(logConfig.ErrorName)

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(warnWriter), warnLevel),
		zapcore.NewCore(zapcore.NewConsoleEncoder(zapConfig),
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), logLevel),
	)
	logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.WarnLevel))
}

func getLevel() zapcore.Level {
	switch logConfig.Level {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "panic":
		return zap.PanicLevel
	case "fatal":
		return zap.FatalLevel
	default:
		fmt.Printf("[Warning] log level [%s] confused, use info level instead", logConfig.Level)
		return zap.InfoLevel
	}
}

func getWriter(filename string) io.Writer {
	hook, err := rotate.New(
		filename+".%Y%m%d%H",
		rotate.WithLinkName(filename),
		rotate.WithMaxAge(time.Hour*24*30),
		rotate.WithRotationTime(time.Hour*24),
	)
	if err != nil {
		log.Println("log error")
		panic(err)
	}
	return hook
}

func Logger() *zap.SugaredLogger {
	once.Do(initLog)
	return logger.Sugar()
}
