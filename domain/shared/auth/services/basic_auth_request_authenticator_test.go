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

package services_test

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"powerdns-auth-proxy/domain/shared/auth/services"
	"testing"
)

func TestBasicAuthSucceeds(t *testing.T) {
	container := getContainer("basic_auth")
	authenticator, _ := services.ProvideBasicAuthenticator(container)

	context, _ := gin.CreateTestContext(&httptest.ResponseRecorder{})
	context.Request = httptest.NewRequest("GET", "/test", nil)
	context.Request.SetBasicAuth("user", "testpassword")

	user, err := authenticator.Authenticate(context)

	if !assert.NoError(t, err, "error should be nil") {
		return
	}

	if !assert.NotNil(t, user, "user should not be nil") {
		return
	}

	if !assert.Equal(t, "user", user.Username, "username should be user") {
		return
	}
}

func TestBasicAuthFailsOnMissingHeader(t *testing.T) {
	container := getContainer("basic_auth")
	authenticator, _ := services.ProvideBasicAuthenticator(container)

	context, _ := gin.CreateTestContext(&httptest.ResponseRecorder{})
	context.Request = httptest.NewRequest("GET", "/test", nil)

	user, err := authenticator.Authenticate(context)

	if !assert.Error(t, err, "error should not be nil") {
		return
	}

	if !assert.Nil(t, user, "user should be nil") {
		return
	}
}

func TestBasicAuthFailsOnWrongUsername(t *testing.T) {
	container := getContainer("basic_auth")
	authenticator, _ := services.ProvideBasicAuthenticator(container)

	context, _ := gin.CreateTestContext(&httptest.ResponseRecorder{})
	context.Request = httptest.NewRequest("GET", "/test", nil)
	context.Request.SetBasicAuth("wronguser", "testpassword")

	user, err := authenticator.Authenticate(context)

	if !assert.Error(t, err, "error should not be nil") {
		return
	}

	if !assert.Nil(t, user, "user should be nil") {
		return
	}
}

func TestBasicAuthFailsOnWrongPassword(t *testing.T) {
	container := getContainer("basic_auth")
	authenticator, _ := services.ProvideBasicAuthenticator(container)

	context, _ := gin.CreateTestContext(&httptest.ResponseRecorder{})
	context.Request = httptest.NewRequest("GET", "/test", nil)
	context.Request.SetBasicAuth("user", "wrongpassword")

	user, err := authenticator.Authenticate(context)

	if !assert.Error(t, err, "error should not be nil") {
		return
	}

	if !assert.Nil(t, user, "user should be nil") {
		return
	}
}
