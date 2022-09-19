package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func getEncoder(isJSON bool) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	if isJSON {
		return zapcore.NewJSONEncoder(encoderConfig)
	}

	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getZapLevel(level string) zapcore.Level {
	switch level {
	case InfoLevel:
		return zapcore.InfoLevel
	case WarnLevel:
		return zapcore.WarnLevel
	case DebugLevel:
		return zapcore.DebugLevel
	case ErrorLevel:
		return zapcore.ErrorLevel
	case FatalLevel:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func newZapLogger(config Configuration) (*zap.SugaredLogger, error) {
	var cores []zapcore.Core

	if config.EnableConsole {
		level := getZapLevel(config.ConsoleLevel)
		writer := zapcore.Lock(os.Stdout)
		core := zapcore.NewCore(getEncoder(config.ConsoleJSONFormat), writer, level)
		cores = append(cores, core)
		fmt.Println("starting logger with ")
	}

	if config.EnableFile {
		level := getZapLevel(config.FileLevel)

		loggerFile, err := os.OpenFile(config.FileLocation, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			loggerFile, err = os.Create(config.FileLocation)
			if err != nil {
				return nil, fmt.Errorf("failed to open logger storage file: %s error: %w", config.FileLocation, err)
			}
		}

		writer := zapcore.AddSync(loggerFile)
		core := zapcore.NewCore(getEncoder(config.FileJSONFormat), writer, level)
		cores = append(cores, core)
	}

	combinedCore := zapcore.NewTee(cores...)

	// AddCallerSkip skips 1 number of callers, this is important else the file that gets
	// logged will always be the wrapped file. In our case zap.go
	const callersToSkip = 1
	logger := zap.New(combinedCore,
		zap.AddCallerSkip(callersToSkip),
		zap.AddCaller(),
	)

	defer logger.Sync() //nolint

	return logger.Sugar(), nil
}
