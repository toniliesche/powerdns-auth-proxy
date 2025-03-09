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
	"powerdns-auth-proxy/domain/shared/model/management"
	"powerdns-auth-proxy/domain/shared/setup"
	"powerdns-auth-proxy/domain/test"
	"testing"
)

func TestConfigureRoutes(t *testing.T) {
	router := gin.New()
	controller, err := getController(false)
	if !assert.NoError(t, err, fmt.Sprintf("could not initialize TestConfigureRoutes: %s", err)) {
		return
	}

	controller.ConfigureGroupRoutes(router.Group("/api"))

	if !assert.Equal(t, 4, len(router.Routes()), "There should be 4 routes configured") {
		return
	}
}

func TestCallListDomainsEndpoint(t *testing.T) {
	controller, err := getController(true)
	if !assert.NoError(t, err, fmt.Sprintf("could not initialize TestCallListDomainsEndpoint: %s", err)) {
		return
	}

	test.RunRequest(t, controller, "/admin", "/admin/v1/domains", "GET", 200)
}

func TestCallCreateDomainEndpoint(t *testing.T) {
	controller, err := getController(false)
	if !assert.NoError(t, err, fmt.Sprintf("could not initialize TestCallCreateDomainEndpoint: %s", err)) {
		return
	}

	test.RunRequest(t, controller, "/admin", "/admin/v1/domains", "POST", 200, &management.DomainCreatePayload{FQDN: "testdomain.com"})
}

func TestCallGetDomainEndpoint(t *testing.T) {
	controller, err := getController(true)
	if !assert.NoError(t, err, fmt.Sprintf("could not initialize TestCallGetDomainEndpoint: %s", err)) {
		return
	}

	test.RunRequest(t, controller, "/admin", "/admin/v1/domains/details/1", "GET", 200)
}

func TestCallDeleteDomainEndpoint(t *testing.T) {
	controller, err := getController(true)
	if !assert.NoError(t, err, fmt.Sprintf("could not initialize TestCallGetDomainEndpoint: %s", err)) {
		return
	}

	test.RunRequest(t, controller, "/admin", "/admin/v1/domains/1", "DELETE", 200)
}

func getController(create bool) (*domains.DomainController, error) {
	container, err := setup.InitContainerTest(&setup.TestConfig{EnableMockAuthentication: true, EnableMockForwardService: true, EnableMockRepositories: true})
	if err != nil {
		return nil, err
	}

	if create {
		payload := &management.DomainCreatePayload{FQDN: "testdomain.com"}
		_, err := container.DomainService.CreateDomain(payload)
		if err != nil {
			return nil, err
		}
	}

	return container.AdminDomainController.(*domains.DomainController), nil
}
