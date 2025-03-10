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
	"powerdns-auth-proxy/domain/shared/database/model"
	errors2 "powerdns-auth-proxy/domain/shared/database/repository/errors"
)

type RoleRepositoryMock struct {
	roles   map[string]*model.Role
	counter uint
}

func (r *RoleRepositoryMock) FetchAllRoles() ([]*model.Role, error) {
	roles := make([]*model.Role, 0, len(r.roles))

	for _, role := range r.roles {
		roles = append(roles, role)
	}

	return roles, nil
}

func (r *RoleRepositoryMock) SaveNewRole(role *model.Role) error {
	if r.CheckExistenceByName(role.Name) {
		return errors2.NewItemAlreadyExistsError("role already exists")
	}

	role.ID = r.counter
	r.roles[role.Name] = role
	r.counter++

	return nil
}

func (r *RoleRepositoryMock) CheckExistenceByName(name string) bool {
	_, found := r.roles[name]

	return found
}

func (r *RoleRepositoryMock) FetchRoleById(roleId uint) (*model.Role, error) {
	for _, role := range r.roles {
		if role.ID == roleId {
			return role, nil
		}
	}

	return nil, errors2.NewItemNotFoundError("role not found")
}

func (r *RoleRepositoryMock) FetchRoleByName(name string) (*model.Role, error) {
	role, found := r.roles[name]
	if !found {
		return nil, errors2.NewItemNotFoundError("role not found")
	}

	return role, nil
}

func (r *RoleRepositoryMock) DeleteRoleByName(name string) error {
	if !r.CheckExistenceByName(name) {
		return errors2.NewItemNotFoundError("role not found")
	}

	delete(r.roles, name)

	return nil
}

func ProvideRoleRepositoryMock() (*RoleRepositoryMock, error) {
	return &RoleRepositoryMock{
		roles:   make(map[string]*model.Role),
		counter: 1,
	}, nil
}
