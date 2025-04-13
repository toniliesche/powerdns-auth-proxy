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
	"powerdns-auth-proxy/domain/shared/model/management"
	"powerdns-auth-proxy/domain/shared/setup"
	"powerdns-auth-proxy/domain/test"
	"testing"
)

func TestListDomainRolesForUserByUserId(t *testing.T) {
	registry := test.NewRegistry()
	service, err := getUserDomainRoleService(registry, true)
	if err != nil {
		t.Error(err)
		return
	}

	domainRoles, err := service.ListDomainRolesForUserByUserId(registry.GetUint("userId"))
	if !assert.NoError(t, err, "failed to list domain roles for user") {
		return
	}

	assert.NotEmpty(t, domainRoles, "expected domain roles to not be empty")
	assert.Equal(t, 1, len(domainRoles), "expected one domain role")
	assert.Equal(t, test.DomainFqdn, domainRoles[0].Domain, "expected domain to match fqdn")
	assert.Equal(t, test.DomainRoleName, domainRoles[0].Role, "expected domain role to match role name")
}

func TestListDomainRolesForUserByUsername(t *testing.T) {
	registry := test.NewRegistry()
	service, err := getUserDomainRoleService(registry, true)
	if err != nil {
		t.Error(err)
		return
	}

	domainRoles, err := service.ListDomainRolesForUser(test.UserUsername)
	if !assert.NoError(t, err, "failed to list domain roles for user") {
		return
	}

	assert.NotEmpty(t, domainRoles, "expected domain roles to not be empty")
	assert.Equal(t, 1, len(domainRoles), "expected one domain role")
	assert.Equal(t, test.DomainFqdn, domainRoles[0].Domain, "expected domain to match fqdn")
	assert.Equal(t, test.DomainRoleName, domainRoles[0].Role, "expected domain role to match role name")
}

func TestListDomainRolesForUserByUserIdPerDomain(t *testing.T) {
	registry := test.NewRegistry()
	service, err := getUserDomainRoleService(registry, true)
	if err != nil {
		t.Error(err)
		return
	}

	domainRoles, err := service.ListDomainRolesForUserByUserIdPerDomain(test.DomainFqdn, registry.GetUint("userId"))
	if !assert.NoError(t, err, "failed to list domain roles for user") {
		return
	}

	assert.NotEmpty(t, domainRoles, "expected domain roles to not be empty")
	assert.Equal(t, 1, len(domainRoles), "expected one domain role")
	assert.Equal(t, test.DomainFqdn, domainRoles[0].Domain, "expected domain to match fqdn")
	assert.Equal(t, test.DomainRoleName, domainRoles[0].Role, "expected domain role to match role name")
}

func TestListDomainRolesForUserPerDomain(t *testing.T) {
	registry := test.NewRegistry()
	service, err := getUserDomainRoleService(registry, true)
	if err != nil {
		t.Error(err)
		return
	}

	domainRoles, err := service.ListDomainRolesForUserPerDomain(test.DomainFqdn, test.UserUsername)
	if !assert.NoError(t, err, "failed to list domain roles for user") {
		return
	}

	assert.NotEmpty(t, domainRoles, "expected domain roles to not be empty")
	assert.Equal(t, 1, len(domainRoles), "expected one domain role")
	assert.Equal(t, test.DomainFqdn, domainRoles[0].Domain, "expected domain to match fqdn")
	assert.Equal(t, test.DomainRoleName, domainRoles[0].Role, "expected domain role to match role name")
}

func TestGrantDomainRoleToUser(t *testing.T) {
	registry := test.NewRegistry()
	service, err := getUserDomainRoleService(registry, false)
	if err != nil {
		t.Error(err)
		return
	}

	err = service.GrantDomainRoleToUser(test.DomainFqdn, test.UserUsername, test.DomainRoleName)
	if !assert.NoError(t, err, "failed to grant domain role to user") {
		return
	}

	domainRoles, err := service.ListDomainRolesForUserByUserId(registry.GetUint("userId"))
	if !assert.NoError(t, err, "failed to list domain roles for user") {
		return
	}

	assert.NotEmpty(t, domainRoles, "expected domain roles to not be empty")
	assert.Equal(t, 1, len(domainRoles), "expected one domain role")
	assert.Equal(t, test.DomainFqdn, domainRoles[0].Domain, "expected domain to match fqdn")
	assert.Equal(t, test.DomainRoleName, domainRoles[0].Role, "expected domain role to match role name")
}

func TestGrantDomainRolesToUser(t *testing.T) {
	registry := test.NewRegistry()
	service, err := getUserDomainRoleService(registry, false)
	if err != nil {
		t.Error(err)
		return
	}

	err = service.GrantDomainRolesToUser(registry.GetUint("userId"), []*management.DomainRole{{Domain: test.DomainFqdn, Role: test.DomainRoleName}})
	if !assert.NoError(t, err, "failed to grant domain roles to user") {
		return
	}

	domainRoles, err := service.ListDomainRolesForUserByUserId(registry.GetUint("userId"))
	if !assert.NoError(t, err, "failed to list domain roles for user") {
		return
	}

	assert.NotEmpty(t, domainRoles, "expected domain roles to not be empty")
	assert.Equal(t, 1, len(domainRoles), "expected one domain role")
	assert.Equal(t, test.DomainFqdn, domainRoles[0].Domain, "expected domain to match fqdn")
	assert.Equal(t, test.DomainRoleName, domainRoles[0].Role, "expected domain role to match role name")
}

