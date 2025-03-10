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

package database

import (
	"powerdns-auth-proxy/domain/shared/basics"
	"powerdns-auth-proxy/domain/shared/database/repository/interfaces"
	"powerdns-auth-proxy/domain/shared/mappers"
	"powerdns-auth-proxy/domain/shared/model/management"
)

type DomainService struct {
	domainRepository interfaces.DomainRepositoryInterface
	mapper           *mappers.DomainMapper
}

func (s *DomainService) CreateDomain(payload *management.DomainCreatePayload) (uint, error) {
	domain := s.mapper.MapCreatePayloadToDb(payload)
	err := s.domainRepository.SaveNewDomain(domain)
	if err != nil {
		return 0, err
	}

	return domain.ID, nil
}

func (s *DomainService) ListDomains() ([]*management.DomainMinimal, error) {
	dbDomains, err := s.domainRepository.FindAll()
	if err != nil {
		return nil, err
	}

	return s.mapper.MapDatabaseToMinimalDtoList(dbDomains), nil
}

func (s *DomainService) GetDomainById(id uint) (*management.Domain, error) {
	dbDomain, err := s.domainRepository.FetchDomainById(id)
	if err != nil {
		return nil, err
	}

	return s.mapper.MapDatabaseToDto(dbDomain), nil
}

func (s *DomainService) DeleteDomain(fqdn string) error {
	return s.domainRepository.DeleteDomain(fqdn)
}

func (s *DomainService) DeleteDomainById(id uint) error {
	return s.domainRepository.DeleteDomainById(id)
}

func ProvideDomainService(container *basics.InjectionContainer) (*DomainService, error) {
	if container.DomainRepository == nil {
		return nil, basics.NewMissingDependencyError("could not provide domain service: domain repository could not be resolved")
	}

	return &DomainService{
		domainRepository: container.DomainRepository,
		mapper:           &mappers.DomainMapper{},
	}, nil
}
