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

package index_test

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	basehttp "net/http"
	"net/http/httptest"
	"powerdns-auth-proxy/domain/admin/index"
	"powerdns-auth-proxy/domain/shared/basics"
	"powerdns-auth-proxy/domain/shared/controller"
	"powerdns-auth-proxy/domain/shared/model"
	"powerdns-auth-proxy/domain/shared/setup"
	"testing"
)

func TestCallNotFound(t *testing.T) {
	router := gin.New()
	indexController, err := getController()
	if !assert.NoError(t, err, fmt.Sprintf("could not initialize TestCallNotFound: %s", err)) {
		return
	}

	indexController.ConfigureEngineRoutes(router)

	w := httptest.NewRecorder()
	req, _ := basehttp.NewRequest("GET", "/api/v1/servers/localhost/cache/flush", nil)
	router.ServeHTTP(w, req)

	if !assert.Equal(t, 404, w.Code, "Response should be 404") {
		return
	}

	body := w.Body.Bytes()
	response := &model.ResponseData{}

	err = json.Unmarshal(body, response)

	if !assert.NoError(t, err, "Failed to unmarshal response") {
		return
	}

	if !assert.Equal(t, "not found", response.Message, "Message should be 'not found'") {
		return
	}

	if !assert.Equal(t, "path '/api/v1/servers/localhost/cache/flush' not found", response.AdvancedMessage, "Error should be 'path '/api/v1/servers/localhost/cache/flush' not found'") {
		return
	}
}

func TestProvideIndexControllerFailsOnMissingBaseController(t *testing.T) {
	container := &basics.InjectionContainer{}

	indexController, err := index.ProvideIndexController(container)
	assert.Error(t, err, "provide index controller method should return an error")
	assert.Equal(t, "index controller could not be created: base controller could not be resolved", err.Error())
	assert.Nil(t, indexController, "provide index controller method should not return a controller")
}

func TestProvideIndexControllerSucceeds(t *testing.T) {
	container := &basics.InjectionContainer{}
	container.BaseController = &controller.BaseController{}

	indexController, err := index.ProvideIndexController(container)
	assert.NoError(t, err, "provide index controller method should succeed")
	assert.NotNil(t, indexController, "provide index controller method should return a controller")
}

func getController() (*index.IndexController, error) {
	container, err := setup.InitContainerTest(&setup.TestConfig{})
	if err != nil {
		return nil, err
	}

	container.ResponseWriter.SetDebug(true)

	return container.AdminIndexController.(*index.IndexController), nil
}
