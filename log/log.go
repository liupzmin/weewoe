// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of weewoe

package log

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/liupzmin/weewoe/util/xcolor"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel = zap.DebugLevel
	// InfoLevel is the default logging priority.
	InfoLevel = zap.InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel = zap.WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-Level logs.
	ErrorLevel = zap.ErrorLevel
	// PanicLevel logs a message, then panics.
	PanicLevel = zap.PanicLevel
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel = zap.FatalLevel
)

// Func ...
type (
	Func   func(string, ...zap.Field)
	Field  = zap.Field
	Level  = zapcore.Level
	Logger struct {
		desugar *zap.Logger
		lv      *zap.AtomicLevel
		config  Config
		sugar   *zap.SugaredLogger
	}
)

var (
	// String ...
	String = zap.String
	// Any ...
	Any = zap.Any
	// Int64 ...
	Int64 = zap.Int64
	// Int ...
	Int = zap.Int
	// Int32 ...
	Int32 = zap.Int32
	// Uint ...
	Uint = zap.Uint
	// Duration ...
	Duration = zap.Duration
	// Durationp ...
	Durationp = zap.Durationp
	// Object ...
	Object = zap.Object
	// Namespace ...
	Namespace = zap.Namespace
	// Reflect ...
	Reflect = zap.Reflect
	// Skip ...
	Skip = zap.Skip()
	// ByteString ...
	ByteString = zap.ByteString
	// Binary
	Binary = zap.Binary
)

func newLogger(config *Config) *Logger {
	zapOptions := make([]zap.Option, 0)
	zapOptions = append(zapOptions, zap.AddStacktrace(zap.ErrorLevel))
	if config.AddCaller {
		zapOptions = append(zapOptions, zap.AddCaller(), zap.AddCallerSkip(config.CallerSkip))
	}
	if len(config.Fields) > 0 {
		zapOptions = append(zapOptions, zap.Fields(config.Fields...))
	}

	var ws zapcore.WriteSyncer
	if config.EnableConsole {
		ws = os.Stdout
	} else {
		ws = zapcore.AddSync(newRotate(config))
	}

	if config.Async {
		ws, config.asyncCloser = Buffer(ws, defaultBufferSize, defaultFlushInterval)
	}

	lv := zap.NewAtomicLevelAt(zapcore.InfoLevel)
	if err := lv.UnmarshalText([]byte(config.Level)); err != nil {
		panic(err)
	}

	encoderConfig := *config.EncoderConfig
	core := config.Core
	if core == nil {
		core = zapcore.NewCore(
			func() zapcore.Encoder {
				if config.EnableConsole {
					return zapcore.NewConsoleEncoder(encoderConfig)
				}
				return zapcore.NewJSONEncoder(encoderConfig)
			}(),
			ws,
			lv,
		)
	}

	zapLogger := zap.New(
		core,
		zapOptions...,
	)
	return &Logger{
		desugar: zapLogger,
		lv:      &lv,
		config:  *config,
		sugar:   zapLogger.Sugar(),
	}
}

// AutoLevel ...
func (logger *Logger) AutoLevel(level string) {
	lvText := strings.ToLower(level)
	if lvText != "" {
		logger.Info("update level", String("level", lvText), String("name", logger.config.Name))
		logger.lv.UnmarshalText([]byte(lvText))
	}
}

// SetLevel ...
func (logger *Logger) SetLevel(lv Level) {
	logger.lv.SetLevel(lv)
}

// Stop ...
func (logger *Logger) Stop() error {
	if logger.config.asyncCloser != nil {
		err := logger.config.asyncCloser()
		if err != nil {
			return err
		}
	}

	if !logger.config.EnableConsole {
		return logger.desugar.Sync()
	}
	return nil
}

// DebugEncodeLevel ...
func DebugEncodeLevel(lv zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	var colorize = xcolor.Red
	switch lv {
	case zapcore.DebugLevel:
		colorize = xcolor.Blue
	case zapcore.InfoLevel:
		colorize = xcolor.Green
	case zapcore.WarnLevel:
		colorize = xcolor.Yellow
	case zapcore.ErrorLevel, zap.PanicLevel, zap.DPanicLevel, zap.FatalLevel:
		colorize = xcolor.Red
	default:
	}
	enc.AppendString(colorize(lv.CapitalString()))
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendInt64(t.Unix())
}

func normalizeMessage(msg string) string {
	return fmt.Sprintf("%-32s", msg)
}

// Debug ...
func (logger *Logger) Debug(msg string, fields ...Field) {
	logger.desugar.Debug(msg, fields...)
}

