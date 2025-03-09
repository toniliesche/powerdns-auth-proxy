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

func TestAPIKeyCanBeSaved(t *testing.T) {
	var err error
	registry := test.NewRegistry()
	repo := ProvideTestMockAPIKeyRepository(registry, t, false)

	if !assert.NotNil(t, repo, "failed to provide test api key repository") {
		return
	}

	err = repo.SaveNewAPIKey(&model.APIKey{UserID: registry.GetUint("userId"), APIKey: test.APIKey})

	if !assert.NoError(t, err, "failed to save new api key") {
		return
	}
}

func TestAPIKeyCanBeFetched(t *testing.T) {
	var err error
	registry := test.NewRegistry()
	repo := ProvideTestMockAPIKeyRepository(registry, t, true)

	if !assert.NotNil(t, repo, "failed to provide test api key repository") {
		return
	}

	apiKey, err := repo.FetchAPIKey(test.APIKey)

	if !assert.NoError(t, err, "failed to fetch api key") {
		return
	}

	if !assert.IsType(t, &model.APIKey{}, apiKey) {
		return
	}

	if !assert.Equal(t, test.APIKey, apiKey.APIKey) {
		return
	}
}

func TestCorrectAPIKeyWillBeFetched(t *testing.T) {
	var err error
	registry := test.NewRegistry()
	repo := ProvideTestMockAPIKeyRepository(registry, t, true)

	if !assert.NotNil(t, repo, "failed to provide test api key repository") {
		return
	}

	apiKey, err := repo.FetchAPIKey(test.NonExistingAPIKey)

	if !assert.Error(t, err, "expected error to be returned") {
		return
	}

	if !assert.Nil(t, apiKey, "expected api key to be nil") {
		return
	}
}

func TestAPIKeyCanBeDeleted(t *testing.T) {
	var err error
	registry := test.NewRegistry()
	repo := ProvideTestMockAPIKeyRepository(registry, t, true)

	if !assert.NotNil(t, repo, "failed to provide test api key repository") {
		return
	}

	err = repo.DeleteAPIKey(test.APIKey)

	if !assert.NoError(t, err, "failed to delete api key") {
		return
	}

	found := repo.CheckExistenceByAPIKey(test.APIKey)

	if !assert.False(t, found, "expected api key to not exist") {
		return
	}
}

func ProvideTestMockAPIKeyRepository(registry *test.Registry, t *testing.T, withData bool) interfaces.APIKeyRepositoryInterface {
	container, err := setup.InitContainerTest(&setup.TestConfig{EnableMockRepositories: true})
	if !assert.NoError(t, err, "failed to setup test database") {
		return nil
	}

	err = test.RepositoryCreateTestUser(container, registry)
	if !assert.NoError(t, err, "failed to create test user") {
		return nil
	}

	if withData {
		err = test.RepositoryCreateTestAPIKey(container, registry)
		if !assert.NoError(t, err, "failed to create test api key") {
			return nil
		}
	}

	return container.APIKeyRepository
}
