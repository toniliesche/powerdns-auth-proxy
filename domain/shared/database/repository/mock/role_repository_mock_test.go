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

func TestRoleCanBeSaved(t *testing.T) {
	var err error
	registry := test.NewRegistry()
	repo := getRoleRepository(registry, t, false)

	if !assert.NotNil(t, repo, "failed to provide test role repository") {
		return
	}

	err = repo.SaveNewRole(&model.Role{Name: test.RoleName})

	if !assert.NoError(t, err, "failed to save new role") {
		return
	}
}

func TestRoleExistenceCanBeChecked(t *testing.T) {
	registry := test.NewRegistry()
	repo := getRoleRepository(registry, t, true)

	if !assert.NotNil(t, repo, "failed to provide test role repository") {
		return
	}

	found := repo.CheckExistenceByName(test.RoleName)

	if !assert.True(t, found, "expected role to exist") {
		return
	}
}

func TestRoleNonExistenceCanBeChecked(t *testing.T) {
	registry := test.NewRegistry()
	repo := getRoleRepository(registry, t, true)

	if !assert.NotNil(t, repo, "failed to provide test role repository") {
		return
	}

	found := repo.CheckExistenceByName(test.NonExistingRoleName)

	if !assert.False(t, found, "expected role to not exist") {
		return
	}
}

func TestRoleCanBeFetchedByName(t *testing.T) {
	var err error
	registry := test.NewRegistry()
	repo := getRoleRepository(registry, t, true)

	if !assert.NotNil(t, repo, "failed to provide test role repository") {
		return
	}

	role, err := repo.FetchRoleByName(test.RoleName)

	if !assert.NoError(t, err, "failed to fetch role") {
		return
	}

	if !assert.IsType(t, &model.Role{}, role, "expected role to be of type model.Role") {
		return
	}
}

func TestCorrectRoleWillBeFetched(t *testing.T) {
	var err error
	registry := test.NewRegistry()
	repo := getRoleRepository(registry, t, true)

	if !assert.NotNil(t, repo, "failed to provide test role repository") {
		return
	}

	role, err := repo.FetchRoleByName(test.NonExistingRoleName)

	if !assert.Error(t, err, "expected error to be returned") {
		return
	}

	if !assert.Nil(t, role, "expected role to be nil") {
		return
	}
}

func TestRoleCanBeDeleted(t *testing.T) {
	var err error
	registry := test.NewRegistry()
	repo := getRoleRepository(registry, t, true)

	if !assert.NotNil(t, repo, "failed to provide test role repository") {
		return
	}

	err = repo.DeleteRoleByName(test.RoleName)

	if !assert.NoError(t, err, "failed to delete role") {
		return
	}

	found := repo.CheckExistenceByName(test.RoleName)

	if !assert.False(t, found, "expected role to not exist") {
		return
	}
}

func getRoleRepository(registry *test.Registry, t *testing.T, withData bool) interfaces.RoleRepositoryInterface {
	container, err := setup.InitContainerTest(&setup.TestConfig{EnableMockRepositories: true})
	if !assert.NoError(t, err, "failed to setup test database") {
		return nil
	}

	if withData {
		err = test.RepositoryCreateTestRole(container, registry)
		if !assert.NoError(t, err, "failed to setup test database") {
			return nil
		}
	}

	return container.RoleRepository
}
