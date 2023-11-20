package log

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewZapLogger(opts *InitOptions) (*zap.Logger, error) {
	encoderConfig := zap.NewProductionEncoderConfig()
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	fileEncoder := zapcore.NewJSONEncoder(encoderConfig)

	var cores []zapcore.Core

	// Setup stdout logger
	if opts.StdoutFile != "" {
		stdoutFile, err := os.Create(opts.StdoutFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to open stdout log file: %v\n", err)
			return nil, err
		}
		cores = append(cores, zapcore.NewCore(fileEncoder, zapcore.AddSync(stdoutFile), zapcore.InfoLevel))
	} else {
		cores = append(cores, zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.InfoLevel))
	}

	// Setup stderr logger
	if opts.StderrFile != "" {
		stderrFile, err := os.Create(opts.StderrFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to open stderr log file: %v\n", err)
			return nil, err
		}
		cores = append(cores, zapcore.NewCore(fileEncoder, zapcore.AddSync(stderrFile), zapcore.ErrorLevel))
	} else {
		cores = append(cores, zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stderr), zapcore.ErrorLevel))
	}

	l := zap.New(zapcore.NewTee(cores...))

	return l, nil
}

func NewZapLoggerManager(l *zap.Logger) ILoggerManager {
	return &zapLoggerManager{
		logger: l,
	}
}

type zapLoggerManager struct {
	logger *zap.Logger
}

func (l *zapLoggerManager) Logger() ILogger {
	return l.logger.Sugar()
}

func (l *zapLoggerManager) Sync() error {
	return l.logger.Sync()
}
