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

package auth

import (
	"fmt"
	"powerdns-auth-proxy/domain/shared/auth/rules"
)

type RuleSetProvider struct {
	RuleSets map[string]RuleSet
}

const (
	Superadmin    = "superadmin"
	Admin         = "admin"
	DomainAdmin   = "domain_admin"
	RecordAdmin   = "record_admin"
	Reader        = "reader"
	SubdomainUser = "subdomain_user"
)

func (p *RuleSetProvider) GetRuleSet(name string) (*RuleSet, error) {
	if ruleSet, ok := p.RuleSets[name]; ok {
		return &ruleSet, nil
	}

	return nil, fmt.Errorf("rule set not found")
}

func GetRuleSetProvider() *RuleSetProvider {
	return &RuleSetProvider{
		RuleSets: map[string]RuleSet{
			Superadmin: {
				rules: []rules.Rule{
					&rules.IsSuperadminRule{},
				},
			},
			Admin: {
				rules: []rules.Rule{
					&rules.IsAdminRule{},
					&rules.IsSuperadminRule{},
				},
			},
			DomainAdmin: {
				rules: []rules.Rule{
					&rules.IsDomainAdminRule{},
					&rules.IsAdminRule{},
					&rules.IsSuperadminRule{},
				},
			},
			RecordAdmin: {
				rules: []rules.Rule{
					&rules.IsRecordAdminRule{},
					&rules.IsDomainAdminRule{},
					&rules.IsAdminRule{},
					&rules.IsSuperadminRule{},
				},
			},
			Reader: {
				rules: []rules.Rule{
					&rules.IsDomainReaderRule{},
					&rules.IsRecordAdminRule{},
					&rules.IsDomainAdminRule{},
					&rules.IsAdminRule{},
					&rules.IsSuperadminRule{},
				},
			},
			SubdomainUser: {
				rules: []rules.Rule{
					&rules.IsSubdomainUserRule{},
					&rules.IsDomainReaderRule{},
					&rules.IsRecordAdminRule{},
					&rules.IsDomainAdminRule{},
					&rules.IsAdminRule{},
					&rules.IsSuperadminRule{},
				},
			},
		},
	}
}
