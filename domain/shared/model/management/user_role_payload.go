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
	"strings"
)

type UserRolePayload struct {
	Roles []string `json:"roles"`
}

func (p *UserRolePayload) ContainsAdminRole() bool {
	for _, role := range p.Roles {
		if role == "admin" || role == "superadmin" {
			return true
		}
	}

	return false
}

func (p *UserRolePayload) Verify() errors.HTTPError {
	if len(p.Roles) == 0 {
		return errors.NewBadRequestError(fmt.Errorf("roles is required"))
	}

	for _, role := range p.Roles {
		if strings.TrimSpace(role) == "" {
			return errors.NewBadRequestError(fmt.Errorf("role cannot be empty"))
		}
	}

	return nil
}
