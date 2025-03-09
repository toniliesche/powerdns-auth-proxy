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
	"powerdns-auth-proxy/domain/shared/model/management"
	"powerdns-auth-proxy/domain/shared/setup"
	"powerdns-auth-proxy/domain/test"
	"testing"
)

func TestCreateRole(t *testing.T) {
	registry := test.NewRegistry()
	service, err := ProvideTestRoleService(registry, false)
	if err != nil {
		t.Error(err)
		return
	}

	payload := &management.RoleCreatePayload{
		Name: test.RoleName,
	}

	_, err = service.CreateRole(payload)
	if !assert.NoError(t, err, "failed to create role") {
		return
	}
}

func TestListRoles(t *testing.T) {
	registry := test.NewRegistry()
	service, err := ProvideTestRoleService(registry, true)
	if err != nil {
		t.Error(err)
		return
	}

	roles, err := service.ListRoles()
	if err != nil {
		t.Error(err)
		return
	}

	if !assert.NotEmpty(t, roles, "expected roles to not be empty") {
		return
	}
}

func TestDeleteRole(t *testing.T) {
	registry := test.NewRegistry()
	service, err := ProvideTestRoleService(registry, true)
	if err != nil {
		t.Error(err)
		return
	}

	err = service.DeleteRole(test.RoleName)
	if !assert.NoError(t, err, "failed to delete role") {
		return
	}
}

func ProvideTestRoleService(registry *test.Registry, withData bool) (*database.RoleService, error) {
	container, err := setup.InitContainerTest(&setup.TestConfig{EnableMockRepositories: true})
	if err != nil {
		return nil, err
	}

	if withData {
		err := test.RepositoryCreateTestRole(container, registry)
		if err != nil {
			return nil, err
		}
	}

	return database.ProvideRoleService(container)
}