// Debugw ...
func (logger *Logger) Debugw(msg string, keysAndValues ...interface{}) {
	logger.sugar.Debugw(msg, keysAndValues...)
}

func sprintf(template string, args ...interface{}) string {
	msg := template
	if msg == "" && len(args) > 0 {
		msg = fmt.Sprint(args...)
	} else if msg != "" && len(args) > 0 {
		msg = fmt.Sprintf(template, args...)
	}
	return msg
}

// StdLog ...
func (logger *Logger) StdLog() *log.Logger {
	return zap.NewStdLog(logger.desugar)
}

// Debugf ...
func (logger *Logger) Debugf(template string, args ...interface{}) {
	logger.sugar.Debugw(sprintf(template, args...))
}

// Info ...
func (logger *Logger) Info(msg string, fields ...Field) {
	logger.desugar.Info(msg, fields...)
}

// Infow ...
func (logger *Logger) Infow(msg string, keysAndValues ...interface{}) {
	logger.sugar.Infow(msg, keysAndValues...)
}

// Infof ...
func (logger *Logger) Infof(template string, args ...interface{}) {
	logger.sugar.Infof(sprintf(template, args...))
}

// Warn ...
func (logger *Logger) Warn(msg string, fields ...Field) {
	logger.desugar.Warn(msg, fields...)
}

// Warnw ...
func (logger *Logger) Warnw(msg string, keysAndValues ...interface{}) {
	logger.sugar.Warnw(msg, keysAndValues...)
}

// Warnf ...
func (logger *Logger) Warnf(template string, args ...interface{}) {
	logger.sugar.Warnf(sprintf(template, args...))
}

// Error ...
func (logger *Logger) Error(msg string, fields ...Field) {
	logger.desugar.Error(msg, fields...)
}

// Errorw ...
func (logger *Logger) Errorw(msg string, keysAndValues ...interface{}) {
	logger.sugar.Errorw(msg, keysAndValues...)
}

// Errorf ...
func (logger *Logger) Errorf(template string, args ...interface{}) {
	logger.sugar.Errorf(sprintf(template, args...))
}

// Panic ...
func (logger *Logger) Panic(msg string, fields ...Field) {
	logger.desugar.Panic(msg, fields...)
}

// Panicw ...
func (logger *Logger) Panicw(msg string, keysAndValues ...interface{}) {
	logger.sugar.Panicw(msg, keysAndValues...)
}

// Panicf ...
func (logger *Logger) Panicf(template string, args ...interface{}) {
	logger.sugar.Panicf(sprintf(template, args...))
}

// DPanic ...
func (logger *Logger) DPanic(msg string, fields ...Field) {
	logger.desugar.DPanic(msg, fields...)
}

// DPanicw ...
func (logger *Logger) DPanicw(msg string, keysAndValues ...interface{}) {
	logger.sugar.DPanicw(msg, keysAndValues...)
}

// DPanicf ...
func (logger *Logger) DPanicf(template string, args ...interface{}) {
	logger.sugar.DPanicf(sprintf(template, args...))
}

// Fatal ...
func (logger *Logger) Fatal(msg string, fields ...Field) {
	logger.desugar.Fatal(msg, fields...)
}

// Fatalw ...
func (logger *Logger) Fatalw(msg string, keysAndValues ...interface{}) {
	logger.sugar.Fatalw(msg, keysAndValues...)
}

// Fatalf ...
func (logger *Logger) Fatalf(template string, args ...interface{}) {
	logger.sugar.Fatalf(sprintf(template, args...))
}

func panicDetail(msg string, fields ...Field) {
	enc := zapcore.NewMapObjectEncoder()
	for _, field := range fields {
		field.AddTo(enc)
	}

	// 控制台输出
	fmt.Printf("%s: \n    %s: %s\n", xcolor.Red("panic"), xcolor.Red("msg"), msg)
	if _, file, line, ok := runtime.Caller(3); ok {
		fmt.Printf("    %s: %s:%d\n", xcolor.Red("loc"), file, line)
	}
	for key, val := range enc.Fields {
		fmt.Printf("    %s: %s\n", xcolor.Red(key), fmt.Sprintf("%+v", val))
	}

}

// With ...
func (logger *Logger) With(fields ...Field) *Logger {
	desugarLogger := logger.desugar.With(fields...)
	return &Logger{
		desugar: desugarLogger,
		lv:      logger.lv,
		sugar:   desugarLogger.Sugar(),
		config:  logger.config,
	}
}
