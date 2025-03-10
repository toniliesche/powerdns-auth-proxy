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

package sessions_test

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"powerdns-auth-proxy/domain/admin/sessions"
	"powerdns-auth-proxy/domain/shared/model/management"
	"powerdns-auth-proxy/domain/shared/setup"
	"powerdns-auth-proxy/domain/test"
	"testing"
)

func TestConfigureRoutes(t *testing.T) {
	router := gin.New()
	registry := test.NewRegistry()
	controller, err := getController(registry, false)
	if !assert.NoError(t, err, fmt.Sprintf("could not initialize TestConfigureRoutes: %s", err)) {
		return
	}

	controller.ConfigureGroupRoutes(router.Group("/api"))

	if !assert.Equal(t, 2, len(router.Routes()), "There should be 2 routes configured") {
		return
	}
}

func TestCallListSessionsEndpoint(t *testing.T) {
	registry := test.NewRegistry()
	controller, err := getController(registry, false)
	if !assert.NoError(t, err, fmt.Sprintf("could not initialize TestCallListSessionsEndpoint: %s", err)) {
		return
	}

	test.RunRequest(t, controller, "/admin", "/admin/v1/users/1/sessions", "GET", 200)
}

func TestCallLogoutSessionEndpoint(t *testing.T) {
	registry := test.NewRegistry()
	controller, err := getController(registry, true)
	if !assert.NoError(t, err, fmt.Sprintf("could not initialize TestCallLogoutSessionEndpoint: %s", err)) {
		return
	}

	test.RunRequest(t, controller, "/admin", "/admin/v1/users/1/sessions/logout", "POST", 200, &management.SessionLogoutPayload{SessionIDs: []string{registry.Get("tokenSession")}})
}

func getController(registry *test.Registry, withData bool) (*sessions.SessionsController, error) {
	container, err := setup.InitContainerTest(&setup.TestConfig{EnableMockAuthentication: true, EnableMockForwardService: true, EnableMockRepositories: true})
	if err != nil {
		return nil, err
	}

	if withData {
		err = test.RepositoryCreateTestUser(container, registry)
		if err != nil {
			return nil, err
		}

		err = test.RepositoryCreateTestTokenSession(container, registry)
		if err != nil {
			return nil, err
		}
	}

	container.ResponseWriter.SetDebug(true)

	return container.AdminSessionsController.(*sessions.SessionsController), nil
}
