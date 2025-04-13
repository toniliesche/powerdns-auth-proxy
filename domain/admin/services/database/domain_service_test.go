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

func TestCreateDomain(t *testing.T) {
	registry := test.NewRegistry()
	service, err := getDomainService(registry, false)
	if err != nil {
		t.Error(err)
		return
	}

	payload := &management.DomainCreatePayload{
		Fqdn: test.DomainFqdn,
	}

	_, err = service.CreateDomain(payload)
	if !assert.NoError(t, err, "failed to create domain") {
		return
	}
}

func TestGetDomainById(t *testing.T) {
	registry := test.NewRegistry()
	service, err := getDomainService(registry, true)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = service.GetDomainById(registry.GetUint("domainId"))
	if !assert.NoError(t, err, "failed to get domain by ID") {
		return
	}
}

func TestListDomains(t *testing.T) {
	registry := test.NewRegistry()
	service, err := getDomainService(registry, true)
	if err != nil {
		t.Error(err)
		return
	}

	domains, err := service.ListDomains()
	if !assert.NoError(t, err, "failed to list domains") {
		return
	}

	if !assert.NotEmpty(t, domains, "expected domains to not be empty") {
		return
	}
}

func TestDeleteDomainById(t *testing.T) {
	registry := test.NewRegistry()
	service, err := getDomainService(registry, true)
	if err != nil {
		t.Error(err)
		return
	}

	err = service.DeleteDomainById(registry.GetUint("domainId"))
	if !assert.NoError(t, err, "failed to delete domain by ID") {
		return
	}
}

func TestDeleteDomainByFqdn(t *testing.T) {
	registry := test.NewRegistry()
	service, err := getDomainService(registry, true)
	if err != nil {
		t.Error(err)
		return
	}

	err = service.DeleteDomain(test.DomainFqdn)
	if !assert.NoError(t, err, "failed to delete domain by Fqdn") {
		return
	}
}

func TestProvideDomainServiceFailsOnMissingDomainRepository(t *testing.T) {
	container := &basics.InjectionContainer{}

	_, err := database.NewDomainService(container)
	assert.Error(t, err, "provide domain service method should return an error")
	assert.Equal(t, "could not provide domain service: domain repository could not be resolved", err.Error())
}

func TestProvideDomainServiceSucceeds(t *testing.T) {
	container := &basics.InjectionContainer{}
	container.DomainRepository = &rdbms.DomainRepository{}

	service, err := database.NewDomainService(container)
	assert.NoError(t, err, "provide domain service method should succeed")
	assert.NotNil(t, service, "provide domain service method should return a service")
}

func getDomainService(registry *test.Registry, withData bool) (*database.DomainService, error) {
	container, err := setup.InitContainerTest(&setup.TestConfig{EnableMockRepositories: true})
	if err != nil {
		return nil, err
	}

	if withData {
		err := test.RepositoryCreateTestDomain(container, registry)
		if err != nil {
			return nil, err
		}
	}

	return container.DomainService.(*database.DomainService), nil
}
