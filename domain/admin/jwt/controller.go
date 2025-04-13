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

package jwt

import (
	"github.com/gin-gonic/gin"
	"powerdns-auth-proxy/domain/shared/basics"
	"powerdns-auth-proxy/domain/shared/interfaces"
)

type JWTController struct {
	interfaces.BaseControllerInterface
	jwtService interfaces.JWTServiceInterface
}

func (c *JWTController) ConfigureEngineRoutes(engine *gin.Engine) {
	engine.POST("/auth/login", c.login)
	engine.POST("/auth/refresh", c.refresh)
}

func NewJwtController(container *basics.InjectionContainer) (*JWTController, error) {
	if container == nil {
		return nil, basics.NewMissingDependencyError("could not provide jwt controller: passed injection container is nil")
	}

	if container.BaseControllerAdminAPI == nil {
		return nil, basics.NewMissingDependencyError("could not provide jwt controller: base controller could not be resolved")
	}

	return &JWTController{container.BaseControllerAdminAPI, container.JWTService}, nil
}
