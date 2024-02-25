package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Field = zapcore.Field

var (
	Int    = zap.Int
	String = zap.String
	Error  = zap.Error
	Any    = zap.Any
)