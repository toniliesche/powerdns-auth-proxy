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

package auth_test

import (
	"github.com/stretchr/testify/assert"
	"powerdns-auth-proxy/domain/shared/auth"
	"testing"
)

func TestRuleSetProvider(t *testing.T) {
	provider := auth.GetRuleSetProvider()

	ruleSet, err := provider.GetRuleSet(auth.Admin)

	if !assert.NoError(t, err, "Error should be nil") {
		return
	}

	if !assert.NotNil(t, ruleSet, "RuleSet should not be nil") {
		return
	}
}

func TestRuleSetProviderInvalidRuleSet(t *testing.T) {
	provider := auth.GetRuleSetProvider()

	ruleSet, err := provider.GetRuleSet("invalid")

	if !assert.Error(t, err, "Error should not be nil") {
		return
	}

	if !assert.Nil(t, ruleSet, "RuleSet should be nil") {
		return
	}
}
