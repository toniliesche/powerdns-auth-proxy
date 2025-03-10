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

package database_test

import (
	"github.com/stretchr/testify/assert"
	"powerdns-auth-proxy/domain/admin/services/database"
	"powerdns-auth-proxy/domain/shared/basics"
	"powerdns-auth-proxy/domain/shared/database/repository/rdbms"
	"powerdns-auth-proxy/domain/shared/setup"
	"powerdns-auth-proxy/domain/test"
	"testing"
)

func TestListUserSessions(t *testing.T) {
	registry := test.NewRegistry()
	service, err := getSessionService(registry, true)
	if err != nil {
		t.Error(err)
		return
	}

	sessions, err := service.ListUserSessions(registry.GetUint("userId"))
	if !assert.NoError(t, err, "failed to list user sessions") {
		return
	}

	assert.NotEmpty(t, sessions, "expected sessions to not be empty")
	assert.Equal(t, len(sessions), 1, "expected one session")
	assert.Equal(t, sessions[0].SessionID, registry.Get("tokenSession"), "expected session ID to match")
}

func TestLogoutUserSession(t *testing.T) {
	registry := test.NewRegistry()
	service, err := getSessionService(registry, true)
	if err != nil {
		t.Error(err)
		return
	}

	err = service.LogoutUserSession([]string{registry.Get("tokenSession")})
	assert.NoError(t, err, "failed to logout user session")

	sessions, err := service.ListUserSessions(registry.GetUint("userId"))
	if !assert.NoError(t, err, "failed to list user sessions") {
		return
	}

	assert.Empty(t, sessions, "expected sessions to be empty")
}

func TestProvideSessionServiceFailsOnMissingSessionRepository(t *testing.T) {
	container := &basics.InjectionContainer{}

	_, err := database.ProvideSessionService(container)
	assert.Error(t, err, "provide session service should return an error")
	assert.Equal(t, "could not provide session service: session repository could not be resolved", err.Error())
}

func TestProvideSessionServiceFailsOnMissingUserRepository(t *testing.T) {
	container := &basics.InjectionContainer{}
	container.SessionRepository = &rdbms.TokenSessionRepository{}

	_, err := database.ProvideSessionService(container)
	assert.Error(t, err, "provide session service should return an error")
	assert.Equal(t, "could not provide session service: user repository could not be resolved", err.Error())
}

func TestProvideSessionServiceSucceeds(t *testing.T) {
	container := &basics.InjectionContainer{}
	container.SessionRepository = &rdbms.TokenSessionRepository{}
	container.UserRepository = &rdbms.UserRepository{}

	service, err := database.ProvideSessionService(container)
	assert.NoError(t, err, "provide session service should not return an error")
	assert.NotNil(t, service, "provide session service should return a service")
}

func getSessionService(registry *test.Registry, withData bool) (*database.SessionService, error) {
	container, err := setup.InitContainerTest(&setup.TestConfig{EnableMockRepositories: true})
	if err != nil {
		return nil, err
	}

	if withData {
		err := test.RepositoryCreateTestUser(container, registry)
		if err != nil {
			return nil, err
		}

		err = test.RepositoryCreateTestTokenSession(container, registry)
		if err != nil {
			return nil, err
		}
	}

	return container.SessionService.(*database.SessionService), nil
}
