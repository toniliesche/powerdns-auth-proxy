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

package cache_test

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"powerdns-auth-proxy/domain/powerdns/cache"
	"powerdns-auth-proxy/domain/shared/setup"
	"powerdns-auth-proxy/domain/test"
	"testing"
)

func TestConfigureRoutes(t *testing.T) {
	router := gin.New()
	controller, err := getController()
	if !assert.NoError(t, err, fmt.Sprintf("could not initialize TestConfigureRoutes: %s", err)) {
		return
	}

	controller.ConfigureGroupRoutes(router.Group("/api"))

	if !assert.Equal(t, 1, len(router.Routes()), "There should be 1 route configured") {
		return
	}
}

func TestCallFlushCache(t *testing.T) {
	controller, err := getController()
	if !assert.NoError(t, err, fmt.Sprintf("could not initialize TestCallFlushCache: %s", err)) {
		return
	}

	test.RunRequest(t, controller, "/api", "/api/v1/servers/localhost/cache/flush", "PUT", 200)
}

func getController() (*cache.CacheController, error) {
	container, err := setup.InitContainerTest(&setup.TestConfig{EnableMockAuthentication: true, EnableMockForwardService: true})
	if err != nil {
		return nil, err
	}

	return container.CacheController.(*cache.CacheController), nil
}
