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

package model_test

import (
	"github.com/stretchr/testify/assert"
	"powerdns-auth-proxy/domain/shared/database/model"
	"testing"
)

func TestCheckForGlobalRole(t *testing.T) {
	user := GetTestUser()

	if !assert.True(t, user.HasGlobalRole(RoleName), "admin role not detected") {
		return
	}
}

func TestCheckForNonExistingGlobalRole(t *testing.T) {
	user := GetTestUser()

	if !assert.False(t, user.HasGlobalRole(NonExistingRoleName), "non-existing role not detected") {
		return
	}
}

func TestGetGlobalRole(t *testing.T) {
	user := GetTestUser()
	role, err := user.GetGlobalRole(RoleName)

	if !assert.NoError(t, err, "failed to get global role") {
		return
	}

	if !assert.NotNil(t, role, "role not found") {
		return
	}
}

func TestGetNonExistingGlobalRole(t *testing.T) {
	user := GetTestUser()
	role, err := user.GetGlobalRole(NonExistingRoleName)

	if !assert.Error(t, err, "non-existing role issue not detected") {
		return
	}

	if !assert.Nil(t, role, "non-existing role issue not detected") {
		return
	}
}

func TestCheckForDomainRole(t *testing.T) {
	user := GetTestUser()

	if !assert.True(t, user.HasDomainRole(DomainFqdn, DomainRoleName), "domain role not detected") {
		return
	}
}

func TestCheckForNonExistingDomainRole(t *testing.T) {
	user := GetTestUser()

	if !assert.False(t, user.HasDomainRole(NonExistingDomainRoleName, DomainRoleName), "non-existing domain role issue not detected") {
		return
	}

	if !assert.False(t, user.HasDomainRole(DomainRoleName, NonExistingDomainFqdn), "non-existing domain role issue not detected") {
		return
	}
}

func TestGetDomainRole(t *testing.T) {
	user := GetTestUser()
	domainRole, err := user.GetDomainRole(DomainFqdn, DomainRoleName)

	if !assert.NoError(t, err, "failed to get domain role") {
		return
	}

	if !assert.NotNil(t, domainRole, "domain role not found") {
		return
	}
}

func TestGetNonExistingDomainRole(t *testing.T) {
	user := GetTestUser()
	domainRole, err := user.GetDomainRole(NonExistingDomainFqdn, DomainRoleName)

	if !assert.Error(t, err, "non-existing domain role issue not detected") {
		return
	}

	if !assert.Nil(t, domainRole, "non-existing domain role issue not detected") {
		return
	}

	domainRole, err = user.GetDomainRole(DomainFqdn, NonExistingDomainRoleName)

	if !assert.Error(t, err, "non-existing domain role issue not detected") {
		return
	}

	if !assert.Nil(t, domainRole, "non-existing domain role issue not detected") {
		return
	}
}

func GetTestUser() *model.User {
	return &model.User{
		UserDomainRoles: []*model.UserDomainRole{
			{
				DomainRole: &model.DomainRole{
					Name: DomainRoleName,
				},
				Domain: &model.Domain{
					Fqdn: DomainFqdn,
				},
			},
		},
		UserRoles: []*model.UserRole{
			{
				Role: &model.Role{
					Name: RoleName,
				},
			},
		},
	}
}
