package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"

	"github.com/wylu1037/go-micro-boilerplate/pkg/config"
)

var defaultLogger zerolog.Logger

func Init(cfg config.LogConfig) error {
	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}

	var output zerolog.LevelWriter
	if cfg.Format == "console" {
		output = zerolog.LevelWriterAdapter{Writer: zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}}
	} else {
		output = zerolog.LevelWriterAdapter{Writer: os.Stdout}
	}

	defaultLogger = zerolog.New(output).
		Level(level).
		With().
		Timestamp().
		Caller().
		Logger()

	return nil
}

func Get() *zerolog.Logger {
	return &defaultLogger
}

func With() zerolog.Context {
	return defaultLogger.With()
}

func Debug() *zerolog.Event {
	return defaultLogger.Debug()
}

func Info() *zerolog.Event {
	return defaultLogger.Info()
}

func Warn() *zerolog.Event {
	return defaultLogger.Warn()
}

func Error() *zerolog.Event {
	return defaultLogger.Error()
}

func Fatal() *zerolog.Event {
	return defaultLogger.Fatal()
}
