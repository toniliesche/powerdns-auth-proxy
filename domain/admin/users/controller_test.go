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

package users_test

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"powerdns-auth-proxy/domain/admin/users"
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

	if !assert.Equal(t, 5, len(router.Routes()), "There should be 5 routes configured") {
		return
	}
}

func TestCallListUsersEndpoint(t *testing.T) {
	registry := test.NewRegistry()
	controller, err := getController(registry, false)
	if !assert.NoError(t, err, fmt.Sprintf("could not initialize TestCallListUsersEndpoint: %s", err)) {
		return
	}

	test.RunRequest(t, controller, "/admin", "/admin/v1/users", "GET", 200)
}

func TestCallCreateUserEndpoint(t *testing.T) {
	registry := test.NewRegistry()
	controller, err := getController(registry, false)
	if !assert.NoError(t, err, fmt.Sprintf("could not initialize TestCallCreateUserEndpoint: %s", err)) {
		return
	}

	payload := &management.UserCreatePayload{
		Username: test.UserUsername,
		Password: test.UserPassword,
	}

	test.RunRequest(t, controller, "/admin", "/admin/v1/users", "POST", 200, payload)
}

func TestCallGetUserEndpoint(t *testing.T) {
	registry := test.NewRegistry()
	controller, err := getController(registry, true)
	if !assert.NoError(t, err, fmt.Sprintf("could not initialize TestCallGetUserEndpoint: %s", err)) {
		return
	}

	test.RunRequest(t, controller, "/admin", fmt.Sprintf("/admin/v1/users/details/%d", registry.GetUint("userId")), "GET", 200)
}

func TestCallDeleteUserEndpoint(t *testing.T) {
	registry := test.NewRegistry()
	controller, err := getController(registry, true)
	if !assert.NoError(t, err, fmt.Sprintf("could not initialize TestCallDeleteUserEndpoint: %s", err)) {
		return
	}

	test.RunRequest(t, controller, "/admin", fmt.Sprintf("/admin/v1/users/%d", registry.GetUint("userId")), "DELETE", 200)
}

func TestCallUpdateUserEndpoint(t *testing.T) {
	registry := test.NewRegistry()
	controller, err := getController(registry, true)
	if !assert.NoError(t, err, fmt.Sprintf("could not initialize TestCallUpdateUserEndpoint: %s", err)) {
		return
	}

	payload := &management.UserUpdatePayload{
		Password: test.UserPassword,
	}

	test.RunRequest(t, controller, "/admin", fmt.Sprintf("/admin/v1/users/%d", registry.GetUint("userId")), "PATCH", 200, payload)
}

func getController(registry *test.Registry, withData bool) (*users.UserController, error) {
	container, err := setup.InitContainerTest(&setup.TestConfig{EnableMockAuthentication: true, EnableMockForwardService: true, EnableMockRepositories: true})
	if err != nil {
		return nil, err
	}

	if withData {
		err = test.RepositoryCreateTestUser(container, registry)
		if err != nil {
			return nil, err
		}
	}

	container.ResponseWriterAdminAPI.SetDebug(true)

	return container.AdminUserController.(*users.UserController), nil
}
