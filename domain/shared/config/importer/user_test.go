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

package importer_test

import (
	"github.com/stretchr/testify/assert"
	"powerdns-auth-proxy/domain/shared/config/importer"
	"testing"
)

func TestUserValidation(t *testing.T) {
	user := &importer.User{
		Username: UserUsername,
	}

	if !assert.NotNil(t, user.Validate(), "should return error when password is empty") {
		return
	}

	user = &importer.User{
		Password: UserPassword,
	}

	if !assert.NotNil(t, user.Validate(), "should return error when username is empty") {
		return
	}

	user = &importer.User{
		Username: UserUsername,
		Password: UserPassword,
	}

	if !assert.Nil(t, user.Validate(), "should not return error when username and password are not empty") {
		return
	}
}

func TestUserValidationForUserRoles(t *testing.T) {
	user := &importer.User{
		Username: UserUsername,
		Password: UserPassword,
		UserRoles: []*importer.UserRole{
			{
				Role: RoleName,
			},
		},
	}

	if !assert.Nil(t, user.Validate(), "should not return error when user role is valid") {
		return
	}

	user = &importer.User{
		Username: UserUsername,
		Password: UserPassword,
		UserRoles: []*importer.UserRole{
			{},
		},
	}

	if !assert.NotNil(t, user.Validate(), "should return error when user role is invalid") {
		return
	}
}

func TestUserValidationForUserDomainRoles(t *testing.T) {
	user := &importer.User{
		Username: UserUsername,
		Password: UserPassword,
		UserDomainRoles: []*importer.UserDomainRole{
			{
				Domain: DomainFQDN,
			},
		},
	}

	if !assert.NotNil(t, user.Validate(), "should return error when domain role is invalid") {
		return
	}

	user = &importer.User{
		Username: UserUsername,
		Password: UserPassword,
		UserRoles: []*importer.UserRole{
			{
				Role: DomainRoleName,
			},
		},
	}

	if !assert.Nil(t, user.Validate(), "should not return error when user role is valid") {
		return
	}
}

func TestUserDetectionForGlobalRoles(t *testing.T) {
	user := &importer.User{
		Username: UserUsername,
		Password: UserPassword,
		UserRoles: []*importer.UserRole{
			{
				Role: RoleName,
			},
		},
	}

	if !assert.True(t, user.HasGlobalRoles(), "should return true when user has global roles") {
		return
	}

	user = &importer.User{
		Username: UserUsername,
		Password: UserPassword,
	}

	if !assert.False(t, user.HasGlobalRoles(), "should return false when user has no global roles") {
		return
	}
}

func TestUserDetectionForDomainRoles(t *testing.T) {
	user := &importer.User{
		Username: UserUsername,
		Password: UserPassword,
		UserDomainRoles: []*importer.UserDomainRole{
			{
				Domain: DomainFQDN,
				Role:   DomainRoleName,
			},
		},
	}

	if !assert.True(t, user.HasDomainRoles(), "should return true when user has domain roles") {
		return
	}

	user = &importer.User{
		Username: UserUsername,
		Password: UserPassword,
	}

	if !assert.False(t, user.HasDomainRoles(), "should return false when user has no domain roles") {
		return
	}
}
