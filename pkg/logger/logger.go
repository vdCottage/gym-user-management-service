package logger

import (
	"io"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps zap.Logger
type Logger struct {
	*zap.Logger
}

// NewLogger creates a new logger
func NewLogger(level, format, outputPath string) *Logger {
	// Set default values
	if level == "" {
		level = "info"
	}
	if format == "" {
		format = "json"
	}

	// Parse log level
	var logLevel zapcore.Level
	switch level {
	case "debug":
		logLevel = zapcore.DebugLevel
	case "info":
		logLevel = zapcore.InfoLevel
	case "warn":
		logLevel = zapcore.WarnLevel
	case "error":
		logLevel = zapcore.ErrorLevel
	default:
		logLevel = zapcore.InfoLevel
	}

	// Create encoder config
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Create writer
	var writer io.Writer
	if outputPath != "" {
		// Ensure directory exists
		if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
			panic(err)
		}

		// Open log file
		file, err := os.OpenFile(outputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		writer = file
	} else {
		writer = os.Stdout
	}

	// Create encoder
	var encoder zapcore.Encoder
	if format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// Create core
	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(writer),
		logLevel,
	)

	// Create logger
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return &Logger{logger}
}

// Sync flushes any buffered log entries
func (l *Logger) Sync() error {
	return l.Logger.Sync()
}

// Info logs an info message
func (l *Logger) Info(msg string, fields ...interface{}) {
	l.Logger.Info(msg, zap.Any("fields", fields))
}

// Error logs an error message
func (l *Logger) Error(msg string, fields ...interface{}) {
	l.Logger.Error(msg, zap.Any("fields", fields))
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, fields ...interface{}) {
	l.Logger.Debug(msg, zap.Any("fields", fields))
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, fields ...interface{}) {
	l.Logger.Warn(msg, zap.Any("fields", fields))
}
