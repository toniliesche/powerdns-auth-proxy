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

func TestUserDomainRoleValidation(t *testing.T) {
	userDomainRole := &importer.UserDomainRole{
		Domain: DomainFqdn,
		Role:   DomainRoleName,
	}

	if !assert.Nil(t, userDomainRole.Validate(), "should not return error when domain role is valid") {
		return
	}

	userDomainRole = &importer.UserDomainRole{
		Domain: DomainFqdn,
	}

	if !assert.NotNil(t, userDomainRole.Validate(), "should return error when role is empty") {
		return
	}

	userDomainRole = &importer.UserDomainRole{
		Role: DomainRoleName,
	}

	if !assert.NotNil(t, userDomainRole.Validate(), "should return error when domain is empty") {
		return
	}
}
