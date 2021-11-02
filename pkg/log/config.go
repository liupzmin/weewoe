package log

import (
	"fmt"
	"strconv"
	"time"

	"github.com/liupzmin/weewoe/pkg/util/xcolor"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config ...
type Config struct {
	// Dir 日志输出目录
	Dir string
	// Name 日志文件名称
	Name string
	// Level 日志初始等级
	Level string
	// 日志初始化字段
	Fields []zap.Field
	// 是否添加调用者信息
	AddCaller bool
	// 日志前缀
	Prefix string
	// 日志输出文件最大长度，超过改值则截断
	MaxSize   int
	MaxAge    int
	MaxBackup int
	// 日志磁盘刷盘间隔
	Interval   time.Duration
	CallerSkip int
	// 异步
	Async         bool
	asyncCloser   func() error
	Queue         bool
	QueueSleep    time.Duration
	Core          zapcore.Core
	EnableConsole bool
	EncoderConfig *zapcore.EncoderConfig
	configKey     string
	// 是否启用可读时间
	EnableTimeLayout bool
}

// Filename ...
func (config *Config) Filename() string {
	return fmt.Sprintf("%s/%s", config.Dir, config.Name)
}

// DefaultZapEncoderConfig ...
func DefaultZapEncoderConfig() *zapcore.EncoderConfig {
	return &zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "lv",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     timeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   ShortCallerEncoder,
	}
}

func NewLogger(config *Config) *Logger {

	if config.EncoderConfig == nil {
		config.EncoderConfig = DefaultZapEncoderConfig()
	}
	if config.EnableConsole {
		config.EncoderConfig.EncodeLevel = DebugEncodeLevel
	}

	var (
		c int8 = 0
		t int8 = 0
	)
	if config.EnableConsole {
		c = 1
	}
	if config.EnableTimeLayout {
		t = 1
	}
	switch c<<1 + t {
	case 0:
		config.EncoderConfig.EncodeTime = TimeEncoder
	case 1:
		config.EncoderConfig.EncodeTime = TimeLayoutEncoder
	case 2:
		config.EncoderConfig.EncodeTime = TimeEncoderColor
	case 3:
		config.EncoderConfig.EncodeTime = TimeLayoutEncoderColor
	}

	return newLogger(config)
}

// ShortCallerEncoder serializes a caller in package/file:line format, trimming
// all but the final directory from the full path.
func ShortCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	// TODO: consider using a byte-oriented API to save an allocation.
	enc.AppendString(fmt.Sprintf("%-32s", caller.TrimmedPath()))
}

func TimeLayoutEncoderColor(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(xcolor.Yellow(t.Format("2006-01-02 15:04:05.000")))
}

func TimeLayoutEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func TimeEncoderColor(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(xcolor.Yellow(strconv.FormatInt(t.Unix(), 10)))
}

func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendInt64(t.Unix())
}