func TestRevokeDomainRoleFromUser(t *testing.T) {
	registry := test.NewRegistry()
	service, err := getUserDomainRoleService(registry, true)
	if err != nil {
		t.Error(err)
		return
	}

	err = service.RevokeDomainRoleFromUser(test.DomainFqdn, test.UserUsername, test.DomainRoleName)
	if !assert.NoError(t, err, "failed to revoke domain role from user") {
		return
	}

	domainRoles, err := service.ListDomainRolesForUserByUserId(registry.GetUint("userId"))
	if !assert.NoError(t, err, "failed to list domain roles for user") {
		return
	}

	assert.Empty(t, domainRoles, "expected domain roles to be empty")
}

func TestRevokeDomainRolesFromUser(t *testing.T) {
	registry := test.NewRegistry()
	service, err := getUserDomainRoleService(registry, true)
	if err != nil {
		t.Error(err)
		return
	}

	err = service.RevokeDomainRolesFromUser(registry.GetUint("userId"), []*management.DomainRole{{Domain: test.DomainFqdn, Role: test.DomainRoleName}})
	if !assert.NoError(t, err, "failed to revoke domain roles from user") {
		return
	}

	domainRoles, err := service.ListDomainRolesForUserByUserId(registry.GetUint("userId"))
	if !assert.NoError(t, err, "failed to list domain roles for user") {
		return
	}

	assert.Empty(t, domainRoles, "expected domain roles to be empty")
}

func TestProvideUserDomainRoleServiceFailsOnMissingDomainRepository(t *testing.T) {
	container := &basics.InjectionContainer{}

	_, err := database.NewUserDomainRoleService(container)
	assert.Error(t, err, "provide user domain role service should return an error")
	assert.Equal(t, "could not provide user domain role service: domain repository could not be resolved", err.Error())
}

func TestProvideUserDomainRoleServiceFailsOnMissingDomainRoleRepository(t *testing.T) {
	container := &basics.InjectionContainer{}
	container.DomainRepository = &rdbms.DomainRepository{}

	_, err := database.NewUserDomainRoleService(container)
	assert.Error(t, err, "provide user domain role service should return an error")
	assert.Equal(t, "could not provide user domain role service: domain role repository could not be resolved", err.Error())
}

func TestProvideUserDomainRoleServiceFailsOnMissingUserRepository(t *testing.T) {
	container := &basics.InjectionContainer{}
	container.DomainRepository = &rdbms.DomainRepository{}
	container.DomainRoleRepository = &rdbms.DomainRoleRepository{}

	_, err := database.NewUserDomainRoleService(container)
	assert.Error(t, err, "provide user domain role service should return an error")
	assert.Equal(t, "could not provide user domain role service: user repository could not be resolved", err.Error())
}

func TestProvideUserDomainRoleServiceFailsOnMissingUserDomainRoleRepository(t *testing.T) {
	container := &basics.InjectionContainer{}
	container.DomainRepository = &rdbms.DomainRepository{}
	container.DomainRoleRepository = &rdbms.DomainRoleRepository{}
	container.UserRepository = &rdbms.UserRepository{}

	_, err := database.NewUserDomainRoleService(container)
	assert.Error(t, err, "provide user domain role service should return an error")
	assert.Equal(t, "could not provide user domain role service: user domain role repository could not be resolved", err.Error())
}

func TestProvideUserDomainRoleServiceSucceeds(t *testing.T) {
	container := &basics.InjectionContainer{}
	container.DomainRepository = &rdbms.DomainRepository{}
	container.DomainRoleRepository = &rdbms.DomainRoleRepository{}
	container.UserRepository = &rdbms.UserRepository{}
	container.UserDomainRoleRepository = &rdbms.UserDomainRoleRepository{}

	service, err := database.NewUserDomainRoleService(container)
	assert.NoError(t, err, "provide user domain role service should not return an error")
	assert.NotNil(t, service, "provide user domain role service should return a service")
}

func getUserDomainRoleService(registry *test.Registry, withData bool) (*database.UserDomainRoleService, error) {
	container, err := setup.InitContainerTest(&setup.TestConfig{EnableMockRepositories: true})
	if err != nil {
		return nil, err
	}

	err = test.RepositoryCreateTestUser(container, registry)
	if err != nil {
		return nil, err
	}

	err = test.RepositoryCreateTestDomain(container, registry)
	if err != nil {
		return nil, err
	}

	err = test.RepositoryCreateTestDomainRole(container, registry)
	if err != nil {
		return nil, err
	}

	if withData {
		err = test.RepositoryCreateTestUserDomainRole(container, registry)
		if err != nil {
			return nil, err
		}
	}

	return container.UserDomainRoleService.(*database.UserDomainRoleService), nil
}
