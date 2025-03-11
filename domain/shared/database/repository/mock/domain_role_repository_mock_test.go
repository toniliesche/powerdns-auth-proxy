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

func TestDomainRoleCanBeSaved(t *testing.T) {
	var err error
	registry := test.NewRegistry()
	repo := getDomainRoleRepository(registry, t, false)

	if !assert.NotNil(t, repo, "failed to setup test database") {
		return
	}

	err = repo.SaveNewDomainRole(&model.DomainRole{Name: test.DomainRoleName})

	if !assert.NoError(t, err, "failed to save new domainRole") {
		return
	}
}

func TestDomainRoleExistenceCanBeChecked(t *testing.T) {
	registry := test.NewRegistry()
	repo := getDomainRoleRepository(registry, t, true)

	if !assert.NotNil(t, repo, "failed to setup test database") {
		return
	}

	found := repo.CheckExistenceByName(test.DomainRoleName)

	if !assert.True(t, found, "expected domainRole to exist") {
		return
	}
}

func TestDomainRoleNonExistenceCanBeChecked(t *testing.T) {
	registry := test.NewRegistry()
	repo := getDomainRoleRepository(registry, t, true)

	if !assert.NotNil(t, repo, "failed to setup test database") {
		return
	}

	found := repo.CheckExistenceByName(test.NonExistingDomainRoleName)

	if !assert.False(t, found, "expected domainRole to not exist") {
		return
	}
}

func TestDomainRoleCanBeFetchedByName(t *testing.T) {
	var err error
	registry := test.NewRegistry()
	repo := getDomainRoleRepository(registry, t, true)

	if !assert.NotNil(t, repo, "failed to setup test database") {
		return
	}

	domainRole, err := repo.FetchDomainRoleByName(test.DomainRoleName)

	if !assert.NoError(t, err, "failed to fetch domainRole") {
		return
	}

	if !assert.IsType(t, &model.DomainRole{}, domainRole, "expected domainRole to be of type model.DomainRole") {
		return
	}
}

func TestCorrectDomainRoleWillBeFetched(t *testing.T) {
	var err error
	registry := test.NewRegistry()
	repo := getDomainRoleRepository(registry, t, true)

	if !assert.NotNil(t, repo, "failed to setup test database") {
		return
	}

	domainRole, err := repo.FetchDomainRoleByName(test.NonExistingDomainRoleName)

	if !assert.Error(t, err, "expected error to be returned") {
		return
	}

	if !assert.Nil(t, domainRole, "expected domainRole to be nil") {
		return
	}
}

func TestDomainRoleCanBeDeleted(t *testing.T) {
	var err error
	registry := test.NewRegistry()
	repo := getDomainRoleRepository(registry, t, true)

	if !assert.NotNil(t, repo, "failed to setup test database") {
		return
	}

	err = repo.DeleteDomainRoleByName(test.DomainRoleName)

	if !assert.NoError(t, err, "failed to delete domainRole") {
		return
	}

	found := repo.CheckExistenceByName(test.DomainRoleName)

	if !assert.False(t, found, "expected domainRole to not exist") {
		return
	}
}

func getDomainRoleRepository(registry *test.Registry, t *testing.T, withData bool) interfaces.DomainRoleRepositoryInterface {
	container, err := setup.InitContainerTest(&setup.TestConfig{EnableMockRepositories: true})
	if !assert.NoError(t, err, "failed to setup test database") {
		return nil
	}

	if withData {
		err = test.RepositoryCreateTestDomainRole(container, registry)
		if !assert.NoError(t, err, "failed to setup test database") {
			return nil
		}
	}

	return container.DomainRoleRepository
}
