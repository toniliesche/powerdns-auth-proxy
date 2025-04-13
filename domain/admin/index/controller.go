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

package index

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"powerdns-auth-proxy/domain/shared/basics"
	"powerdns-auth-proxy/domain/shared/http/errors"
	"powerdns-auth-proxy/domain/shared/interfaces"
)

type IndexController struct {
	interfaces.BaseControllerInterface
}

func (c *IndexController) ConfigureEngineRoutes(router *gin.Engine) {
	router.NoRoute(c.NotFound)
}

func (c *IndexController) NotFound(context *gin.Context) {
	c.HandleError(context, errors.NewNotFoundError(fmt.Errorf("path '%s' not found", context.Request.URL.Path)))
}

func NewIndexController(container *basics.InjectionContainer) (*IndexController, error) {
	if container == nil {
		return nil, basics.NewMissingDependencyError("could not provide index controller: passed injection container is nil")
	}

	if container.BaseControllerAdminAPI == nil {
		return nil, basics.NewMissingDependencyError("could not provide index controller: base controller could not be resolved")
	}

	return &IndexController{container.BaseControllerAdminAPI}, nil
}
