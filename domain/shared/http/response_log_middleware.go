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

package http

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"powerdns-auth-proxy/domain/shared/basics"
)

type ResponseLogMiddleware struct {
	logger *zerolog.Logger
}

func (m *ResponseLogMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		dump := fmt.Sprintf("Status: %d\n", c.Writer.Status())
		dump = dump + "Headers:\n"

		for k, v := range c.Writer.Header() {
			dump = dump + fmt.Sprintf("  %s: %s\n", k, v)
		}

		dump = dump + "Body:\n" + blw.body.String()

		m.logger.Trace().
			Msg("Incoming response")
		m.logger.Trace().
			Msg(dump)
	}
}

func NewResponseLogMiddleware(container *basics.InjectionContainer) (*ResponseLogMiddleware, error) {
	if container == nil {
		return nil, basics.NewMissingDependencyError("could not provide request log middleware: passed injection container is nil")
	}

	return &ResponseLogMiddleware{
		logger: container.Logger,
	}, nil
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
