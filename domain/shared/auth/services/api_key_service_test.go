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
	"powerdns-auth-proxy/domain/shared/interfaces"
	"powerdns-auth-proxy/domain/shared/setup"
	"powerdns-auth-proxy/domain/test"
	"testing"
)

func TestCreateKey(t *testing.T) {
	registry := test.NewRegistry()
	apiKeyService, err := getApiKeyService(registry, false)
	if !assert.NoError(t, err, "failed to setup api key service") {
		return
	}

	apiKey, err := apiKeyService.Create(test.UserUsername)
	if !assert.NoError(t, err, "failed to create api key") {
		return
	}
	if !assert.NotNil(t, apiKey, "api key is nil") {
		return
	}
}

func TestDeleteKey(t *testing.T) {
	registry := test.NewRegistry()
	apiKeyService, err := getApiKeyService(registry, true)
	if !assert.NoError(t, err, "failed to setup api key service") {
		return
	}

	err = apiKeyService.Delete(registry.Get("apiKey"))
	if !assert.NoError(t, err, "failed to delete api key") {
		return
	}
}

func TestDeleteKeyByIdentifier(t *testing.T) {
	registry := test.NewRegistry()
	apiKeyService, err := getApiKeyService(registry, true)
	if !assert.NoError(t, err, "failed to setup api key service") {
		return
	}

	err = apiKeyService.DeleteByIdentifier(registry.Get("apiKeyIdentifier"))
	if !assert.NoError(t, err, "failed to delete api key by identifier") {
		return
	}
}

func TestListKeys(t *testing.T) {
	registry := test.NewRegistry()
	apiKeyService, err := getApiKeyService(registry, true)
	if !assert.NoError(t, err, "failed to setup api key service") {
		return
	}

	keys, err := apiKeyService.List(test.UserUsername)
	if !assert.NoError(t, err, "failed to list api keys") {
		return
	}
	if !assert.NotNil(t, keys, "keys are nil") {
		return
	}
	if !assert.Len(t, keys, 1, "keys length is not 1") {
		return
	}
}

func getApiKeyService(registry *test.Registry, withData bool) (interfaces.ApiKeyServiceInterface, error) {
	container, err := setup.InitContainerTest(&setup.TestConfig{EnableMockRepositories: true, AuthenticationType: "api_key"})
	if err != nil {
		return nil, err
	}

	err = test.RepositoryCreateTestUser(container, registry)
	if err != nil {
		return nil, err
	}

	if withData {
		err = test.ServiceCreateApiKey(container, registry)
		if err != nil {
			return nil, err
		}
	}

	return container.ApiKeyService, nil
}
