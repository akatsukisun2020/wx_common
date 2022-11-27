package logger

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 日志建设，参考： https://www.jb51.net/article/244746.htm
// 有时间详细看看封装过程.

var sugarLogger *zap.SugaredLogger

func init() {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./log/grpc.log",
		MaxSize:    50,
		MaxBackups: 10,
		MaxAge:     30,
		Compress:   false,
	}
	logWriter := zapcore.AddSync(lumberJackLogger)

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	core := zapcore.NewCore(encoder, logWriter, zapcore.DebugLevel)
	logger := zap.New(core, zap.AddCaller())
	sugarLogger = logger.Sugar()
}

func Debugf(template string, args ...interface{}) {
	sugarLogger.Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	sugarLogger.Infof(template, args...)
}

func Errorf(template string, args ...interface{}) {
	sugarLogger.Errorf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	sugarLogger.Fatalf(template, args...)
}
