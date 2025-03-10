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

package mock

import (
	"powerdns-auth-proxy/domain/shared/database/model"
	errors2 "powerdns-auth-proxy/domain/shared/database/repository/errors"
)

type DomainRepositoryMock struct {
	domains map[string]*model.Domain
	counter uint
}

func (r *DomainRepositoryMock) SaveNewDomain(domain *model.Domain) error {
	if r.CheckExistenceByFqdn(domain.Fqdn) {
		return errors2.NewItemAlreadyExistsError("domain already exists")
	}

	domain.ID = r.counter
	r.domains[domain.Fqdn] = domain
	r.counter++

	return nil
}

func (r *DomainRepositoryMock) CheckExistenceByFqdn(fqdn string) bool {
	_, found := r.domains[fqdn]

	return found
}

func (r *DomainRepositoryMock) FetchDomainById(id uint) (*model.Domain, error) {
	for _, domain := range r.domains {
		if domain.ID == id {
			return domain, nil
		}
	}

	return nil, errors2.NewItemNotFoundError("domain not found")
}

func (r *DomainRepositoryMock) FetchDomainByFqdn(fqdn string) (*model.Domain, error) {
	domain, found := r.domains[fqdn]
	if !found {
		return nil, errors2.NewItemNotFoundError("domain not found")
	}

	return domain, nil
}

func (r *DomainRepositoryMock) FindAll() ([]*model.Domain, error) {
	domains := make([]*model.Domain, 0, len(r.domains))
	for _, domain := range r.domains {
		domains = append(domains, domain)
	}

	return domains, nil
}

func (r *DomainRepositoryMock) DeleteDomain(fqdn string) error {
	if !r.CheckExistenceByFqdn(fqdn) {
		return errors2.NewItemNotFoundError("domain not found")
	}

	delete(r.domains, fqdn)

	return nil
}

func (r *DomainRepositoryMock) DeleteDomainById(id uint) error {
	for _, domain := range r.domains {
		if domain.ID == id {
			delete(r.domains, domain.Fqdn)
			return nil
		}
	}

	return errors2.NewItemNotFoundError("domain not found")
}

func ProvideDomainRepositoryMock() (*DomainRepositoryMock, error) {
	return &DomainRepositoryMock{
		domains: make(map[string]*model.Domain),
		counter: 1,
	}, nil
}
