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
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	nethttp "net/http"
	"powerdns-auth-proxy/domain/shared/basics"
	"powerdns-auth-proxy/domain/shared/http/errors"
	"powerdns-auth-proxy/domain/shared/model"
)

type ResponseWriterAdminAPI struct {
	MessageMapper       *MessageMapper
	StatusCodeMapper    *StatusCodeMapper
	StatusMessageMapper *StatusMessageMapper
	debug               bool
}

func (w *ResponseWriterAdminAPI) HandleSuccess(context *gin.Context) {
	w.WriteResponse(context, w.createResponse(context, nil, nil))
}

func (w *ResponseWriterAdminAPI) HandleData(context *gin.Context, data interface{}) {
	w.WriteResponse(context, w.createResponse(context, data, nil))
}

func (w *ResponseWriterAdminAPI) HandleError(context *gin.Context, err error) {
	w.WriteResponse(context, w.createResponse(context, nil, err))
}

func (w *ResponseWriterAdminAPI) WriteResponse(context *gin.Context, response *nethttp.Response) {
	if response.Body != nil {
		defer response.Body.Close()
	}

	// Copy headers
	for k, v := range response.Header {
		for _, vv := range v {
			context.Header(k, vv)
		}
	}

	// Set the status code
	context.Status(response.StatusCode)

	if response.Body == nil {
		return
	}

	_, err := io.Copy(context.Writer, response.Body)
	if err != nil {
		w.HandleError(context, err)
	}
}

func (w *ResponseWriterAdminAPI) WriteResponseModified(context *gin.Context, response *nethttp.Response) {
	if response.Body != nil {
		defer response.Body.Close()
	}

	// Copy headers
	for k, v := range response.Header {
		for _, vv := range v {
			if k == "Content-Length" || k == "Transfer-Encoding" {
				continue
			}

			context.Header(k, vv)
		}
	}

	// Set the status code
	context.Status(response.StatusCode)

	if response.Body == nil {
		return
	}

	_, err := io.Copy(context.Writer, response.Body)
	if err != nil {
		w.HandleError(context, err)
	}
}

func (w *ResponseWriterAdminAPI) createResponse(context *gin.Context, data interface{}, err error) *nethttp.Response {
	response := model.ResponseData{
		Message:   w.StatusMessageMapper.MapMessage(err),
		RequestID: context.GetString("request_id"),
		Data:      data,
	}

	if w.debug {
		response.AdvancedMessage = w.MessageMapper.MapAdvancedMessage(err)
	}

	jsonResponse, _ := json.Marshal(response)
	body := io.NopCloser(
		bytes.NewBuffer(jsonResponse),
	)

	return &nethttp.Response{
		StatusCode: w.StatusCodeMapper.MapStatusCode(err),
		Request:    context.Request,
		Body:       body,
		Header: map[string][]string{
			"Content-Type": {
				"application/json",
			},
		},
	}
}

func (w *ResponseWriterAdminAPI) ForbiddenError(context *gin.Context) {
	w.HandleError(context, errors.NewForbiddenError(fmt.Errorf("you are not authorized to access this resource")))
}

func (w *ResponseWriterAdminAPI) SetDebug(debug bool) {
	w.debug = debug
}

func NewResponseWriterAdminAPI(container *basics.InjectionContainer) (*ResponseWriterAdminAPI, error) {
	if container == nil {
		return nil, basics.NewMissingDependencyError("could not provide admin api response writer: passed injection container is nil")
	}

	if container.Config == nil {
		return nil, basics.NewMissingDependencyError("could not provide admin api response writer: config could not be resolved")
	}

	return &ResponseWriterAdminAPI{
		MessageMapper:       &MessageMapper{},
		StatusCodeMapper:    &StatusCodeMapper{},
		StatusMessageMapper: &StatusMessageMapper{},
		debug:               container.Config.Debug,
	}, nil
}
