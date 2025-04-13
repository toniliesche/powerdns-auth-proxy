// MIT License
// Copyright (c) 2025 Toni Liesche
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

package log

import (
	"fmt"
	"github.com/rs/zerolog"
	"io"
	"os"
	"powerdns-auth-proxy/domain/shared/basics"
	"time"
)

func NewTempLogger() *zerolog.Logger {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	return &logger
}

func NewLogger(container *basics.InjectionContainer) (*zerolog.Logger, error) {
	if container == nil {
		return nil, basics.NewMissingDependencyError("could not provide logger: passed injection container is nil")
	}

	if container.Config == nil {
		return nil, basics.NewMissingDependencyError("could not provide logger: system config could not be resolved")
	}

	cfg := container.Config

	var err error
	var writer io.Writer

	if cfg.Debug {
		zerolog.TimeFieldFormat = time.RFC3339Nano
		writer = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "2006-01-02 15:04:05.000",
			FormatLevel: func(i interface{}) string {
				if ll, ok := i.(string); ok {
					return "[" + ll + "]"
				}
				return "[???]"
			},
		}
	} else {
		var logFile string
		if cfg.LogPath == "" {
			logFile = "/dev/stdout"
		} else {
			logFile = cfg.LogPath
		}

		var fileErr error
		writer, fileErr = os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if fileErr != nil {
			fmt.Printf("Error opening log file %s: %v\n", logFile, fileErr)
			return nil, err
		}
		zerolog.TimeFieldFormat = "2006-01-02 15:04:05.000"
	}

	logLevel, err := resolveLogLevel(cfg.LogLevel)
	if err != nil {
		return nil, err
	}

	logger := zerolog.
		New(writer).
		Level(zerolog.DebugLevel).
		With().
		Timestamp().
		Str("component.id", "powerdns-auth-proxy").
		Logger()

	logger.Debug().
		Msg("Logger setup complete. Switching to application logger")

	logger.Debug().
		Msgf("Setting log level to %s", logLevel.String())

	logger = logger.Level(logLevel)

	return &logger, nil
}

func resolveLogLevel(level string) (zerolog.Level, error) {
	switch level {
	case "trace":
		return zerolog.TraceLevel, nil
	case "debug":
		return zerolog.DebugLevel, nil
	case "info", "":
		return zerolog.InfoLevel, nil
	case "warn":
		return zerolog.WarnLevel, nil
	case "error", "fatal", "panic":
		return zerolog.ErrorLevel, nil
	default:
		return zerolog.NoLevel, fmt.Errorf("invalid log level: %s", level)
	}
}
