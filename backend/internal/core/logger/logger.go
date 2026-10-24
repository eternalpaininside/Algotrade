package corelogger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger

	file *os.File
}

func FromContext(ctx context.Context) *Logger {
	log, ok := ctx.Value("logger").(*Logger)
	if !ok {
		panic("not found logger in context")
	}

	return log
}

func NewLogger(
	config LoggerConfig,
) (*Logger, error) {
	level := zap.NewAtomicLevel()
	if err := level.UnmarshalText(
		[]byte(config.Level),
	); err != nil {
		return nil, fmt.Errorf(
			"unmarshal log level: %w",
			err,
		)
	}

	if err := os.MkdirAll(
		config.Folder, 0o755,
	); err != nil {
		return nil, fmt.Errorf(
			"mkdir log: %w",
			err,
		)
	}

	timestamp := time.Now().Format("2006-01-02T15-04-05.000000")
	logPath := filepath.Join(
		config.Folder,
		fmt.Sprintf("%s.log", timestamp),
	)

	file, err := os.OpenFile(
		logPath,
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0o0644,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"open log file: %w",
			err,
		)
	}

	zapCfg := zap.NewDevelopmentEncoderConfig()
	zapCfg.EncodeTime = zapcore.TimeEncoderOfLayout(
		"2006-01-02T15:04:05.000000",
	)

	encoder := zapcore.NewConsoleEncoder(zapCfg)
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level),
		zapcore.NewCore(encoder, zapcore.AddSync(file), level),
	)

	logger := zap.New(
		core,
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.DPanicLevel),
	)

	return &Logger{
		Logger: logger,
		file:   file,
	}, nil
}

func (l *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{
		Logger: l.Logger.With(fields...),
		file:   l.file,
	}
}

func (l *Logger) Close() {
	if err := l.file.Close(); err != nil {
		fmt.Printf("failed close logger: %v", err)
	}
}
