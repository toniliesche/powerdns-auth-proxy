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

func TestCreateDomain(t *testing.T) {
	registry := test.NewRegistry()
	service, err := ProvideTestDomainService(registry, false)
	if err != nil {
		t.Error(err)
		return
	}

	payload := &management.DomainCreatePayload{
		FQDN: test.DomainFQDN,
	}

	_, err = service.CreateDomain(payload)
	if !assert.NoError(t, err, "failed to create domain") {
		return
	}
}

func TestGetDomainById(t *testing.T) {
	registry := test.NewRegistry()
	service, err := ProvideTestDomainService(registry, true)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = service.GetDomainByID(registry.GetUint("domainId"))
	if !assert.NoError(t, err, "failed to get domain by ID") {
		return
	}
}

func TestListDomains(t *testing.T) {
	registry := test.NewRegistry()
	service, err := ProvideTestDomainService(registry, true)
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
	service, err := ProvideTestDomainService(registry, true)
	if err != nil {
		t.Error(err)
		return
	}

	err = service.DeleteDomainByID(registry.GetUint("domainId"))
	if !assert.NoError(t, err, "failed to delete domain by ID") {
		return
	}
}

func TestDeleteDomainByFQDN(t *testing.T) {
	registry := test.NewRegistry()
	service, err := ProvideTestDomainService(registry, true)
	if err != nil {
		t.Error(err)
		return
	}

	err = service.DeleteDomain(test.DomainFQDN)
	if !assert.NoError(t, err, "failed to delete domain by FQDN") {
		return
	}
}

func ProvideTestDomainService(registry *test.Registry, withData bool) (*database.DomainService, error) {
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
