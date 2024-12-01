package zlog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerWriter struct {
	out   *zap.Logger
	level zapcore.Level
}

func (k *LoggerWriter) Update(logger ...*zap.Logger) {
	if len(logger) == 0 {
		k.out = DefaultLogger
		return
	}
	k.out = logger[0]
}

func (lw *LoggerWriter) Write(p []byte) (n int, err error) {
	switch lw.level {
	case zapcore.ErrorLevel:
		lw.out.Error(string(p))
	case zapcore.WarnLevel:
		lw.out.Warn(string(p))
	case zapcore.InfoLevel:
		lw.out.Info(string(p))
	default:
		lw.out.Debug(string(p))
	}
	return len(p), nil
}
func NewWriter(logger *zap.Logger, level zapcore.Level) *LoggerWriter {
	return &LoggerWriter{
		out:   logger,
		level: level,
	}
}
