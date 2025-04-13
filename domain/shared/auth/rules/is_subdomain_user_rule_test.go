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

func TestIsUserOnExactDomain(t *testing.T) {
	user := &model.User{
		Username: "user",
		DomainRoles: []*model.UserDomainRole{
			{
				Role:   "admin",
				Domain: "example.com",
			},
		},
	}

	rule := rules.IsSubdomainUserRule{}

	if !assert.True(t, rule.CheckAccessOnResource(user, "example.com"), "User should be domain user on example.com") {
		return
	}
}

func TestIsUserOnSubdomain(t *testing.T) {
	user := &model.User{
		Username: "user",
		DomainRoles: []*model.UserDomainRole{
			{
				Role:   "readonly",
				Domain: "dev.example.com",
			},
		},
	}

	rule := rules.IsSubdomainUserRule{}

	if !assert.True(t, rule.CheckAccessOnResource(user, "example.com"), "User should be domain user on dev.example.com") {
		return
	}
}

func TestIsUserOnAnyDomain(t *testing.T) {
	user := &model.User{
		Username: "user",
		DomainRoles: []*model.UserDomainRole{
			{
				Role:   "readonly",
				Domain: "dev.example.com",
			},
		},
	}

	rule := rules.IsSubdomainUserRule{}

	if !assert.True(t, rule.CheckAccessOnResource(user, "*"), "User should be domain user on *") {
		return
	}
}

func TestIsUserOnDifferentSubdomainFails(t *testing.T) {
	user := &model.User{
		Username: "user",
		DomainRoles: []*model.UserDomainRole{
			{
				Role:   "readonly",
				Domain: "dev.example.com",
			},
		},
	}

	rule := rules.IsSubdomainUserRule{}

	if !assert.False(t, rule.CheckAccessOnResource(user, "test.example.com"), "User should not be domain user on test.example.com") {
		return
	}
}

func TestIsUserOnTopDomainFails(t *testing.T) {
	user := &model.User{
		Username: "user",
		DomainRoles: []*model.UserDomainRole{
			{
				Role:   "admin",
				Domain: "example.com",
			},
		},
	}

	rule := rules.IsSubdomainUserRule{}

	if !assert.False(t, rule.CheckAccessOnResource(user, "dev.example.com"), "User should not be domain user on dev.example.com") {
		return
	}
}
