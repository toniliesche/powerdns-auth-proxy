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
	dbDomain := &model.Domain{Fqdn: DomainFqdn}
	err := container.DomainRepository.SaveNewDomain(dbDomain)
	if err != nil {
		return err
	}

	registry.Set("domainId", fmt.Sprintf("%d", dbDomain.ID))
	registry.Set("domainFqdn", dbDomain.Fqdn)

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

func RepositoryCreateTestTokenSession(container *basics.InjectionContainer, registry *Registry) error {
	if container.UserRepository == nil {
		return fmt.Errorf("could not create test user session: user repository could not be resolved")
	}

	if container.SessionRepository == nil {
		return fmt.Errorf("could not create test user session: session repository could not be resolved")
	}

	dbUser, err := container.UserRepository.FetchUserById(registry.GetUint("userId"))
	if err != nil {
		return fmt.Errorf("could not create test user session: %w", err)
	}

	dbTokenSession := &model.TokenSession{
		User:      dbUser,
		UserID:    dbUser.ID,
		SessionID: TokenSession,
	}

	err = container.SessionRepository.SaveNewTokenSession(dbTokenSession)
	if err != nil {
		return err
	}

	registry.Set("tokenSessionId", fmt.Sprintf("%d", dbTokenSession.ID))
	registry.Set("tokenSession", dbTokenSession.SessionID)

	return nil
}

func RepositoryCreateTestApiKey(container *basics.InjectionContainer, registry *Registry) error {
	if container.UserRepository == nil {
		return fmt.Errorf("could not create test api key: user repository could not be resolved")
	}

	if container.ApiKeyRepository == nil {
		return fmt.Errorf("could not create test api key: api key repository could not be resolved")
	}

	dbUser, err := container.UserRepository.FetchUserById(registry.GetUint("userId"))
	if err != nil {
		return fmt.Errorf("could not create test api key: %w", err)
	}

	dbApiKey := &model.ApiKey{
		ApiKey: ApiKey,
		UserID: dbUser.ID,
		User:   dbUser,
	}

	err = container.ApiKeyRepository.SaveNewApiKey(dbApiKey)
	if err != nil {
		return err
	}

	registry.Set("apiKey", dbApiKey.ApiKey)
	registry.Set("apiKeyId", fmt.Sprintf("%d", dbApiKey.ID))

	return nil
}

func RepositoryCreateTestUserDomainRole(container *basics.InjectionContainer, registry *Registry) error {
	if container.UserRepository == nil {
		return fmt.Errorf("could not create test user domain role: user repository could not be resolved")
	}

	if container.DomainRepository == nil {
		return fmt.Errorf("could not create test user domain role: domain repository could not be resolved")
	}

	if container.DomainRoleRepository == nil {
		return fmt.Errorf("could not create test user domain role: domain role repository could not be resolved")
	}

	if container.UserDomainRoleRepository == nil {
		return fmt.Errorf("could not create test user domain role: user domain role repository could not be resolved")
	}

	dbUser, err := container.UserRepository.FetchUserById(registry.GetUint("userId"))
	if err != nil {
		return fmt.Errorf("could not create test user domain role: %w", err)
	}

	dbDomain, err := container.DomainRepository.FetchDomainById(registry.GetUint("domainId"))
	if err != nil {
		return fmt.Errorf("could not create test user domain role: %w", err)
	}

	dbDomainRole, err := container.DomainRoleRepository.FetchDomainRoleByName(DomainRoleName)
	if err != nil {
		return fmt.Errorf("could not create test user domain role: %w", err)
	}

	dbUserDomainRole := &model.UserDomainRole{
		UserID:       dbUser.ID,
		User:         dbUser,
		DomainID:     dbDomain.ID,
		Domain:       dbDomain,
		DomainRoleID: dbDomainRole.ID,
		DomainRole:   dbDomainRole,
	}

	err = container.UserDomainRoleRepository.SaveNewUserDomainRole(dbUserDomainRole)
	if err != nil {
		return err
	}

	dbUser.UserDomainRoles = append(dbUser.UserDomainRoles, dbUserDomainRole)

	return nil
}

func RepositoryCreateTestUserRole(container *basics.InjectionContainer, registry *Registry) error {
	if container.UserRepository == nil {
		return fmt.Errorf("could not create test user role: user repository could not be resolved")
	}

	if container.RoleRepository == nil {
		return fmt.Errorf("could not create test user role: role repository could not be resolved")
	}

	if container.UserRoleRepository == nil {
		return fmt.Errorf("could not create test user role: user role repository could not be resolved")
	}

	dbUser, err := container.UserRepository.FetchUserById(registry.GetUint("userId"))
	if err != nil {
		return fmt.Errorf("could not create test user role: %w", err)
	}

	dbRole, err := container.RoleRepository.FetchRoleById(registry.GetUint("roleId"))
	if err != nil {
		return fmt.Errorf("could not create test user role: %w", err)
	}

	dbUserRole := &model.UserRole{
		UserID: dbUser.ID,
		User:   dbUser,
		RoleID: dbRole.ID,
		Role:   dbRole,
	}

	err = container.UserRoleRepository.SaveNewUserRole(dbUserRole)
	if err != nil {
		return err
	}

	dbUser.UserRoles = append(dbUser.UserRoles, dbUserRole)

	return nil
}
