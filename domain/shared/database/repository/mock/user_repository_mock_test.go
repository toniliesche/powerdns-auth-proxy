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

func TestUserCanBeSaved(t *testing.T) {
	var err error
	registry := test.NewRegistry()
	repo := getUserRepository(registry, t, false)

	if !assert.NotNil(t, repo, "failed to provide test user repository") {
		return
	}

	err = repo.SaveNewUser(&model.User{Username: test.UserUsername, Password: test.UserPassword})

	if !assert.NoError(t, err, "failed to save new user") {
		return
	}
}

func TestUserExistenceCanBeChecked(t *testing.T) {
	registry := test.NewRegistry()
	repo := getUserRepository(registry, t, true)

	if !assert.NotNil(t, repo, "failed to provide test user repository") {
		return
	}

	found := repo.CheckExistenceByUsername(test.UserUsername)

	if !assert.True(t, found, "expected user to exist") {
		return
	}
}

func TestUserNonExistenceCanBeChecked(t *testing.T) {
	registry := test.NewRegistry()
	repo := getUserRepository(registry, t, true)

	if !assert.NotNil(t, repo, "failed to provide test user repository") {
		return
	}

	found := repo.CheckExistenceByUsername(test.NoneExistingUserUsername)

	if !assert.False(t, found, "expected user to not exist") {
		return
	}
}

func TestUserCanBeFetchedByUsername(t *testing.T) {
	var err error
	registry := test.NewRegistry()
	repo := getUserRepository(registry, t, true)

	if !assert.NotNil(t, repo, "failed to provide test user repository") {
		return
	}

	user, err := repo.FetchUser(test.UserUsername)
	if !assert.NoError(t, err, "failed to fetch user") {
		return
	}

	if !assert.IsType(t, &model.User{}, user, "expected user to be of type model.User") {
		return
	}

	if !assert.Equal(t, test.UserUsername, user.Username, "expected user to have correct username") {
		return
	}
}

func TestCorrectMockUserWillBeFetched(t *testing.T) {
	var err error
	registry := test.NewRegistry()
	repo := getUserRepository(registry, t, true)

	if !assert.NotNil(t, repo, "failed to provide test user repository") {
		return
	}

	user, err := repo.FetchUser(test.NoneExistingUserUsername)

	if !assert.Error(t, err, "expected error to be returned") {
		return
	}

	if !assert.Nil(t, user, "expected user to be nil") {
		return
	}
}

func TestUserCanBeFetchedByApiKey(t *testing.T) {
	var err error
	registry := test.NewRegistry()
	repo := getUserRepository(registry, t, true)

	if !assert.NotNil(t, repo, "failed to provide test user repository") {
		return
	}

	user, err := repo.FetchUserByApiKey(test.ApiKey)

	if !assert.NoError(t, err, "failed to fetch user") {
		return
	}

	if !assert.IsType(t, &model.User{}, user, "expected user to be of type model.User") {
		return
	}

	if !assert.Equal(t, test.UserUsername, user.Username, "expected user to have correct username") {
		return
	}
}

func TestCorrectMockUserWillBeFetchedByApiKey(t *testing.T) {
	var err error
	registry := test.NewRegistry()
	repo := getUserRepository(registry, t, true)

	if !assert.NotNil(t, repo, "failed to provide test user repository") {
		return
	}

	user, err := repo.FetchUserByApiKey(test.NonExistingApiKey)

	if !assert.Error(t, err, "expected error to be returned") {
		return
	}

	if !assert.Nil(t, user, "expected user to be nil") {
		return
	}
}

func TestUserCanBeDeleted(t *testing.T) {
	var err error
	registry := test.NewRegistry()
	repo := getUserRepository(registry, t, true)

	if !assert.NotNil(t, repo, "failed to provide test user repository") {
		return
	}

	err = repo.DeleteUser(test.UserUsername)

	if !assert.NoError(t, err, "failed to delete user") {
		return
	}

	found := repo.CheckExistenceByUsername(test.UserUsername)

	if !assert.False(t, found, "expected user to not exist") {
		return
	}
}

func getUserRepository(registry *test.Registry, t *testing.T, init bool) interfaces.UserRepositoryInterface {
	container, err := setup.InitContainerTest(&setup.TestConfig{EnableMockRepositories: true})
	if !assert.NoError(t, err, "failed to setup test database") {
		return nil
	}

	if init {
		err = test.RepositoryCreateTestUser(container, registry)
		if !assert.NoError(t, err, "failed to create test user") {
			return nil
		}

		err = test.RepositoryCreateTestApiKey(container, registry)
		if !assert.NoError(t, err, "failed to create test api key") {
			return nil
		}

		apiKey, err := container.ApiKeyRepository.FetchApiKey(test.ApiKey)
		if !assert.NoError(t, err, "failed to fetch api key") {
			return nil
		}

		user, err := container.UserRepository.FetchUser(test.UserUsername)
		if !assert.NoError(t, err, "failed to fetch user") {
			return nil
		}

		user.ApiKeys = append(user.ApiKeys, apiKey)
	}

	return container.UserRepository
}
