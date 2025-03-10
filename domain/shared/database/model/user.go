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

package model

import (
	"fmt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username        string `gorm:"unique"`
	Password        string
	ApiKeys         []*ApiKey
	UserRoles       []*UserRole
	UserDomainRoles []*UserDomainRole
}

func (u *User) HasGlobalRole(roleName string) bool {
	for _, role := range u.UserRoles {
		if role.Role.Name == roleName {
			return true
		}
	}

	return false
}

func (u *User) GetGlobalRole(roleName string) (*UserRole, error) {
	for _, role := range u.UserRoles {
		if role.Role.Name == roleName {
			return role, nil
		}
	}

	return nil, fmt.Errorf("role not found: %s", roleName)
}

func (u *User) HasDomainRole(domain string, role string) bool {
	for _, domainRole := range u.UserDomainRoles {
		if domainRole.Domain.Fqdn == domain && domainRole.DomainRole.Name == role {
			return true
		}
	}

	return false
}

func (u *User) GetDomainRole(domain string, role string) (*UserDomainRole, error) {
	for _, domainRole := range u.UserDomainRoles {
		if domainRole.Domain.Fqdn == domain && domainRole.DomainRole.Name == role {
			return domainRole, nil
		}
	}

	return nil, fmt.Errorf("domain role not found: %s (%s)", role, domain)
}

func (u *User) IsAdminUser() bool {
	if u.HasGlobalRole("superadmin") {
		return true
	}

	return u.HasGlobalRole("admin")
}
