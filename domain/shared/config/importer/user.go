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

package importer

import "fmt"

type User struct {
	Username        string            `yaml:"username"`
	Password        string            `yaml:"password"`
	Deleted         bool              `yaml:"deleted"`
	APIKey          bool              `yaml:"api_key"`
	UserDomainRoles []*UserDomainRole `yaml:"user_domain_roles"`
	UserRoles       []*UserRole       `yaml:"user_roles"`
}

func (u *User) Validate() error {
	var err error
	if u.Username == "" {
		return fmt.Errorf("username is required")
	}

	if u.Password == "" {
		return fmt.Errorf("password is required")
	}

	if u.UserDomainRoles != nil {
		for _, domainRole := range u.UserDomainRoles {
			if err = domainRole.Validate(); err != nil {
				return fmt.Errorf("invalid domain role config: %w", err)
			}
		}
	}

	if u.UserRoles != nil {
		for _, role := range u.UserRoles {
			if err = role.Validate(); err != nil {
				return fmt.Errorf("invalid role config: %w", err)
			}
		}
	}

	return nil
}

func (u *User) HasGlobalRoles() bool {
	return len(u.UserRoles) > 0
}

func (u *User) HasDomainRoles() bool {
	return (u.UserDomainRoles != nil) && (len(u.UserDomainRoles) > 0)
}
