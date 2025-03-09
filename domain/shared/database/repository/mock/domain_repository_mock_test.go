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

func TestDomainCanBeSaved(t *testing.T) {
	var err error
	registry := test.NewRegistry()
	repo := ProvideTestMockDomainRepository(registry, t, false)

	if !assert.NotNil(t, repo, "failed to provide test domain repository") {
		return
	}

	err = repo.SaveNewDomain(&model.Domain{FQDN: test.DomainFQDN})

	if !assert.NoError(t, err, "failed to save new domain") {
		return
	}
}

func TestDomainExistenceCanBeChecked(t *testing.T) {
	registry := test.NewRegistry()
	repo := ProvideTestMockDomainRepository(registry, t, true)

	if !assert.NotNil(t, repo, "failed to provide test domain repository") {
		return
	}

	found := repo.CheckExistenceByFQDN(test.DomainFQDN)

	if !assert.True(t, found, "expected domain to exist") {
		return
	}
}

func TestDomainNonExistenceCanBeChecked(t *testing.T) {
	registry := test.NewRegistry()
	repo := ProvideTestMockDomainRepository(registry, t, true)

	if !assert.NotNil(t, repo, "failed to provide test domain repository") {
		return
	}

	found := repo.CheckExistenceByFQDN(test.NonExistingDomainFQDN)

	if !assert.False(t, found, "expected domain to not exist") {
		return
	}
}

func TestDomainCanBeFetchedByFQDN(t *testing.T) {
	var err error
	registry := test.NewRegistry()
	repo := ProvideTestMockDomainRepository(registry, t, true)

	if !assert.NotNil(t, repo, "failed to provide test domain repository") {
		return
	}

	domain, err := repo.FetchDomainByFQDN(test.DomainFQDN)

	if !assert.NoError(t, err, "failed to fetch domain") {
		return
	}

	if !assert.IsType(t, &model.Domain{}, domain, "expected domain to be of type model.Domain") {
		return
	}

	if !assert.Equal(t, test.DomainFQDN, domain.FQDN, "expected domain to have correct FQDN") {
		return
	}
}

func TestCorrectDomainWillBeFetched(t *testing.T) {
	var err error
	registry := test.NewRegistry()
	repo := ProvideTestMockDomainRepository(registry, t, true)

	if !assert.NotNil(t, repo, "failed to provide test domain repository") {
		return
	}

	domain, err := repo.FetchDomainByFQDN(test.NonExistingDomainFQDN)

	if !assert.NotNilf(t, err, "expected error to be returned") {
		return
	}

	if !assert.Nil(t, domain, "expected domain to be nil") {
		return
	}
}

func TestDomainCanBeDeleted(t *testing.T) {
	var err error
	registry := test.NewRegistry()
	repo := ProvideTestMockDomainRepository(registry, t, true)

	if !assert.NotNil(t, repo, "failed to provide test domain repository") {
		return
	}

	err = repo.DeleteDomain(test.DomainFQDN)

	if !assert.NoError(t, err, "failed to delete domain") {
		return
	}

	found := repo.CheckExistenceByFQDN(test.DomainFQDN)

	if !assert.False(t, found, "expected domain to not exist") {
		return
	}
}

func ProvideTestMockDomainRepository(registry *test.Registry, t *testing.T, withData bool) interfaces.DomainRepositoryInterface {
	container, err := setup.InitContainerTest(&setup.TestConfig{EnableMockRepositories: true})
	if !assert.NoError(t, err, "failed to setup test database") {
		return nil
	}

	if withData {
		err = test.RepositoryCreateTestDomain(container, registry)
		if !assert.NoError(t, err, "failed to setup test database") {
			return nil
		}
	}

	return container.DomainRepository
}
