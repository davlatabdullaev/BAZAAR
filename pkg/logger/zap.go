package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newZapLogger(namespace string) *zap.Logger {

	stdout := zapcore.AddSync(os.Stdout)

	level := zap.NewAtomicLevelAt(zap.InfoLevel)

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel=zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)

	consoleEncoder.OpenNamespace(namespace)

	core := zapcore.NewTee(zapcore.NewCore(consoleEncoder, stdout, level))

	newZap := zap.New(core)

	return newZap

}
