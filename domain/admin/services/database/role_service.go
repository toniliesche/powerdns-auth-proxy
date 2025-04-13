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
	"powerdns-auth-proxy/domain/shared/database/repository/interfaces"
	"powerdns-auth-proxy/domain/shared/mappers"
	"powerdns-auth-proxy/domain/shared/model/management"
)

type RoleService struct {
	roleRepository interfaces.RoleRepositoryInterface
	mapper         *mappers.RoleMapper
}

func (s *RoleService) ListRoles() ([]*management.Role, error) {
	roles, err := s.roleRepository.FetchAllRoles()
	if err != nil {
		return nil, err
	}

	return s.mapper.MapDbToDtoList(roles), nil
}

func (s *RoleService) CreateRole(payload *management.RoleCreatePayload) (uint, error) {
	role := s.mapper.MapCreatePayloadToDb(payload)

	err := s.roleRepository.SaveNewRole(role)
	if err != nil {
		return 0, err
	}

	return role.ID, nil
}

func (s *RoleService) DeleteRole(name string) error {
	return s.roleRepository.DeleteRoleByName(name)
}

func NewRoleService(container *basics.InjectionContainer) (*RoleService, error) {
	if container == nil {
		return nil, basics.NewMissingDependencyError("could not provide role service: passed injection container is nil")
	}
	if container.RoleRepository == nil {
		return nil, basics.NewMissingDependencyError("could not provide role service: role repository could not be resolved")
	}

	return &RoleService{
		roleRepository: container.RoleRepository,
		mapper:         &mappers.RoleMapper{},
	}, nil
}
