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

type DomainRoleRepositoryMock struct {
	roles   map[string]*model.DomainRole
	counter uint
}

func (r *DomainRoleRepositoryMock) SaveNewDomainRole(role *model.DomainRole) error {
	if r.CheckExistenceByName(role.Name) {
		return errors2.NewItemAlreadyExistsError("domain role already exists")
	}

	role.ID = r.counter
	r.roles[role.Name] = role
	r.counter++

	return nil
}

func (r *DomainRoleRepositoryMock) CheckExistenceByName(name string) bool {
	_, found := r.roles[name]

	return found
}

func (r *DomainRoleRepositoryMock) FetchDomainRoleByName(name string) (*model.DomainRole, error) {
	role, found := r.roles[name]
	if !found {
		return nil, errors2.NewItemNotFoundError("domain name not found")
	}

	return role, nil
}

func (r *DomainRoleRepositoryMock) DeleteDomainRoleByName(roleName string) error {
	if !r.CheckExistenceByName(roleName) {
		return errors2.NewItemNotFoundError("domain role not found")
	}

	delete(r.roles, roleName)

	return nil
}

func ProvideDomainRoleRepositoryMock() (*DomainRoleRepositoryMock, error) {
	return &DomainRoleRepositoryMock{
		roles:   make(map[string]*model.DomainRole),
		counter: 1,
	}, nil
}
