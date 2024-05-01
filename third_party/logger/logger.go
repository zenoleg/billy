package logger

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
)

func init() {
	zerolog.ErrorStackMarshaler = marshalStack
}

func NewLogger(cfg LogConfig, appVersion string) zerolog.Logger {
	var logger zerolog.Logger
	if cfg.Fmt() == "console" {
		logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout})
	} else {
		logger = zerolog.New(os.Stdout)
	}

	logger = logger.With().
		Str("ver", appVersion).
		Timestamp().
		Logger()

	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil || level == zerolog.NoLevel {
		logger.Error().Err(err).Msg("failed to parse logger.level")
		logger.Level(zerolog.InfoLevel)
	} else {
		logger.Level(level)
	}

	return logger
}

type LogConfig struct {
	Level  string
	Format string
}

func NewConfig() LogConfig {
	level, exists := os.LookupEnv("MIGRATOR_LOGGER_LEVEL")
	if !exists {
		level = "info"
	}

	format, exists := os.LookupEnv("MIGRATOR_LOGGER_FORMAT")
	if !exists {
		format = "console"
	}

	return LogConfig{
		Level:  level,
		Format: format,
	}
}

func (c *LogConfig) Fmt() string {
	result, err := FormatEnumString(c.Format)
	if err != nil {
		fmt.Fprintln(os.Stderr, "incorrect logger.format string", err)
		result = Console
	}

	return result.String()
}
