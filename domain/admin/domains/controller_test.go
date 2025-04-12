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

package domains_test

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"powerdns-auth-proxy/domain/admin/domains"
	"powerdns-auth-proxy/domain/admin/services/database"
	"powerdns-auth-proxy/domain/shared/basics"
	"powerdns-auth-proxy/domain/shared/controller"
	"powerdns-auth-proxy/domain/shared/model/management"
	"powerdns-auth-proxy/domain/shared/setup"
	"powerdns-auth-proxy/domain/test"
	"testing"
)

func TestConfigureRoutes(t *testing.T) {
	router := gin.New()
	domainController, err := getController(false)
	if !assert.NoError(t, err, fmt.Sprintf("could not initialize TestConfigureRoutes: %s", err)) {
		return
	}

	domainController.ConfigureGroupRoutes(router.Group("/api"))

	if !assert.Equal(t, 4, len(router.Routes()), "There should be 4 routes configured") {
		return
	}
}

func TestCallListDomainsEndpoint(t *testing.T) {
	domainController, err := getController(true)
	if !assert.NoError(t, err, fmt.Sprintf("could not initialize TestCallListDomainsEndpoint: %s", err)) {
		return
	}

	test.RunRequest(t, domainController, "/admin", "/admin/v1/domains", "GET", 200)
}

func TestCallCreateDomainEndpoint(t *testing.T) {
	domainController, err := getController(false)
	if !assert.NoError(t, err, fmt.Sprintf("could not initialize TestCallCreateDomainEndpoint: %s", err)) {
		return
	}

	test.RunRequest(t, domainController, "/admin", "/admin/v1/domains", "POST", 200, &management.DomainCreatePayload{Fqdn: "testdomain.com"})
}

func TestCallGetDomainEndpoint(t *testing.T) {
	domainController, err := getController(true)
	if !assert.NoError(t, err, fmt.Sprintf("could not initialize TestCallGetDomainEndpoint: %s", err)) {
		return
	}

	test.RunRequest(t, domainController, "/admin", "/admin/v1/domains/details/1", "GET", 200)
}

func TestCallDeleteDomainEndpoint(t *testing.T) {
	domainController, err := getController(true)
	if !assert.NoError(t, err, fmt.Sprintf("could not initialize TestCallGetDomainEndpoint: %s", err)) {
		return
	}

	test.RunRequest(t, domainController, "/admin", "/admin/v1/domains/1", "DELETE", 200)
}

func TestProvideDomainControllerFailsOnMissingBaseController(t *testing.T) {
	container := &basics.InjectionContainer{}

	domainController, err := domains.ProvideDomainController(container)
	assert.Error(t, err, "provide domain controller method should return an error")
	assert.Equal(t, "domains controller could not be created: base controller could not be resolved", err.Error())
	assert.Nil(t, domainController, "provide domain controller method should not return a controller")
}

func TestProvideDomainControllerFailsOnMissingDomainService(t *testing.T) {
	container := &basics.InjectionContainer{}
	container.BaseControllerAdminAPI = &controller.BaseController{}

	domainController, err := domains.ProvideDomainController(container)
	assert.Error(t, err, "provide domain controller method should return an error")
	assert.Equal(t, "domains controller could not be created: domain service could not be resolved", err.Error())
	assert.Nil(t, domainController, "provide domain controller method should not return a controller")
}

func TestProvideDomainControllerSucceeds(t *testing.T) {
	container := &basics.InjectionContainer{}
	container.BaseControllerAdminAPI = &controller.BaseController{}
	container.DomainService = &database.DomainService{}

	domainController, err := domains.ProvideDomainController(container)
	assert.NoError(t, err, "provide domain controller method should succeed")
	assert.NotNil(t, domainController, "provide domain controller method should return a controller")
}

func getController(create bool) (*domains.DomainController, error) {
	container, err := setup.InitContainerTest(&setup.TestConfig{EnableMockAuthentication: true, EnableMockForwardService: true, EnableMockRepositories: true})
	if err != nil {
		return nil, err
	}

	if create {
		payload := &management.DomainCreatePayload{Fqdn: "testdomain.com"}
		_, err := container.DomainService.CreateDomain(payload)
		if err != nil {
			return nil, err
		}
	}

	return container.AdminDomainController.(*domains.DomainController), nil
}
