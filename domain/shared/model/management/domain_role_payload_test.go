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

package management_test

import (
	"github.com/stretchr/testify/assert"
	"powerdns-auth-proxy/domain/shared/model/management"
	"testing"
)

func TestVerifyDomainRolePayloadFailsOnMissingRolesList(t *testing.T) {
	payload := &management.UserDomainRolePayload{
		DomainRoles: nil,
	}

	err := payload.Verify()
	assert.Error(t, err, "verify should return an error")
	assert.Equal(t, err.Error(), "domain_roles is required")
}

func TestVerifyDomainRolePayloadFailsOnEmptyRolesList(t *testing.T) {
	payload := &management.UserDomainRolePayload{
		DomainRoles: []*management.DomainRole{},
	}

	err := payload.Verify()
	assert.Error(t, err, "verify should return an error")
	assert.Equal(t, err.Error(), "domain_roles must contain at least one item")
}

func TestVerifyDomainRolePayloadFailsOnInvalidRolesItem(t *testing.T) {
	payload := &management.UserDomainRolePayload{
		DomainRoles: []*management.DomainRole{
			{
				Domain: "",
				Role:   "admin",
			},
		},
	}

	err := payload.Verify()
	assert.Error(t, err, "verify should return an error")
	assert.Equal(t, err.Error(), "domain_roles[0]: domain is required")
}

func TestVerifyDomainRolePayloadSucceeds(t *testing.T) {
	payload := &management.UserDomainRolePayload{
		DomainRoles: []*management.DomainRole{
			{
				Domain: "example.com",
				Role:   "admin",
			},
		},
	}

	err := payload.Verify()
	assert.NoError(t, err, "verify should not return an error")
}
