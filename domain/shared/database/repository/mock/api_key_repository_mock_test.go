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

package mock_test

import (
	"github.com/stretchr/testify/assert"
	"powerdns-auth-proxy/domain/shared/database/model"
	"powerdns-auth-proxy/domain/shared/database/repository/interfaces"
	"powerdns-auth-proxy/domain/shared/setup"
	"powerdns-auth-proxy/domain/test"
	"testing"
)

func TestApiKeyCanBeSaved(t *testing.T) {
	var err error
	registry := test.NewRegistry()
	repo := ProvideTestMockApiKeyRepository(registry, t, false)

	if !assert.NotNil(t, repo, "failed to provide test api key repository") {
		return
	}

	err = repo.SaveNewApiKey(&model.ApiKey{UserId: registry.GetUint("userId"), ApiKey: test.ApiKey})

	if !assert.NoError(t, err, "failed to save new api key") {
		return
	}
}

func TestApiKeyCanBeFetched(t *testing.T) {
	var err error
	registry := test.NewRegistry()
	repo := ProvideTestMockApiKeyRepository(registry, t, true)

	if !assert.NotNil(t, repo, "failed to provide test api key repository") {
		return
	}

	apiKey, err := repo.FetchApiKey(test.ApiKey)

	if !assert.NoError(t, err, "failed to fetch api key") {
		return
	}

	if !assert.IsType(t, &model.ApiKey{}, apiKey) {
		return
	}

	if !assert.Equal(t, test.ApiKey, apiKey.ApiKey) {
		return
	}
}

func TestApiKeyCanBeFetchedByIdentifier(t *testing.T) {
	var err error
	registry := test.NewRegistry()
	repo := ProvideTestMockApiKeyRepository(registry, t, true)

	if !assert.NotNil(t, repo, "failed to provide test api key repository") {
		return
	}

	apiKey, err := repo.FetchApiKeyByIdentifier(test.ApiKeyIdentifier)

	if !assert.NoError(t, err, "failed to fetch api key by identifier") {
		return
	}

	if !assert.IsType(t, &model.ApiKey{}, apiKey) {
		return
	}

	if !assert.Equal(t, test.ApiKeyIdentifier, apiKey.Identifier) {
		return
	}
}

func TestCorrectApiKeyWillBeFetched(t *testing.T) {
	var err error
	registry := test.NewRegistry()
	repo := ProvideTestMockApiKeyRepository(registry, t, true)

	if !assert.NotNil(t, repo, "failed to provide test api key repository") {
		return
	}

	apiKey, err := repo.FetchApiKey(test.NonExistingApiKey)

	if !assert.Error(t, err, "expected error to be returned") {
		return
	}

	if !assert.Nil(t, apiKey, "expected api key to be nil") {
		return
	}
}

func TestApiKeyCanBeDeleted(t *testing.T) {
	var err error
	registry := test.NewRegistry()
	repo := ProvideTestMockApiKeyRepository(registry, t, true)

	if !assert.NotNil(t, repo, "failed to provide test api key repository") {
		return
	}

	err = repo.DeleteApiKey(test.ApiKey)

	if !assert.NoError(t, err, "failed to delete api key") {
		return
	}

	found := repo.CheckExistenceByApiKey(test.ApiKey)

	if !assert.False(t, found, "expected api key to not exist") {
		return
	}
}

func ProvideTestMockApiKeyRepository(registry *test.Registry, t *testing.T, withData bool) interfaces.ApiKeyRepositoryInterface {
	container, err := setup.InitContainerTest(&setup.TestConfig{EnableMockRepositories: true})
	if !assert.NoError(t, err, "failed to setup test database") {
		return nil
	}

	err = test.RepositoryCreateTestUser(container, registry)
	if !assert.NoError(t, err, "failed to create test user") {
		return nil
	}

	if withData {
		err = test.RepositoryCreateTestApiKey(container, registry)
		if !assert.NoError(t, err, "failed to create test api key") {
			return nil
		}
	}

	return container.ApiKeyRepository
}
