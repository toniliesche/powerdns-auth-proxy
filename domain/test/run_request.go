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

package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	basehttp "net/http"
	"net/http/httptest"
	"powerdns-auth-proxy/domain/shared/interfaces"
	"testing"
)

func RunRequest(t *testing.T, controller interfaces.ControllerInterface, prefix string, path string, method string, code int, payload ...interface{}) {
	router := gin.New()
	controller.ConfigureEngineRoutes(router)
	controller.ConfigureGroupRoutes(router.Group(prefix))

	var body []byte
	if len(payload) > 0 {
		body, _ = json.Marshal(payload[0])
	}

	w := httptest.NewRecorder()
	req, _ := basehttp.NewRequest(method, path, bytes.NewBuffer(body))
	router.ServeHTTP(w, req)

	if !assert.Equal(t, code, w.Code, fmt.Sprintf("Response should be %d", code)) {
		fmt.Println(w.Body.String())

		return
	}

	responseBody := w.Body.Bytes()
	response := &Response{}

	err := json.Unmarshal(responseBody, response)

	if !assert.NoError(t, err, "Failed to unmarshal response") {
		return
	}

	if response.Path == "" || response.StatusCode == 0 || response.Method == "" {
		return
	}

	if !assert.Equal(t, code, response.StatusCode, fmt.Sprintf("Response status should be %d", code)) {
		return
	}

	if !assert.Equal(t, method, response.Method, fmt.Sprintf("Method should be %s", method)) {
		return
	}

	if !assert.Equal(t, path, response.Path, fmt.Sprintf("Path should be %s", path)) {
		return
	}
}
