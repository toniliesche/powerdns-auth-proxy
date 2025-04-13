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
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"net/http/httputil"
	"powerdns-auth-proxy/domain/shared/basics"
)

type RequestContentLogMiddleware struct {
	logger *zerolog.Logger
}

func (m *RequestContentLogMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		dump, err := httputil.DumpRequest(c.Request, true)
		if err == nil {
			m.logger.Trace().
				Msg("Incoming request")
			m.logger.Trace().
				Msg(string(dump))
		}

		c.Next()
	}
}

func NewRequestContentLogMiddleware(container *basics.InjectionContainer) (*RequestContentLogMiddleware, error) {
	if container == nil {
		return nil, basics.NewMissingDependencyError("could not provide request content log middleware: passed injection container is nil")
	}

	return &RequestContentLogMiddleware{
		logger: container.Logger,
	}, nil
}
