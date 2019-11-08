package main

import (
	stdlog "log"
	"go.uber.org/zap"
)

// Logger struct
type Logger struct {
    *zap.SugaredLogger
    zlogger *zap.Logger
}

// Sync is alias of zlogger.Sync() function
// set in main.go:
//      defer zlogger.Sync()
func (l *Logger) Sync() {
    l.zlogger.Sync()
}

func defaultLogger(debug bool) *Logger {
	var zlogger *zap.Logger
	var err error
	if debug {
		zlogger, err = zap.NewDevelopment()
	} else {
		zlogger, err = zap.NewProduction()
	}
	if err != nil {
		stdlog.Fatal("can't init logger", err)
	}

	log := zlogger.Sugar().With(zap.String("service_name", serviceName), zap.String("build", buildTag))
	log.Info("starting app")
	log.Debug("debug mode enabled")

    return &Logger{log, zlogger}
}