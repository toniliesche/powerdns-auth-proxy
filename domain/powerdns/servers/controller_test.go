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

package servers_test

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"powerdns-auth-proxy/domain/powerdns/servers"
	"powerdns-auth-proxy/domain/shared/setup"
	"powerdns-auth-proxy/domain/test"
	"testing"
)

func TestConfigureRoutes(t *testing.T) {
	router := gin.New()
	controller, err := getController()
	if !assert.NoError(t, err, "could not initialize TestConfigureRoutes") {
		return
	}
	controller.ConfigureGroupRoutes(router.Group("/api"))

	if !assert.Equal(t, 2, len(router.Routes()), "There should be 2 route configured") {
		return
	}
}

func TestCallListServers(t *testing.T) {
	controller, err := getController()
	if !assert.NoError(t, err, fmt.Sprintf("could not initialize TestCallListServers: %s", err)) {
		return
	}

	test.RunRequest(t, controller, "/api", "/api/v1/servers", "GET", 200)
}

func TestCallGetServer(t *testing.T) {
	controller, err := getController()
	if !assert.NoError(t, err, fmt.Sprintf("could not initialize TestCallGetServer: %s", err)) {
		return
	}

	test.RunRequest(t, controller, "/api", "/api/v1/servers/localhost", "GET", 200)
}

func getController() (*servers.ServersController, error) {
	container, err := setup.InitContainerTest(&setup.TestConfig{EnableMockAuthentication: true, EnableMockForwardService: true})
	if err != nil {
		return nil, err
	}

	return servers.NewServersController(container)
}
