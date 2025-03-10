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

package management

import (
	"fmt"
	"powerdns-auth-proxy/domain/shared/http/errors"
)

type UserDomainRolePayload struct {
	DomainRoles []*DomainRole `json:"domain_roles"`
}

func (p *UserDomainRolePayload) Verify() errors.HTTPError {
	if p.DomainRoles == nil {
		return errors.NewBadRequestError(fmt.Errorf("domain_roles is required"))
	}

	if len(p.DomainRoles) == 0 {
		return errors.NewBadRequestError(fmt.Errorf("domain_roles must contain at least one item"))
	}

	for index, role := range p.DomainRoles {
		if err := role.Verify(); err != nil {
			return errors.NewBadRequestError(fmt.Errorf("domain_roles[%d]: %w", index, err))
		}
	}

	return nil
}

func (p *UserDomainRolePayload) GetDomainRoles() []*DomainRole {
	domainRoles := make([]*DomainRole, 0, len(p.DomainRoles))
	for _, domainRole := range p.DomainRoles {
		domainRoles = append(domainRoles, &DomainRole{
			Domain: domainRole.Domain,
			Role:   domainRole.Role,
		})
	}

	return domainRoles
}
