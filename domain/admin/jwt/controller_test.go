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

package jwt_test

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"powerdns-auth-proxy/domain/admin/jwt"
	"powerdns-auth-proxy/domain/shared/auth/model"
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

	controller.ConfigureEngineRoutes(router)

	if !assert.Equal(t, 2, len(router.Routes()), "There should be 5 route configured") {
		return
	}
}

func TestCallLogin(t *testing.T) {
	controller, err := getController()
	if !assert.NoError(t, err, fmt.Sprintf("could not initialize TestCallLogin: %s", err)) {
		return
	}

	test.RunRequest(t, controller, "/api", "/auth/login", "POST", 200, &model.LoginPayload{Username: "user", Password: "password"})
}

func TestCallRefresh(t *testing.T) {
	controller, err := getController()
	if !assert.NoError(t, err, fmt.Sprintf("could not initialize TestCallRefresh: %s", err)) {
		return
	}

	test.RunRequest(t, controller, "/api", "/auth/refresh", "POST", 200, &model.RefreshPayload{RefreshToken: "token"})
}

func getController() (*jwt.JWTController, error) {
	container, err := setup.InitContainerTest(&setup.TestConfig{EnableMockJwtService: true, EnableMockAuthentication: true, EnableMockAuthorization: true, EnableMockRepositories: true})
	if err != nil {
		return nil, err
	}

	return jwt.ProvideJwtController(container)
}
