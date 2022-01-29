package logger

import (
	"fmt"
	"os"

	"github.com/alexandr-lakeev/otus-final-project/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	logg *zap.Logger
}

func New(cfg config.LoggerConf) (*Logger, error) {
	logg, err := initLogger(cfg)
	if err != nil {
		return nil, err
	}

	return &Logger{
		logg: logg,
	}, nil
}

func (l Logger) Debug(msg string) {
	l.logg.Debug(msg)
}

func (l Logger) Info(msg string) {
	l.logg.Info(msg)
}

func (l Logger) Warning(msg string) {
	l.logg.Warn(msg)
}

func (l Logger) Error(msg string) {
	l.logg.Error(msg)
}

func (l Logger) Panic(msg string) {
	l.logg.Panic(msg)
}

func initLogger(cfg config.LoggerConf) (*zap.Logger, error) {
	levelMap := map[string]zapcore.Level{
		"DEBUG":   zap.DebugLevel,
		"INFO":    zap.InfoLevel,
		"WARNING": zap.WarnLevel,
		"ERROR":   zap.ErrorLevel,
		"PANIC":   zap.PanicLevel,
	}

	atom := zap.NewAtomicLevel()

	var logger *zap.Logger

	if cfg.Env == "prod" {
		logger = zap.New(zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			zapcore.Lock(os.Stdout),
			atom,
		))
	} else {
		logger = zap.New(zapcore.NewCore(
			zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
			zapcore.Lock(os.Stdout),
			atom,
		))
	}

	zapLevel, ok := levelMap[cfg.Level]
	if !ok {
		return nil, fmt.Errorf("wrong log level: %s", cfg.Level)
	}

	atom.SetLevel(zapLevel)

	return logger, nil
}
