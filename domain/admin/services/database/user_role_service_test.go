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

func TestListRolesForUser(t *testing.T) {
	registry := test.NewRegistry()
	service, err := getUserRoleService(registry, true)
	if err != nil {
		t.Error(err)
		return
	}

	roles, err := service.ListRolesForUser(test.UserUsername)
	if !assert.NoError(t, err, "failed to list roles for user") {
		return
	}

	assert.NotEmpty(t, roles, "expected roles to not be empty")
	assert.Equal(t, 1, len(roles), "expected one role")
	assert.Equal(t, test.RoleName, roles[0].Role, "expected role to match role name")
}

func TestListRolesForUserByUserId(t *testing.T) {
	registry := test.NewRegistry()
	service, err := getUserRoleService(registry, true)
	if err != nil {
		t.Error(err)
		return
	}

	roles, err := service.ListRolesForUserByUserId(registry.GetUint("userId"))
	if !assert.NoError(t, err, "failed to list roles for user") {
		return
	}

	assert.NotEmpty(t, roles, "expected roles to not be empty")
	assert.Equal(t, 1, len(roles), "expected one role")
	assert.Equal(t, test.RoleName, roles[0].Role, "expected role to match role name")
}

func TestGrantRoleToUser(t *testing.T) {
	registry := test.NewRegistry()
	service, err := getUserRoleService(registry, false)
	if err != nil {
		t.Error(err)
		return
	}

	err = service.GrantRoleToUser(test.UserUsername, test.RoleName)
	if !assert.NoError(t, err, "failed to grant role to user") {
		return
	}

	roles, err := service.ListRolesForUser(test.UserUsername)
	if !assert.NoError(t, err, "failed to list roles for user") {
		return
	}

	assert.NotEmpty(t, roles, "expected roles to not be empty")
	assert.Equal(t, 1, len(roles), "expected one role")
	assert.Equal(t, test.RoleName, roles[0].Role, "expected role to match role name")
}

func TestGrantRolesToUser(t *testing.T) {
	registry := test.NewRegistry()
	service, err := getUserRoleService(registry, false)
	if err != nil {
		t.Error(err)
		return
	}

	err = service.GrantRolesToUser(registry.GetUint("userId"), []string{test.RoleName})
	if !assert.NoError(t, err, "failed to grant roles to user") {
		return
	}

	roles, err := service.ListRolesForUser(test.UserUsername)
	if !assert.NoError(t, err, "failed to list roles for user") {
		return
	}

	assert.NotEmpty(t, roles, "expected roles to not be empty")
	assert.Equal(t, 1, len(roles), "expected one role")
	assert.Equal(t, test.RoleName, roles[0].Role, "expected role to match role name")
}

func TestRevokeRoleFromUser(t *testing.T) {
	registry := test.NewRegistry()
	service, err := getUserRoleService(registry, true)
	if err != nil {
		t.Error(err)
		return
	}

	err = service.RevokeRoleFromUser(test.UserUsername, test.RoleName)
	if !assert.NoError(t, err, "failed to revoke role from user") {
		return
	}

	roles, err := service.ListRolesForUser(test.UserUsername)
	if !assert.NoError(t, err, "failed to list roles for user") {
		return
	}

	assert.Empty(t, roles, "expected roles to be empty")
}

func TestRevokeRolesFromUser(t *testing.T) {
	registry := test.NewRegistry()
	service, err := getUserRoleService(registry, true)
	if err != nil {
		t.Error(err)
		return
	}

	err = service.RevokeRolesFromUser(registry.GetUint("userId"), []string{test.RoleName})
	if !assert.NoError(t, err, "failed to revoke roles from user") {
		return
	}

	roles, err := service.ListRolesForUser(test.UserUsername)
	if !assert.NoError(t, err, "failed to list roles for user") {
		return
	}

	assert.Empty(t, roles, "expected roles to be empty")
}

func TestProvideUserRoleServiceFailsOnMissingRoleRepository(t *testing.T) {
	container := &basics.InjectionContainer{}

	service, err := database.NewUserRoleService(container)
	assert.Error(t, err, "provide user role service should return an error")
	assert.Equal(t, "could not provide user role service: role repository could not be resolved", err.Error())
	assert.Nil(t, service, "provide user role service should return nil")
}

func TestProvideUserRoleServiceFailsOnMissingUserRepository(t *testing.T) {
	container := &basics.InjectionContainer{}
	container.RoleRepository = &rdbms.RoleRepository{}

	service, err := database.NewUserRoleService(container)
	assert.Error(t, err, "provide user role service should return an error")
	assert.Equal(t, "could not provide user role service: user repository could not be resolved", err.Error())
	assert.Nil(t, service, "provide user role service should return nil")
}

func TestProvideUserRoleServiceFailsOnMissingUserRoleRepository(t *testing.T) {
	container := &basics.InjectionContainer{}
	container.RoleRepository = &rdbms.RoleRepository{}
	container.UserRepository = &rdbms.UserRepository{}

	service, err := database.NewUserRoleService(container)
	assert.Error(t, err, "provide user role service should return an error")
	assert.Equal(t, "could not provide user role service: user role repository could not be resolved", err.Error())
	assert.Nil(t, service, "provide user role service should return nil")
}

func TestProvideUserRoleServiceSucceeds(t *testing.T) {
	container := &basics.InjectionContainer{}
	container.RoleRepository = &rdbms.RoleRepository{}
	container.UserRepository = &rdbms.UserRepository{}
	container.UserRoleRepository = &rdbms.UserRoleRepository{}

	service, err := database.NewUserRoleService(container)
	assert.NoError(t, err, "provide user role service should not return an error")
	assert.NotNil(t, service, "provide user role service should return a service")
}

func getUserRoleService(registry *test.Registry, withData bool) (*database.UserRoleService, error) {
	container, err := setup.InitContainerTest(&setup.TestConfig{EnableMockRepositories: true})
	if err != nil {
		return nil, err
	}

	service, err := database.NewUserRoleService(container)
	if err != nil {
		return nil, err
	}

	err = test.RepositoryCreateTestUser(container, registry)
	if err != nil {
		return nil, err
	}

	err = test.RepositoryCreateTestRole(container, registry)
	if err != nil {
		return nil, err
	}

	if withData {
		err = test.RepositoryCreateTestUserRole(container, registry)
		if err != nil {
			return nil, err
		}
	}

	return service, nil
}
