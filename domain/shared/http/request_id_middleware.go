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
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"powerdns-auth-proxy/domain/shared/basics"
)

type RequestIDMiddleware struct {
	app    string
	logger *zerolog.Logger
}

func (m *RequestIDMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("request_id", m.app+"-"+uuid.New().String())
		c.Next()
	}
}

func NewRequestIDMiddleware(container *basics.InjectionContainer) (*RequestIDMiddleware, error) {
	if container == nil {
		return nil, basics.NewMissingDependencyError("could not provide request id middleware: passed injection container is nil")
	}

	return &RequestIDMiddleware{
		app:    "powerdns-auth-proxy",
		logger: container.Logger,
	}, nil
}
