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

package database

import (
	"powerdns-auth-proxy/domain/shared/basics"
	"powerdns-auth-proxy/domain/shared/database/model"
	interfaces2 "powerdns-auth-proxy/domain/shared/database/repository/interfaces"
	"powerdns-auth-proxy/domain/shared/mappers"
	"powerdns-auth-proxy/domain/shared/model/management"
)

type UserDomainRoleService struct {
	domainRepository         interfaces2.DomainRepositoryInterface
	domainRoleRepository     interfaces2.DomainRoleRepositoryInterface
	userRepository           interfaces2.UserRepositoryInterface
	userDomainRoleRepository interfaces2.UserDomainRoleRepositoryInterface
	mapper                   *mappers.DomainRoleMapper
}

func (s *UserDomainRoleService) ListDomainRolesForUserByUserId(userId uint) ([]*management.DomainRole, error) {
	user, err := s.userRepository.FetchUserById(userId)
	if err != nil {
		return nil, err
	}

	dbDomainRoles, err := s.userDomainRoleRepository.FindDomainRolesForUser(user)
	if err != nil {
		return nil, err
	}

	return s.mapper.MapDatabaseToDtoList(dbDomainRoles), nil
}

func (s *UserDomainRoleService) ListDomainRolesForUser(username string) ([]*management.DomainRole, error) {
	user, err := s.userRepository.FetchUser(username)
	if err != nil {
		return nil, err
	}

	dbDomainRoles, err := s.userDomainRoleRepository.FindDomainRolesForUser(user)
	if err != nil {
		return nil, err
	}

	return s.mapper.MapDatabaseToDtoList(dbDomainRoles), nil
}

func (s *UserDomainRoleService) ListDomainRolesForUserByUserIdPerDomain(fqdn string, userId uint) ([]*management.DomainRole, error) {
	domain, err := s.domainRepository.FetchDomainByFqdn(fqdn)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepository.FetchUserById(userId)
	if err != nil {
		return nil, err
	}

	dbDomainRoles, err := s.userDomainRoleRepository.FindDomainRolesForUserAndDomain(user, domain)
	if err != nil {
		return nil, err
	}

	return s.mapper.MapDatabaseToDtoList(dbDomainRoles), nil
}

func (s *UserDomainRoleService) ListDomainRolesForUserPerDomain(fqdn string, username string) ([]*management.DomainRole, error) {
	domain, err := s.domainRepository.FetchDomainByFqdn(fqdn)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepository.FetchUser(username)
	if err != nil {
		return nil, err
	}

	dbDomainRoles, err := s.userDomainRoleRepository.FindDomainRolesForUserAndDomain(user, domain)
	if err != nil {
		return nil, err
	}

	return s.mapper.MapDatabaseToDtoList(dbDomainRoles), nil
}

func (s *UserDomainRoleService) GrantDomainRolesToUser(userId uint, roles []*management.DomainRole) error {
	user, err := s.userRepository.FetchUserById(userId)
	if err != nil {
		return err
	}

	for _, domainRole := range roles {
		err = s.GrantDomainRoleToUser(domainRole.Domain, user.Username, domainRole.Role)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *UserDomainRoleService) GrantDomainRoleToUser(fqdn string, username string, roleName string) error {
	domain, err := s.domainRepository.FetchDomainByFqdn(fqdn)
	if err != nil {
		return err
	}

	user, err := s.userRepository.FetchUser(username)
	if err != nil {
		return err
	}

	domainRole, err := s.domainRoleRepository.FetchDomainRoleByName(roleName)
	if err != nil {
		return err
	}

	userDomainRole := model.UserDomainRole{
		UserID:       user.ID,
		User:         user,
		DomainID:     domain.ID,
		Domain:       domain,
		DomainRoleID: domainRole.ID,
		DomainRole:   domainRole,
	}

	return s.userDomainRoleRepository.SaveNewUserDomainRole(&userDomainRole)
}

func (s *UserDomainRoleService) RevokeDomainRolesFromUser(userId uint, roles []*management.DomainRole) error {
	user, err := s.userRepository.FetchUserById(userId)
	if err != nil {
		return err
	}

	for _, domainRole := range roles {
		err = s.RevokeDomainRoleFromUser(domainRole.Domain, user.Username, domainRole.Role)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *UserDomainRoleService) RevokeDomainRoleFromUser(fqdn string, username string, roleName string) error {
	domain, err := s.domainRepository.FetchDomainByFqdn(fqdn)
	if err != nil {
		return err
	}

	user, err := s.userRepository.FetchUser(username)
	if err != nil {
		return err
	}

	domainRole, err := s.domainRoleRepository.FetchDomainRoleByName(roleName)
	if err != nil {
		return err
	}

	userDomainRole := model.UserDomainRole{
		UserID:       user.ID,
		User:         user,
		DomainID:     domain.ID,
		Domain:       domain,
		DomainRoleID: domainRole.ID,
		DomainRole:   domainRole,
	}

	return s.userDomainRoleRepository.DeleteUserDomainRole(&userDomainRole)
}

func ProvideUserDomainRoleService(container *basics.InjectionContainer) (*UserDomainRoleService, error) {
	if container.DomainRepository == nil {
		return nil, basics.NewMissingDependencyError("could not provide user domain role service: domain repository could not be resolved")
	}

	if container.DomainRoleRepository == nil {
		return nil, basics.NewMissingDependencyError("could not provide user domain role service: domain role repository could not be resolved")
	}

	if container.UserRepository == nil {
		return nil, basics.NewMissingDependencyError("could not provide user domain role service: user repository could not be resolved")
	}

	if container.UserDomainRoleRepository == nil {
		return nil, basics.NewMissingDependencyError("could not provide user domain role service: user domain role repository could not be resolved")
	}

	return &UserDomainRoleService{
		domainRepository:         container.DomainRepository,
		domainRoleRepository:     container.DomainRoleRepository,
		userRepository:           container.UserRepository,
		userDomainRoleRepository: container.UserDomainRoleRepository,
		mapper:                   &mappers.DomainRoleMapper{},
	}, nil
}
