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
	"github.com/stretchr/testify/assert"
	"powerdns-auth-proxy/domain/shared/database/model"
	"powerdns-auth-proxy/domain/shared/interfaces"
	"testing"
)

func TestCreateKey(t *testing.T) {
	apiKeyService, err := setupApiKeyService()
	if !assert.NoError(t, err, "failed to setup api key service") {
		return
	}

	apiKey, err := apiKeyService.Create("test-user")
	if !assert.NoError(t, err, "failed to create api key") {
		return
	}

	if !assert.NotNil(t, apiKey, "api key is nil") {
		return
	}
}

func setupApiKeyService() (interfaces.ApiKeyServiceInterface, error) {
	container := getContainer("api_key")

	dbUser := &model.User{
		Username: "test-user",
	}

	if err := container.UserRepository.SaveNewUser(dbUser); err != nil {
		return nil, err
	}

	return container.ApiKeyService, nil
}
