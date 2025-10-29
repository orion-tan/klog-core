package utils

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"klog-backend/internal/config"

	"github.com/DeRuina/timberjack"
)

var Logger *zap.Logger
var SugarLogger *zap.SugaredLogger

func InitLogger() {
	encoder := getEncoder()
	writer := getLogWriter()
	core := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(writer, zapcore.AddSync(os.Stdout)), getLogLevel())
	Logger = zap.New(core, zap.AddCaller())
	SugarLogger = Logger.Sugar()
	defer Logger.Sync()
	defer SugarLogger.Sync()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter() zapcore.WriteSyncer {

	writer := &timberjack.Logger{
		Filename:    config.Cfg.Logger.Path,
		MaxSize:     config.Cfg.Logger.MaxSize,
		MaxBackups:  config.Cfg.Logger.MaxBackups,
		MaxAge:      config.Cfg.Logger.MaxAge,
		Compression: "none",
	}
	return zapcore.AddSync(writer)
}

func getLogLevel() zapcore.Level {
	switch config.Cfg.Logger.Level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	}
	return zapcore.InfoLevel
}
