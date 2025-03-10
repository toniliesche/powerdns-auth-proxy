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

func TestUserDomainRoleCanBeSaved(t *testing.T) {
	var err error
	registry := test.NewRegistry()
	repo := ProvideTestUserDomainRoleRepository(registry, t, false)

	if !assert.NotNil(t, repo, "failed to provide test user domain role repository") {
		return
	}

	err = repo.SaveNewUserDomainRole(&model.UserDomainRole{UserID: registry.GetUint("userId"), DomainID: registry.GetUint("domainId"), Domain: &model.Domain{Fqdn: test.DomainFqdn}, DomainRoleID: registry.GetUint("domainRoleId")})

	if !assert.NoError(t, err, "failed to save new user domain role") {
		return
	}
}

func TestUserDomainRoleCanBeDeleted(t *testing.T) {
	var err error
	registry := test.NewRegistry()
	repo := ProvideTestUserDomainRoleRepository(registry, t, true)

	if !assert.NotNil(t, repo, "failed to provide test user domain role repository") {
		return
	}

	err = repo.DeleteUserDomainRole(&model.UserDomainRole{UserID: registry.GetUint("userId"), DomainID: registry.GetUint("domainId"), Domain: &model.Domain{Fqdn: test.DomainFqdn}, DomainRoleID: registry.GetUint("domainRoleId")})

	if !assert.NoError(t, err, "failed to delete user domain role") {
		return
	}
}

func ProvideTestUserDomainRoleRepository(registry *test.Registry, t *testing.T, init bool) interfaces.UserDomainRoleRepositoryInterface {
	container, err := setup.InitContainerTest(&setup.TestConfig{RunDatabaseMigrations: true})
	if !assert.NoError(t, err, "failed to setup test database") {
		return nil
	}

	err = test.RepositoryCreateTestDomain(container, registry)
	if !assert.NoError(t, err, "failed to create test domain") {
		return nil
	}

	err = test.RepositoryCreateTestDomainRole(container, registry)
	if !assert.NoError(t, err, "failed to create test domain role") {
		return nil
	}

	err = test.RepositoryCreateTestUser(container, registry)
	if !assert.NoError(t, err, "failed to create test user") {
		return nil
	}

	if init {
		err = test.RepositoryCreateTestUserDomainRole(container, registry)
		if !assert.NoError(t, err, "failed to create test user domain role") {
			return nil
		}
	}

	return container.UserDomainRoleRepository
}
