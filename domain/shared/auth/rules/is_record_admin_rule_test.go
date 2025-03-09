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

func TestIsRecordAdmin(t *testing.T) {
	user := &model.User{
		Username: "user",
		DomainRoles: []*model.UserDomainRole{
			{
				Role:   "record_admin",
				Domain: "example.com",
			},
		},
	}

	rule := &rules.IsRecordAdminRule{}

	if !assert.True(t, rule.CheckAccessOnResource(user, "example.com"), "User should be record admin") {
		return
	}
}

func TestIsNotRecordAdmin(t *testing.T) {
	user := &model.User{
		Username: "user",
	}

	rule := &rules.IsRecordAdminRule{}

	if !assert.False(t, rule.CheckAccessOnResource(user, "example.com"), "User should not be record admin") {
		return
	}
}

func TestIsRecordAdminOnOtherDomain(t *testing.T) {
	user := &model.User{
		Username: "user",
		DomainRoles: []*model.UserDomainRole{
			{
				Role:   "record_admin",
				Domain: "otherexample.com",
			},
		},
	}

	rule := &rules.IsRecordAdminRule{}

	if !assert.False(t, rule.CheckAccessOnResource(user, "example.com"), "User should not be record admin") {
		return
	}
}

func TestIsReadOnlyOnDomain(t *testing.T) {
	user := &model.User{
		Username: "user",
		DomainRoles: []*model.UserDomainRole{
			{
				Role:   "readonly",
				Domain: "example.com",
			},
		},
	}

	rule := &rules.IsRecordAdminRule{}

	if !assert.False(t, rule.CheckAccessOnResource(user, "example.com"), "User should not be record admin") {
		return
	}
}
