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

package services_test

import (
	"github.com/stretchr/testify/assert"
	"powerdns-auth-proxy/domain/shared/auth/model"
	"powerdns-auth-proxy/domain/shared/auth/services"
	dbmodel "powerdns-auth-proxy/domain/shared/database/model"
	"testing"
)

func TestMapDatabaseUserToUser(t *testing.T) {
	mapper := services.UserMapper{}

	dbUser := &dbmodel.User{
		Username: "user",
		UserRoles: []*dbmodel.UserRole{
			{
				Role: &dbmodel.Role{
					Name: "admin",
				},
				User: &dbmodel.User{
					Username: "user",
				},
			},
		},
		UserDomainRoles: []*dbmodel.UserDomainRole{
			{
				Domain: &dbmodel.Domain{
					Fqdn: "example.com",
				},
				DomainRole: &dbmodel.DomainRole{
					Name: "admin",
				},
				User: &dbmodel.User{
					Username: "user",
				},
			},
		},
	}

	user, err := mapper.FromDatabase(dbUser)
	if !assert.NoError(t, err, "unexpected error") {
		return
	}

	if !assert.IsType(t, &model.User{}, user, "unexpected user type") {
		return
	}

	if !assert.Equal(t, dbUser.Username, user.Username, "unexpected username") {
		return
	}

	if !assert.Len(t, user.Roles, 1, "unexpected number of roles") {
		return
	}

	if !assert.Equal(t, "admin", user.Roles[0], "unexpected role") {
		return
	}

	if !assert.Len(t, user.DomainRoles, 1, "unexpected number of domain roles") {
		return
	}

	if !assert.Equal(t, "example.com", user.DomainRoles[0].Domain, "unexpected domain") {
		return
	}

	if !assert.Equal(t, "admin", user.DomainRoles[0].Role, "unexpected role") {
		return
	}
}
