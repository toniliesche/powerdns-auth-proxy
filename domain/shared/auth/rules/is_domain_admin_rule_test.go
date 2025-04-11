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

package rules_test

import (
	"github.com/stretchr/testify/assert"
	"powerdns-auth-proxy/domain/shared/auth/model"
	"powerdns-auth-proxy/domain/shared/auth/rules"
	"testing"
)

func TestUserIsDomainAdmin(t *testing.T) {
	user := &model.User{
		Username: "admin",
		DomainRoles: []*model.UserDomainRole{
			{
				Domain: "example.com",
				Role:   "admin",
			},
		},
	}

	rule := &rules.IsDomainAdminRule{}

	if !assert.True(t, rule.CheckAccessOnResource(user, "example.com"), "User should be domain admin on example.com") {
		return
	}
}

func TestUserIsNotDomainAdmin(t *testing.T) {
	user := &model.User{
		Username: "user",
		DomainRoles: []*model.UserDomainRole{
			{
				Domain: "example.com",
				Role:   "user",
			},
		},
	}

	rule := &rules.IsDomainAdminRule{}

	if !assert.False(t, rule.CheckAccessOnResource(user, "example.com"), "User should not be domain admin on example.com") {
		return
	}

}

func TestUserIsDomainAdminOnOtherDomain(t *testing.T) {
	user := &model.User{
		Username: "user",
		DomainRoles: []*model.UserDomainRole{
			{
				Domain: "otherexample.com",
			},
		},
	}

	rule := &rules.IsDomainAdminRule{}

	if !assert.False(t, rule.CheckAccessOnResource(user, "example.com"), "User should not be domain admin on example.com") {
		return
	}
}

func TestUserIsReadonlyOnDomain(t *testing.T) {
	user := &model.User{
		Username: "user",
		DomainRoles: []*model.UserDomainRole{
			{
				Domain: "example.com",
				Role:   "readonly",
			},
		},
	}

	rule := &rules.IsDomainAdminRule{}

	if !assert.False(t, rule.CheckAccessOnResource(user, "example.com"), "User should not be domain admin on example.com") {
		return
	}
}
