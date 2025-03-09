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

package autoprimaries

import (
	"github.com/gin-gonic/gin"
	"powerdns-auth-proxy/domain/shared/basics"
	"powerdns-auth-proxy/domain/shared/interfaces"
)

type AutoprimariesController struct {
	interfaces.BaseControllerInterface
}

func (c *AutoprimariesController) ConfigureGroupRoutes(router *gin.RouterGroup) {
	router.GET("/v1/servers/:server/zones/:zone/autoprimaries", c.listAutoPrimaries)
	router.POST("/v1/servers/:server/zones/:zone/autoprimaries", c.createAutoPrimary)
	router.DELETE("/v1/servers/:server/zones/:zone/autoprimaries/:ip/:nameserver", c.deleteAutoPrimary)
}

func ProvideAutoprimariesController(container *basics.InjectionContainer) (*AutoprimariesController, error) {
	if container.BaseController == nil {
		return nil, basics.NewMissingDependencyError("autoprimaries controller could not be created: base controller could not be resolved")
	}

	return &AutoprimariesController{container.BaseController}, nil
}
