package log

import (
	"errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"syscall"
)

type Logger interface {
	InitLogger()
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	DPanic(args ...interface{})
	DPanicf(template string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
}

type apiLogger struct {
	sugar *zap.SugaredLogger
}

func NewLogger() Logger {
	return &apiLogger{}
}

func (l *apiLogger) InitLogger() {
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.AddSync(os.Stderr),
		zap.NewAtomicLevelAt(zapcore.InfoLevel))
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	l.sugar = logger.Sugar()
	defer func(sugar *zap.SugaredLogger) {
		err := l.sugar.Sync()
		if err != nil && !errors.Is(err, syscall.ENOTTY) {
			l.sugar.Error(err)
		}
	}(l.sugar)
}

func (l *apiLogger) Debug(args ...interface{}) {
	l.sugar.Debug(args...)
}

func (l *apiLogger) Debugf(template string, args ...interface{}) {
	l.sugar.Debugf(template, args...)
}

func (l *apiLogger) Info(args ...interface{}) {
	l.sugar.Info(args...)
}

func (l *apiLogger) Infof(template string, args ...interface{}) {
	l.sugar.Infof(template, args...)
}

func (l *apiLogger) Warn(args ...interface{}) {
	l.sugar.Warn(args...)
}

func (l *apiLogger) Warnf(template string, args ...interface{}) {
	l.sugar.Warnf(template, args...)
}

func (l *apiLogger) Error(args ...interface{}) {
	l.sugar.Error(args...)
}

func (l *apiLogger) Errorf(template string, args ...interface{}) {
	l.sugar.Errorf(template, args...)
}

func (l *apiLogger) DPanic(args ...interface{}) {
	l.sugar.DPanic(args...)
}

func (l *apiLogger) DPanicf(template string, args ...interface{}) {
	l.sugar.DPanicf(template, args...)
}

func (l *apiLogger) Panic(args ...interface{}) {
	l.sugar.Panic(args...)
}

func (l *apiLogger) Panicf(template string, args ...interface{}) {
	l.sugar.Panicf(template, args...)
}

func (l *apiLogger) Fatal(args ...interface{}) {
	l.sugar.Fatal(args...)
}

func (l *apiLogger) Fatalf(template string, args ...interface{}) {
	l.sugar.Fatalf(template, args...)
}
