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

package auth_test

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"powerdns-auth-proxy/domain/shared/auth"
	"powerdns-auth-proxy/domain/shared/auth/model"
	"powerdns-auth-proxy/domain/shared/auth/services"
	"powerdns-auth-proxy/domain/shared/basics"
	dbmodel "powerdns-auth-proxy/domain/shared/database/model"
	"powerdns-auth-proxy/domain/shared/interfaces"
	"powerdns-auth-proxy/domain/shared/setup"
	"testing"
)

func TestAuthentication(t *testing.T) {
	service, err := getAuthenticationService()
	if !assert.NoError(t, err, fmt.Sprintf("could not initialize TestAuthentication: %s", err)) {
		return
	}

	context, _ := gin.CreateTestContext(&httptest.ResponseRecorder{})
	context.Request = httptest.NewRequest("GET", "/test", nil)
	context.Request.SetBasicAuth("user", "testpassword")

	authFunction := service.Authentication(false)
	authFunction(context)

	user, found := context.Get("user")

	if !assert.True(t, found, "user should be found in context") {
		return
	}

	if !assert.NotNil(t, user, "user should not be nil") {
		return
	}
}

func TestAuthenticationFails(t *testing.T) {
	service, err := getAuthenticationService()
	if !assert.NoError(t, err, fmt.Sprintf("could not initialize TestAuthenticationFails: %s", err)) {
		return
	}

	context, _ := gin.CreateTestContext(&httptest.ResponseRecorder{})
	context.Request = httptest.NewRequest("GET", "/test", nil)

	authFunction := service.Authentication(false)
	authFunction(context)

	user, found := context.Get("user")

	if !assert.False(t, found, "user should not be found in context") {
		return
	}

	if !assert.Nil(t, user, "user should be nil") {
		return
	}
}

func TestCheckAccessOnResource(t *testing.T) {
	user := &model.User{
		Username: "user",
		Roles: []string{
			"admin",
		},
	}

	service, err := getAuthenticationService()
	if !assert.NoError(t, err, fmt.Sprintf("could not initialize TestCheckAccessOnResource: %s", err)) {
		return
	}

	context, _ := gin.CreateTestContext(&httptest.ResponseRecorder{})
	context.Request = httptest.NewRequest("GET", "/test", nil)
	context.Request.SetBasicAuth("user", "testpassword")
	context.Set("user", user)

	if !assert.True(t, service.CheckAccessOnResource(context, auth.Admin, ""), "user should have admin role") {
		return
	}
}

func TestCheckAccessOnResourceFailsDueToInvalidRuleSet(t *testing.T) {
	user := &model.User{
		Username: "user",
		Roles: []string{
			"admin",
		},
	}

	service, err := getAuthenticationService()
	if !assert.NoError(t, err, fmt.Sprintf("could not initialize TestCheckAccessOnResourceFailsDueToInvalidRuleSet: %s", err)) {
		return
	}

	context, _ := gin.CreateTestContext(&httptest.ResponseRecorder{})
	context.Request = httptest.NewRequest("GET", "/test", nil)
	context.Request.SetBasicAuth("user", "testpassword")
	context.Set("user", user)

	if !assert.False(t, service.CheckAccessOnResource(context, "invalid", ""), "user should not have admin role") {
		return
	}
}

func TestCheckAccessOnResourceFailsDueToInsufficientPermissions(t *testing.T) {
	user := &model.User{
		Username: "user",
		Roles: []string{
			"user",
		},
	}

	service, err := getAuthenticationService()
	if !assert.NoError(t, err, fmt.Sprintf("could not initialize TestCheckAccessOnResourceFailsDueToInsufficientPermissions: %s", err)) {
		return
	}

	context, _ := gin.CreateTestContext(&httptest.ResponseRecorder{})
	context.Request = httptest.NewRequest("GET", "/test", nil)
	context.Request.SetBasicAuth("user", "testpassword")
	context.Set("user", user)

	if !assert.False(t, service.CheckAccessOnResource(context, auth.Admin, ""), "user should not have admin role") {
		return
	}
}

func TestCheckAccessOnResourceFailsDueToAnonymousAccess(t *testing.T) {
	service, err := getAuthenticationService()
	if !assert.NoError(t, err, fmt.Sprintf("could not initialize TestCheckAccessOnResourceFailsDueToAnonymousAccess: %s", err)) {
		return
	}

	context, _ := gin.CreateTestContext(&httptest.ResponseRecorder{})
	context.Request = httptest.NewRequest("GET", "/test", nil)

	if !assert.False(t, service.CheckAccessOnResource(context, auth.Admin, ""), "anonymous user should not have admin role") {
		return
	}
}

func getAuthenticationService() (interfaces.AuthenticationServiceInterface, error) {
	container, err := setup.InitContainerTest(&setup.TestConfig{EnableMockRepositories: true})
	if err != nil {
		return nil, err
	}

	user := &dbmodel.User{
		Username: "user",
		Password: "$2a$10$0yI64XAHi8q2SVZ.yuGpYeEu2Ufsyz0VgjvBRy5TknrQ9j4glrRQq",
		ApiKeys: []*dbmodel.ApiKey{
			{
				ApiKey: "testapikey",
			},
		},
		UserRoles: []*dbmodel.UserRole{
			{
				Role: &dbmodel.Role{
					Name: "admin",
				},
			},
		},
	}

	if container.UserRepository == nil {
		return nil, basics.NewMissingDependencyError(fmt.Sprintf("user repository could not be provided: %s", err))
	}
	container.UserRepository.SaveNewUser(user)

	container.Authenticator, err = services.NewBasicAuthenticator(container)
	if err != nil {
		return nil, basics.NewMissingDependencyError(fmt.Sprintf("authenticator could not be provided: %s", err))
	}

	container.AuthService, err = auth.NewAuthenticationService(container)
	if err != nil {
		return nil, basics.NewMissingDependencyError(fmt.Sprintf("authentication service could not be provided: %s", err))
	}

	return container.AuthService, nil
}
