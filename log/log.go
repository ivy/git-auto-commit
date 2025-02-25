package log

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	mu          sync.RWMutex
	logger      *zap.SugaredLogger
	atomicLevel zap.AtomicLevel
)

func init() {
	atomicLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
	cfg := zap.NewDevelopmentConfig()
	cfg.Level = atomicLevel

	l, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	logger = l.Sugar()
}

// SetLogger allows replacing the default logger with a custom one.
func SetLogger(l *zap.SugaredLogger) {
	mu.Lock()
	defer mu.Unlock()
	logger = l
}

// SetLevel allows changing the log level dynamically.
func SetLevel(level zapcore.Level) {
	atomicLevel.SetLevel(level)
}

// getLogger safely returns the current logger.
func getLogger() *zap.SugaredLogger {
	mu.RLock()
	defer mu.RUnlock()
	return logger
}

// Debug logs a debug-level message with optional key/value pairs.
func Debug(args ...interface{}) {
	getLogger().Debug(args...)
}

// Debugf logs a debug-level message with formatting.
func Debugf(template string, args ...interface{}) {
	getLogger().Debugf(template, args...)
}

// Debugw logs a debug-level message with optional structured context.
func Debugw(msg string, keysAndValues ...interface{}) {
	getLogger().Debugw(msg, keysAndValues...)
}

// Info logs an info-level message with optional key/value pairs.
func Info(args ...interface{}) {
	getLogger().Info(args...)
}

// Infof logs an info-level message with formatting.
func Infof(template string, args ...interface{}) {
	getLogger().Infof(template, args...)
}

// Infow logs an info-level message with optional structured context.
func Infow(msg string, keysAndValues ...interface{}) {
	getLogger().Infow(msg, keysAndValues...)
}

// Warn logs a warn-level message with optional key/value pairs.
func Warn(args ...interface{}) {
	getLogger().Warn(args...)
}

// Warnf logs a warn-level message with formatting.
func Warnf(template string, args ...interface{}) {
	getLogger().Warnf(template, args...)
}

// Warnw logs a warn-level message with optional structured context.
func Warnw(msg string, keysAndValues ...interface{}) {
	getLogger().Warnw(msg, keysAndValues...)
}

// Error logs an error-level message with optional key/value pairs.
func Error(args ...interface{}) {
	getLogger().Error(args...)
}

// Errorf logs an error-level message with formatting.
func Errorf(template string, args ...interface{}) {
	getLogger().Errorf(template, args...)
}

// Errorw logs an error-level message with optional structured context.
func Errorw(msg string, keysAndValues ...interface{}) {
	getLogger().Errorw(msg, keysAndValues...)
}

// Panic logs a panic-level message, then calls panic.
func Panic(args ...interface{}) {
	getLogger().Panic(args...)
}

// Panicf logs a panic-level message with formatting, then calls panic.
func Panicf(template string, args ...interface{}) {
	getLogger().Panicf(template, args...)
}

// Panicw logs a panic-level message with optional structured context, then calls panic.
func Panicw(msg string, keysAndValues ...interface{}) {
	getLogger().Panicw(msg, keysAndValues...)
}

// Fatal logs a fatal-level message, then calls os.Exit(1).
func Fatal(args ...interface{}) {
	getLogger().Fatal(args...)
}

// Fatalf logs a fatal-level message with formatting, then calls os.Exit(1).
func Fatalf(template string, args ...interface{}) {
	getLogger().Fatalf(template, args...)
}

// Fatalw logs a fatal-level message with optional structured context, then calls os.Exit(1).
func Fatalw(msg string, keysAndValues ...interface{}) {
	getLogger().Fatalw(msg, keysAndValues...)
}

// Sync flushes any buffered log entries.
func Sync() error {
	return getLogger().Sync()
}
