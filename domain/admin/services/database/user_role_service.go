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
	"fmt"
	"powerdns-auth-proxy/domain/shared/basics"
	"powerdns-auth-proxy/domain/shared/database/model"
	interfaces2 "powerdns-auth-proxy/domain/shared/database/repository/interfaces"
	"powerdns-auth-proxy/domain/shared/mappers"
	"powerdns-auth-proxy/domain/shared/model/management"
)

type UserRoleService struct {
	roleRepository     interfaces2.RoleRepositoryInterface
	userRepository     interfaces2.UserRepositoryInterface
	userRoleRepository interfaces2.UserRoleRepositoryInterface
	mapper             *mappers.UserRoleMapper
}

func (s *UserRoleService) ListRolesForUser(username string) ([]*management.UserRoleMinimal, error) {
	user, err := s.userRepository.FetchUser(username)
	if err != nil {
		return nil, err
	}

	dbUserRoles, err := s.userRoleRepository.FindRolesForUser(user)
	if err != nil {
		return nil, err
	}

	return s.mapper.MapDatabaseToDtoList(dbUserRoles), nil
}

func (s *UserRoleService) ListRolesForUserByID(id uint) ([]*management.UserRoleMinimal, error) {
	user, err := s.userRepository.FetchUserByID(id)
	if err != nil {
		return nil, err
	}

	dbUserRoles, err := s.userRoleRepository.FindRolesForUser(user)
	if err != nil {
		return nil, err
	}

	return s.mapper.MapDatabaseToDtoList(dbUserRoles), nil
}

func (s *UserRoleService) GrantRoleToUser(username string, roleName string) error {
	user, err := s.userRepository.FetchUser(username)
	if err != nil {
		return err
	}

	role, err := s.roleRepository.FetchRoleByName(roleName)
	if err != nil {
		return err
	}

	userRole := model.UserRole{
		UserID: user.ID,
		User:   user,
		RoleID: role.ID,
		Role:   role,
	}

	return s.userRoleRepository.SaveNewUserRole(&userRole)
}

func (s *UserRoleService) RevokeRoleFromUser(username string, roleName string) error {
	user, err := s.userRepository.FetchUser(username)
	if err != nil {
		return err
	}

	role, err := s.roleRepository.FetchRoleByName(roleName)
	if err != nil {
		return err
	}

	userRole := model.UserRole{
		UserID: user.ID,
		RoleID: role.ID,
	}

	return s.userRoleRepository.DeleteUserRole(&userRole)
}

func (s *UserRoleService) GrantRolesToUser(userId uint, roles []string) error {
	user, err := s.userRepository.FetchUserByID(userId)
	if err != nil {
		return err
	}

	for _, roleName := range roles {
		role, err := s.roleRepository.FetchRoleByName(roleName)
		if err != nil {
			return err
		}

		userRole := model.UserRole{
			UserID: user.ID,
			User:   user,
			RoleID: role.ID,
			Role:   role,
		}

		err = s.userRoleRepository.SaveNewUserRole(&userRole)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *UserRoleService) RevokeRolesFromUser(userId uint, roles []string) error {
	user, err := s.userRepository.FetchUserByID(userId)
	if err != nil {
		return err
	}

	for _, roleName := range roles {
		role, err := s.roleRepository.FetchRoleByName(roleName)
		if err != nil {
			return err
		}

		userRole := model.UserRole{
			UserID: user.ID,
			RoleID: role.ID,
		}

		err = s.userRoleRepository.DeleteUserRole(&userRole)
		if err != nil {
			return err
		}
	}

	return nil
}

func ProvideUserRoleService(container *basics.InjectionContainer) (*UserRoleService, error) {
	if container.RoleRepository == nil {
		return nil, fmt.Errorf("could not provide user role service: role repository could not be resolved")
	}

	if container.UserRepository == nil {
		return nil, fmt.Errorf("could not provide user role service: user repository could not be resolved")
	}

	if container.UserRoleRepository == nil {
		return nil, fmt.Errorf("could not provide user role service: user role repository could not be resolved")
	}

	return &UserRoleService{
		roleRepository:     container.RoleRepository,
		userRepository:     container.UserRepository,
		userRoleRepository: container.UserRoleRepository,
		mapper:             &mappers.UserRoleMapper{},
	}, nil
}
