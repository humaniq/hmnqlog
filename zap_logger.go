// Package hmnqlog ...
package hmnqlog

import (
	"errors"
	"os"

	"github.com/getsentry/raven-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type logger struct {
	zl     *zap.Logger
	sentry *raven.Client
}

// Logger ...
type Logger interface {
	Debug(string, ...zapcore.Field)
	Info(string, ...zapcore.Field)
	Warn(string, ...zapcore.Field)
	Error(string, ...zapcore.Field)
	Fatal(string, ...zapcore.Field)
}

// ZapOptions ...
type ZapOptions struct {
	AppName     string
	AppEnv      string
	AppRevision string
	Hostname    string
	PID         int
	LogLevel    zapcore.Level
}

// NewZapLogger ...
func NewZapLogger(zo ZapOptions) (Logger, error) {
	if zo.AppEnv == "" {
		return nil, errors.New("Options.Env must be set to use this logger")
	}

	if zo.LogLevel == 0 {
		if zo.AppEnv == "staging" || zo.AppEnv == "production" {
			zo.LogLevel = zap.InfoLevel
		} else {
			zo.LogLevel = zap.DebugLevel
		}
	}

	if zo.Hostname == "" {
		var err error
		zo.Hostname, err = os.Hostname()
		if err != nil {
			return nil, err
		}
	}

	if zo.PID == 0 {
		zo.PID = os.Getpid()
	}

	context := []zapcore.Field{
		zap.String("app_name", zo.AppName),
		zap.String("hostname", zo.Hostname),
		zap.String("version", zo.AppRevision),
		zap.String("env", zo.AppEnv),
		zap.Int("PID", zo.PID),
	}

	var zl *zap.Logger
	zl, zErr := zap.NewProduction(zap.Fields(context...))
	if zErr != nil {
		return nil, zErr
	}

	return &logger{zl: zl}, nil
}

// Debug logs the given message that ends with an os.Exit(1). If a sentry
// client was set it will report the given message prior to logging and
// subsequently calling os.Exit(1)
func (l *logger) Debug(message string, fields ...zapcore.Field) {
	l.zl.Debug(message, fields...)
}

// Info logs the given message and fields with a log level of Warn and sends
// the given message and fields to sentry if a sentry client is set.
func (l *logger) Info(message string, fields ...zapcore.Field) {
	l.zl.Info(message, fields...)
}

// Warn logs the given message and fields with a log level of Warn and sends
// the given message and fields to sentry if a sentry client is set.
func (l *logger) Warn(message string, fields ...zapcore.Field) {
	l.zl.Warn(message, fields...)
}

// Error logs the given message that ends with an os.Exit(1). If a sentry
// client was set it will report the given message prior to logging and
// subsequently calling os.Exit(1)
func (l *logger) Error(message string, fields ...zapcore.Field) {
	l.zl.Error(message, fields...)
}

// Fatal logs the given message that ends with an os.Exit(1). If a sentry
// client was set it will report the given message prior to logging and
// subsequently calling os.Exit(1)
func (l *logger) Fatal(message string, fields ...zapcore.Field) {
	l.zl.Fatal(message, fields...)
}
