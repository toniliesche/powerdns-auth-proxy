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

package mock

import (
	"fmt"
	"powerdns-auth-proxy/domain/shared/basics"
	"powerdns-auth-proxy/domain/shared/database/model"
	errors2 "powerdns-auth-proxy/domain/shared/database/repository/errors"
)

type UserRoleRepositoryMock struct {
	userRoles map[string]*model.UserRole
}

func (r *UserRoleRepositoryMock) SaveNewUserRole(role *model.UserRole) error {
	key := fmt.Sprintf("%d-%d", role.UserId, role.RoleId)
	if _, found := r.userRoles[key]; found {
		return errors2.NewItemAlreadyExistsError("User role already exists")
	}

	r.userRoles[key] = role

	return nil
}

func (r *UserRoleRepositoryMock) DeleteUserRole(role *model.UserRole) error {
	key := fmt.Sprintf("%d-%d", role.UserId, role.RoleId)
	if _, found := r.userRoles[key]; !found {
		return errors2.NewItemNotFoundError("User role not found")
	}

	delete(r.userRoles, key)

	return nil
}

func (r *UserRoleRepositoryMock) FindRolesForUser(user *model.User) ([]*model.UserRole, error) {
	roles := make([]*model.UserRole, 0)

	for _, role := range r.userRoles {
		if role.UserId == user.ID {
			roles = append(roles, role)
		}
	}

	return roles, nil
}

func NewUserRoleRepositoryMock(container *basics.InjectionContainer) (*UserRoleRepositoryMock, error) {
	if container == nil {
		return nil, basics.NewMissingDependencyError("could not provide user role repository mock: passed injection container is nil")
	}

	return &UserRoleRepositoryMock{
		userRoles: make(map[string]*model.UserRole),
	}, nil
}
