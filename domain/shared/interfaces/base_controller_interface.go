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

package interfaces

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type BaseControllerInterface interface {
	CheckAccessOnResource(context *gin.Context, ruleSetName string, resource string) bool
	ConfigureGroupRoutes(router *gin.RouterGroup)
	ConfigureEngineRoutes(engine *gin.Engine)
	ParseAndCheck(context *gin.Context, ruleSet string, resource string, target PayloadInterface) bool
	GetIntParameter(context *gin.Context, name string) (uint, bool)
	GetStringParameter(context *gin.Context, name string) (string, bool)
	ParsePayload(context *gin.Context, target PayloadInterface) bool
	CheckAccess(context *gin.Context, ruleSet string, resource string) bool
	HandleSuccess(context *gin.Context)
	HandleData(context *gin.Context, data interface{})
	HandleError(context *gin.Context, err error)
	WriteResponse(context *gin.Context, response *http.Response)
	WriteResponseModified(context *gin.Context, response *http.Response)
	ForbiddenError(context *gin.Context)
	ForwardRequest(request *http.Request) (*http.Response, error)
}
