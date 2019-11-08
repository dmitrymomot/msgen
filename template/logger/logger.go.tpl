package logger

import (
	stdlog "log"

	"go.uber.org/zap"
)

type (
	// Logger struct
	Logger struct {
		*zap.SugaredLogger
		zlogger *zap.Logger
	}

	// Options struct
	Options struct {
		Debug       bool
		ServiceName string
		BuildTag    string
	}
)

// Sync is alias of zlogger.Sync() function
// set in main.go:
//      defer zlogger.Sync()
func (l *Logger) Sync() {
	l.zlogger.Sync()
}

// DefaultLogger factory
func DefaultLogger(opt Options) *Logger {
	var zlogger *zap.Logger
	var err error
	if opt.Debug {
		zlogger, err = zap.NewDevelopment()
	} else {
		zlogger, err = zap.NewProduction()
	}
	if err != nil {
		stdlog.Fatal("can't init logger", err)
	}

	log := zlogger.Sugar().With(zap.String("service_name", opt.ServiceName), zap.String("build", opt.BuildTag))
	log.Info("starting app")
	log.Debug("debug mode enabled")

	return &Logger{log, zlogger}
}
