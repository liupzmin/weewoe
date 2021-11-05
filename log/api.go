package log

import (
	"errors"

	"google.golang.org/grpc/grpclog"
)

// defaultLogger default logger
// Biz Log
// debug=true as default, will be
var defaultLogger = NewLogger(&Config{
	EnableConsole:    true,
	EnableTimeLayout: true,
	Level:            "debug",
})

func SetLogger(logger *Logger) {
	defaultLogger = logger
}

// SetGRPCLogger sets grpclog's logger using defaultLogger
func SetGRPCLogger() error {
	if defaultLogger != nil {
		grpclog.SetLoggerV2(&loggerWrapper{logger: defaultLogger})
		return nil
	}
	return errors.New("default logger is nil")
}

// Info ...
func Info(msg string, fields ...Field) {
	defaultLogger.Info(msg, fields...)
}

// Debug ...
func Debug(msg string, fields ...Field) {
	defaultLogger.Debug(msg, fields...)
}

// Warn ...
func Warn(msg string, fields ...Field) {
	defaultLogger.Warn(msg, fields...)
}

// Error ...
func Error(msg string, fields ...Field) {
	defaultLogger.Error(msg, fields...)
}

// Panic ...
func Panic(msg string, fields ...Field) {
	defaultLogger.Panic(msg, fields...)
}

// DPanic ...
func DPanic(msg string, fields ...Field) {
	defaultLogger.DPanic(msg, fields...)
}

// Fatal ...
func Fatal(msg string, fields ...Field) {
	defaultLogger.Fatal(msg, fields...)
}

// Debugw ...
func Debugw(msg string, keysAndValues ...interface{}) {
	defaultLogger.Debugw(msg, keysAndValues...)
}

// Infow ...
func Infow(msg string, keysAndValues ...interface{}) {
	defaultLogger.Infow(msg, keysAndValues...)
}

// Warnw ...
func Warnw(msg string, keysAndValues ...interface{}) {
	defaultLogger.Warnw(msg, keysAndValues...)
}

// Errorw ...
func Errorw(msg string, keysAndValues ...interface{}) {
	defaultLogger.Errorw(msg, keysAndValues...)
}

// Panicw ...
func Panicw(msg string, keysAndValues ...interface{}) {
	defaultLogger.Panicw(msg, keysAndValues...)
}

// DPanicw ...
func DPanicw(msg string, keysAndValues ...interface{}) {
	defaultLogger.DPanicw(msg, keysAndValues...)
}

// Fatalw ...
func Fatalw(msg string, keysAndValues ...interface{}) {
	defaultLogger.Fatalw(msg, keysAndValues...)
}

// Debugf ...
func Debugf(msg string, args ...interface{}) {
	defaultLogger.Debugf(msg, args...)
}

// Infof ...
func Infof(msg string, args ...interface{}) {
	defaultLogger.Infof(msg, args...)
}

// Warnf ...
func Warnf(msg string, args ...interface{}) {
	defaultLogger.Warnf(msg, args...)
}

// Errorf ...
func Errorf(msg string, args ...interface{}) {
	defaultLogger.Errorf(msg, args...)
}

// Panicf ...
func Panicf(msg string, args ...interface{}) {
	defaultLogger.Panicf(msg, args...)
}

// DPanicf ...
func DPanicf(msg string, args ...interface{}) {
	defaultLogger.DPanicf(msg, args...)
}

// Fatalf ...
func Fatalf(msg string, args ...interface{}) {
	defaultLogger.Fatalf(msg, args...)
}

// Log ...
func (fn Func) Log(msg string, fields ...Field) {
	fn(msg, fields...)
}

// With ...
func With(fields ...Field) *Logger {
	return defaultLogger.With(fields...)
}

func Stop() error {
	return defaultLogger.Stop()
}
