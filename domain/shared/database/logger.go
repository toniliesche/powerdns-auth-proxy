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

package database

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"gorm.io/gorm/logger"
	"time"
)

type GormLogger struct {
	logger *zerolog.Logger
}

func (l *GormLogger) LogMode(logger.LogLevel) logger.Interface {
	return l
}

func (l *GormLogger) Error(ctx context.Context, msg string, opts ...interface{}) {
	l.logger.Error().Msg(fmt.Sprintf(msg, opts...))
}

func (l *GormLogger) Warn(ctx context.Context, msg string, opts ...interface{}) {
	l.logger.Warn().Msg(fmt.Sprintf(msg, opts...))
}

func (l *GormLogger) Info(ctx context.Context, msg string, opts ...interface{}) {
	l.logger.Info().Msg(fmt.Sprintf(msg, opts...))
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, f func() (string, int64), err error) {
	zl := l.logger
	var event *zerolog.Event

	if err != nil {
		event = zl.Debug()
	} else {
		event = zl.Trace()
	}

	var dur_key string

	switch zerolog.DurationFieldUnit {
	case time.Nanosecond:
		dur_key = "elapsed_ns"
	case time.Microsecond:
		dur_key = "elapsed_us"
	case time.Millisecond:
		dur_key = "elapsed_ms"
	case time.Second:
		dur_key = "elapsed"
	case time.Minute:
		dur_key = "elapsed_min"
	case time.Hour:
		dur_key = "elapsed_hr"
	default:
		zl.Error().Interface("zerolog.DurationFieldUnit", zerolog.DurationFieldUnit).Msg("gormzerolog encountered a mysterious, unknown value for DurationFieldUnit")
		dur_key = "elapsed_"
	}

	event.Dur(dur_key, time.Since(begin))

	sql, rows := f()
	if sql != "" {
		event.Str("sql", sql)
	}
	if rows > -1 {
		event.Int64("rows", rows)
	}

	event.Send()

	return
}

func newGormLogger(logger *zerolog.Logger) logger.Interface {
	return &GormLogger{logger: logger}
}
