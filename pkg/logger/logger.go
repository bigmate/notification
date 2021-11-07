package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.SugaredLogger

func init() {
	localEnabled := false
	// localEnabled := *flag.Bool("local-enabled", false, "--local-enabled")
	// flag.Parse()

	var zLogger *zap.Logger

	if localEnabled {
		l, err := zap.NewDevelopment()
		if err != nil {
			panic(err)
		}
		zLogger = l
	} else {
		l, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		zLogger = l
	}

	zLogger = zLogger.WithOptions(
		zap.WithCaller(false),
		zap.AddStacktrace(&stackDisabler{}),
	)
	logger = zLogger.Sugar()
}

type stackDisabler struct{}

func (s *stackDisabler) Enabled(_ zapcore.Level) bool {
	return false
}

func Debug(msg interface{}) {
	logger.Debug(msg)
}
func Info(msg interface{}) {
	logger.Info(msg)
}
func Warn(msg interface{}) {
	logger.Warn(msg)
}
func Error(msg interface{}) {
	logger.Error(msg)
}
func Fatal(msg interface{}) {
	logger.Fatal(msg)
}
func Debugf(tmpl string, args ...interface{}) {
	logger.Debugf(tmpl, args...)
}
func Infof(tmpl string, args ...interface{}) {
	logger.Infof(tmpl, args...)
}
func Warnf(tmpl string, args ...interface{}) {
	logger.Warnf(tmpl, args...)
}
func Errorf(tmpl string, args ...interface{}) {
	logger.Errorf(tmpl, args...)
}
func Fatalf(tmpl string, args ...interface{}) {
	logger.Fatalf(tmpl, args...)
}
