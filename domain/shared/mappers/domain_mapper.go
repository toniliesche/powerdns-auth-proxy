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

package mappers

import (
	"powerdns-auth-proxy/domain/shared/database/model"
	"powerdns-auth-proxy/domain/shared/model/management"
)

type DomainMapper struct {
}

func (m *DomainMapper) MapDatabaseToMinimalDto(domain *model.Domain) *management.DomainMinimal {
	return &management.DomainMinimal{
		ID:        domain.ID,
		Fqdn:      domain.Fqdn,
		UpdatedAt: domain.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func (m *DomainMapper) MapDatabaseToMinimalDtoList(domains []*model.Domain) []*management.DomainMinimal {
	dtos := make([]*management.DomainMinimal, 0)

	for _, domain := range domains {
		dtos = append(dtos, m.MapDatabaseToMinimalDto(domain))
	}

	return dtos
}

func (m *DomainMapper) MapDatabaseToDto(domain *model.Domain) *management.Domain {
	return &management.Domain{
		ID:        domain.ID,
		Fqdn:      domain.Fqdn,
		CreatedAt: domain.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: domain.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func (m *DomainMapper) MapDatabaseToDtoList(domains []*model.Domain) []*management.Domain {
	dtos := make([]*management.Domain, 0)

	for _, domain := range domains {
		dtos = append(dtos, m.MapDatabaseToDto(domain))
	}

	return dtos
}

func (m *DomainMapper) MapCreatePayloadToDb(domain *management.DomainCreatePayload) *model.Domain {
	return &model.Domain{
		Fqdn: domain.Fqdn,
	}
}
