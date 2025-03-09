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

package test

import (
	"fmt"
	"powerdns-auth-proxy/domain/shared/basics"
	"powerdns-auth-proxy/domain/shared/database/model"
)

func RepositoryCreateTestDomain(container *basics.InjectionContainer, registry *Registry) error {
	dbDomain := &model.Domain{FQDN: DomainFQDN}
	err := container.DomainRepository.SaveNewDomain(dbDomain)
	if err != nil {
		return err
	}

	registry.Set("domainId", fmt.Sprintf("%d", dbDomain.ID))
	registry.Set("domainFQDN", dbDomain.FQDN)

	return nil
}

func RepositoryCreateTestDomainRole(container *basics.InjectionContainer, registry *Registry) error {
	dbDomainRole := &model.DomainRole{Name: DomainRoleName}
	err := container.DomainRoleRepository.SaveNewDomainRole(dbDomainRole)
	if err != nil {
		return err
	}

	registry.Set("domainRoleId", fmt.Sprintf("%d", dbDomainRole.ID))
	registry.Set("domainRoleName", dbDomainRole.Name)

	return nil

}

func RepositoryCreateTestRole(container *basics.InjectionContainer, registry *Registry) error {
	dbRole := &model.Role{Name: RoleName}
	err := container.RoleRepository.SaveNewRole(dbRole)
	if err != nil {
		return err
	}

	registry.Set("roleId", fmt.Sprintf("%d", dbRole.ID))
	registry.Set("roleName", dbRole.Name)

	return nil
}

func RepositoryCreateTestUser(container *basics.InjectionContainer, registry *Registry) error {
	dbUser := &model.User{Username: UserUsername, Password: UserPassword}
	err := container.UserRepository.SaveNewUser(dbUser)
	if err != nil {
		return err
	}

	registry.Set("userId", fmt.Sprintf("%d", dbUser.ID))
	registry.Set("username", dbUser.Username)
	registry.Set("userPassword", dbUser.Password)

	return nil
}

func RepositoryCreateTestSession(container *basics.InjectionContainer, registry *Registry) error {
	dbTokenSession := &model.TokenSession{
		UserID:    registry.GetUint("userId"),
		SessionID: TokenSession,
	}

	err := container.SessionRepository.SaveNewTokenSession(dbTokenSession)
	if err != nil {
		return err
	}

	registry.Set("tokenSessionId", fmt.Sprintf("%d", dbTokenSession.ID))
	registry.Set("tokenSession", dbTokenSession.SessionID)

	return nil
}

func RepositoryCreateTestAPIKey(container *basics.InjectionContainer, registry *Registry) error {
	dbAPIKey := &model.APIKey{
		APIKey: APIKey,
		UserID: registry.GetUint("userId"),
	}

	err := container.APIKeyRepository.SaveNewAPIKey(dbAPIKey)
	if err != nil {
		return err
	}

	registry.Set("apiKey", dbAPIKey.APIKey)
	registry.Set("apiKeyId", fmt.Sprintf("%d", dbAPIKey.ID))

	return nil
}

func RepositoryCreateTestTokenSession(container *basics.InjectionContainer, registry *Registry) error {
	tokenSession := &model.TokenSession{
		UserID:    registry.GetUint("userId"),
		SessionID: TokenSession,
	}

	err := container.SessionRepository.SaveNewTokenSession(tokenSession)
	if err != nil {
		return err
	}

	registry.Set("tokenSession", tokenSession.SessionID)
	registry.Set("tokenSessionId", fmt.Sprintf("%d", tokenSession.ID))

	return nil
}

func RepositoryCreateTestUserDomainRole(container *basics.InjectionContainer, registry *Registry) error {
	dbUserDomainRole := &model.UserDomainRole{
		UserID:       registry.GetUint("userId"),
		DomainID:     registry.GetUint("domainId"),
		Domain:       &model.Domain{FQDN: "example.com"},
		DomainRoleID: registry.GetUint("domainRoleId"),
	}

	err := container.UserDomainRoleRepository.SaveNewUserDomainRole(dbUserDomainRole)
	if err != nil {
		return err
	}

	return nil
}

func RepositoryCreateTestUserRole(container *basics.InjectionContainer, registry *Registry) error {
	dbUserRole := &model.UserRole{
		UserID: registry.GetUint("userId"),
		RoleID: registry.GetUint("roleId"),
	}

	err := container.UserRoleRepository.SaveNewUserRole(dbUserRole)
	if err != nil {
		return err
	}

	return nil
}
