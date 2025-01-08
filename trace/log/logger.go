package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"sync"
)

const logLevelEnvVar string = "OTEL_LOG_LEVEL"

var (
	defaultLogger *zap.Logger
	level         zap.AtomicLevel
	levellock     sync.Mutex
)

// ZapLogger is a wrapper around zap.Logger that adds a traceID field to all logs.
type ZapLogger struct {
	zap *zap.Logger
}

func init() {
	level = zap.NewAtomicLevel()
	//level.SetLevel(zapcore.InfoLevel)

	envLevel, exists := os.LookupEnv(logLevelEnvVar)
	// if level is wrong, default to Info
	if exists {
		SetLogLevel(envLevel)
	} else {
		level.SetLevel(zapcore.InfoLevel) // Default level
	}
	config := zap.Config{
		Level:            level,
		Development:      false,
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	var err error
	defaultLogger, err = config.Build(zap.AddCaller(), zap.AddCallerSkip(1))
	if err != nil {
		log.Printf("Failed to init log: %v\n", err)
	}

	//GlobalLogger = NewLogger()
}

// SetLogLevel sets the log level for the log.
// Valid levels are:   debug | info | warn | error
// If the level is invalid, the log level will default to Info.
func SetLogLevel(newLevel string) {
	levellock.Lock()
	defer levellock.Unlock()

	var l zapcore.Level
	if err := l.UnmarshalText([]byte(newLevel)); err != nil {
		// if level is wrong, default to Info
		level.SetLevel(zapcore.InfoLevel)
	}

	level.SetLevel(l)
}

// NewLogger creates and returns a new ZapLogger instance.
// The newly created log will use the defaultLogger as its underlying Zap log.
func NewLogger() *ZapLogger {
	return &ZapLogger{zap: defaultLogger.With(zap.String("type", "trace"))}
}

// WithTrace returns a new ZapLogger instance that includes additional fields
// for tracing. Specifically, it adds a "type" field set to "trace" and a "traceID"
// field set to the provided traceID string.
// This function is useful for generating a trace-specific log from a generic ZapLogger instance.
func (l *ZapLogger) WithTrace(traceID string) *ZapLogger {
	newZapLogger := l.zap.With(zap.String("type", "trace"), zap.String("traceID", traceID))
	return &ZapLogger{zap: newZapLogger}
}

// Debug logs a debug-level message.
func (l *ZapLogger) Debug(message string, fields ...zap.Field) {
	l.zap.Debug(message, fields...)
}

// Info logs an info-level message.
func (l *ZapLogger) Info(message string, fields ...zap.Field) {
	l.zap.Info(message, fields...)
}

// Warn logs a warning-level message.
func (l *ZapLogger) Warn(message string, fields ...zap.Field) {
	l.zap.Warn(message, fields...)
}

// Error logs an error-level message.
func (l *ZapLogger) Error(message string, fields ...zap.Field) {
	l.zap.Error(message, fields...)
}

func (l *ZapLogger) Debugf(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	l.zap.Debug(msg)
}
func (l *ZapLogger) Infof(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	l.zap.Info(msg)
}
func (l *ZapLogger) Warnf(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	l.zap.Warn(msg)
}

func (l *ZapLogger) Errorf(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	l.zap.Error(msg)
}
