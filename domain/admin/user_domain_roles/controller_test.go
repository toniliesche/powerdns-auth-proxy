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

package user_domain_roles_test

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"powerdns-auth-proxy/domain/admin/user_domain_roles"
	"powerdns-auth-proxy/domain/shared/model/management"
	"powerdns-auth-proxy/domain/shared/setup"
	"powerdns-auth-proxy/domain/test"
	"testing"
)

func TestConfigureRoutes(t *testing.T) {
	router := gin.New()
	registry := test.NewRegistry()
	controller, err := getController(registry, false)
	if !assert.NoError(t, err, fmt.Sprintf("could not initialize TestConfigureRoutes: %s", err)) {
		return
	}

	controller.ConfigureGroupRoutes(router.Group("/api"))

	if !assert.Equal(t, 3, len(router.Routes()), "There should be 3 routes configured") {
		return
	}
}

func TestCallListUserDomainRolesEndpoint(t *testing.T) {
	registry := test.NewRegistry()
	controller, err := getController(registry, false)
	if !assert.NoError(t, err, fmt.Sprintf("could not initialize TestCallListUserDomainRolesEndpoint: %s", err)) {
		return
	}

	test.RunRequest(t, controller, "/admin", fmt.Sprintf("/admin/v1/users/%s/domain-roles", registry.Get("userId")), "GET", 200)
}

func TestCallGrantUserDomainRoleEndpoint(t *testing.T) {
	registry := test.NewRegistry()
	controller, err := getController(registry, false)
	if !assert.NoError(t, err, fmt.Sprintf("could not initialize TestCallAddUserDomainRoleEndpoint: %s", err)) {
		return
	}

	payload := &management.UserDomainRolePayload{
		DomainRoles: []*management.DomainRole{
			{
				Domain: test.DomainFqdn,
				Role:   test.DomainRoleName,
			},
		},
	}
	test.RunRequest(t, controller, "/admin", fmt.Sprintf("/admin/v1/users/%s/domain-roles/grant", registry.Get("userId")), "POST", 200, payload)
}

func TestCallRevokeUserDomainRoleEndpoint(t *testing.T) {
	registry := test.NewRegistry()
	controller, err := getController(registry, true)
	if !assert.NoError(t, err, fmt.Sprintf("could not initialize TestCallRevokeUserDomainRoleEndpoint: %s", err)) {
		return
	}

	payload := &management.UserDomainRolePayload{
		DomainRoles: []*management.DomainRole{
			{
				Domain: test.DomainFqdn,
				Role:   test.DomainRoleName,
			},
		},
	}
	test.RunRequest(t, controller, "/admin", fmt.Sprintf("/admin/v1/users/%s/domain-roles/revoke", registry.Get("userId")), "POST", 200, payload)
}

func getController(registry *test.Registry, withData bool) (*user_domain_roles.UserDomainRolesController, error) {
	container, err := setup.InitContainerTest(&setup.TestConfig{EnableMockAuthentication: true, EnableMockForwardService: true, EnableMockRepositories: true})
	if err != nil {
		return nil, err
	}

	err = test.RepositoryCreateTestUser(container, registry)
	if err != nil {
		return nil, err
	}

	err = test.RepositoryCreateTestDomainRole(container, registry)
	if err != nil {
		return nil, err
	}

	err = test.RepositoryCreateTestDomain(container, registry)
	if err != nil {
		return nil, err
	}

	if withData {
		err = test.RepositoryCreateTestUserDomainRole(container, registry)
		if err != nil {
			return nil, err
		}
	}

	container.ResponseWriter.SetDebug(true)

	return container.AdminUserDomainRolesController.(*user_domain_roles.UserDomainRolesController), nil
}
