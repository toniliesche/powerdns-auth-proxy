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

package rdbms_test

import (
	"github.com/stretchr/testify/assert"
	"powerdns-auth-proxy/domain/shared/database/model"
	"powerdns-auth-proxy/domain/shared/database/repository/interfaces"
	"powerdns-auth-proxy/domain/shared/setup"
	"powerdns-auth-proxy/domain/test"
	"testing"
)

func TestTokenSessionCanBeSaved(t *testing.T) {
	var err error
	registry := test.NewRegistry()
	repo := ProvideTestTokenSessionRepository(registry, t, false)

	if !assert.NotNil(t, repo, "failed to provide test token session repository") {
		return
	}

	err = repo.SaveNewTokenSession(&model.TokenSession{UserId: registry.GetUint("userId"), SessionId: test.TokenSession})

	if !assert.NoError(t, err, "failed to save new token session") {
		return
	}
}

func TestTokenSessionCanBeFetchedBySessionId(t *testing.T) {
	registry := test.NewRegistry()
	repo := ProvideTestTokenSessionRepository(registry, t, true)

	if !assert.NotNil(t, repo, "failed to provide test token session repository") {
		return
	}

	tokenSession, err := repo.FetchTokenSession(test.TokenSession)

	if !assert.NoError(t, err, "failed to get token session by session ID") {
		return
	}

	if !assert.NotNil(t, tokenSession, "expected token session to not be nil") {
		return
	}
}

func TestCorrectTokenSessionCanBeFetched(t *testing.T) {
	registry := test.NewRegistry()
	repo := ProvideTestTokenSessionRepository(registry, t, true)

	if !assert.NotNil(t, repo, "failed to provide test token session repository") {
		return
	}

	tokenSession, err := repo.FetchTokenSession(test.NonExistingTokenSession)

	if !assert.Error(t, err, "expected error when fetching non-existing token session") {
		return
	}

	if !assert.Nil(t, tokenSession, "expected token session to be nil") {
		return
	}
}

func TestTokenSessionCanBeDeleted(t *testing.T) {
	var err error
	registry := test.NewRegistry()
	repo := ProvideTestTokenSessionRepository(registry, t, true)

	if !assert.NotNil(t, repo, "failed to provide test token session repository") {
		return
	}

	err = repo.DeleteTokenSession(test.TokenSession)

	if !assert.NoError(t, err, "failed to delete token session") {
		return
	}
}

func ProvideTestTokenSessionRepository(registry *test.Registry, t *testing.T, init bool) interfaces.TokenSessionRepositoryInterface {
	container, err := setup.InitContainerTest(&setup.TestConfig{RunDatabaseMigrations: true})
	if !assert.NoError(t, err, "failed to setup test database") {
		return nil
	}

	err = test.RepositoryCreateTestUser(container, registry)
	if init {
		err = test.RepositoryCreateTestTokenSession(container, registry)
		if !assert.NoError(t, err, "failed to create test token session") {
			return nil
		}
	}

	return container.SessionRepository
}
