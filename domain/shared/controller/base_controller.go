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

package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"powerdns-auth-proxy/domain/shared/basics"
	httperrors "powerdns-auth-proxy/domain/shared/http/errors"
	"powerdns-auth-proxy/domain/shared/interfaces"
	"strconv"
)

type BaseController struct {
	interfaces.ResponseWriterInterface
	AuthenticationService interfaces.AuthenticationServiceInterface
	ForwardService        interfaces.ForwardServiceInterface
}

func (c *BaseController) CheckAccessOnResource(context *gin.Context, ruleSetName string, resource string) bool {
	return c.AuthenticationService.CheckAccessOnResource(context, ruleSetName, resource)
}

func (c *BaseController) ForwardRequest(request *http.Request) (*http.Response, error) {
	return c.ForwardService.ForwardRequest(request)
}

func (c *BaseController) ConfigureGroupRoutes(router *gin.RouterGroup) {
}

func (c *BaseController) ConfigureEngineRoutes(engine *gin.Engine) {
}

func (c *BaseController) ParseAndCheck(context *gin.Context, ruleSet string, resource string, target interfaces.PayloadInterface) bool {
	if !c.CheckAccess(context, ruleSet, resource) {
		return false
	}

	if !c.ParsePayload(context, target) {
		return false
	}

	return true
}

func (c *BaseController) GetIntParameter(context *gin.Context, name string) (uint, bool) {
	value := context.Param(name)
	uintValue, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		c.HandleError(context, httperrors.NewBadRequestError(fmt.Errorf("invalid %s", name)))
		return 0, false
	}

	return uint(uintValue), true
}

func (c *BaseController) GetStringParameter(context *gin.Context, name string) (string, bool) {
	value := context.Param(name)
	if value == "" {
		c.HandleError(context, httperrors.NewBadRequestError(fmt.Errorf("invalid %s", name)))
		return "", false
	}

	return value, true
}

func (c *BaseController) ParsePayload(context *gin.Context, target interfaces.PayloadInterface) bool {
	if err := context.ShouldBindJSON(target); err != nil {
		c.HandleError(context, httperrors.NewBadRequestError(fmt.Errorf("failed parsing payload")))
		return false
	}

	if err := target.Verify(); err != nil {
		c.HandleError(context, httperrors.NewBadRequestError(err))
		return false
	}

	return true
}

func (c *BaseController) CheckAccess(context *gin.Context, ruleSet string, resource string) bool {
	if !c.AuthenticationService.CheckAccessOnResource(context, ruleSet, resource) {
		c.HandleError(context, httperrors.NewForbiddenError(fmt.Errorf("insufficient permissions to access resource")))
		return false
	}

	return true
}

func ProvideControllerAdminAPI(container *basics.InjectionContainer) (*BaseController, error) {
	if container.ResponseWriterAdminAPI == nil {
		return nil, fmt.Errorf("could not provide base controller: response writer could not be resolved")
	}

	if container.AuthService == nil {
		return nil, fmt.Errorf("could not provide base controller: auth service could not be resolved")
	}

	if container.ForwardService == nil {
		return nil, fmt.Errorf("could not provide base controller: forward service could not be resolved")
	}

	return &BaseController{container.ResponseWriterAdminAPI, container.AuthService, container.ForwardService}, nil
}

func ProvideControllerPowerDNS(container *basics.InjectionContainer) (*BaseController, error) {
	if container.ResponseWriterPowerDNS == nil {
		return nil, fmt.Errorf("could not provide base controller: response writer could not be resolved")
	}

	if container.AuthService == nil {
		return nil, fmt.Errorf("could not provide base controller: auth service could not be resolved")
	}

	if container.ForwardService == nil {
		return nil, fmt.Errorf("could not provide base controller: forward service could not be resolved")
	}

	return &BaseController{container.ResponseWriterPowerDNS, container.AuthService, container.ForwardService}, nil
}
