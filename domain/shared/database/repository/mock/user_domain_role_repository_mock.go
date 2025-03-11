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
	"powerdns-auth-proxy/domain/shared/database/repository/interfaces"
)

type UserDomainRoleRepositoryMock struct {
	userDomainRoles  map[string]map[string]*model.UserDomainRole
	domainRepository interfaces.DomainRepositoryInterface
}

func (r *UserDomainRoleRepositoryMock) SaveNewUserDomainRole(role *model.UserDomainRole) error {
	domain, err := r.domainRepository.FetchDomainById(role.DomainId)
	if err != nil {
		return err
	}

	domainKey := domain.Fqdn
	if _, found := r.userDomainRoles[domainKey]; !found {
		r.userDomainRoles[domainKey] = make(map[string]*model.UserDomainRole)
	}

	key := fmt.Sprintf("%d-%d", role.UserId, role.DomainRoleId)
	if _, found := r.userDomainRoles[domainKey][key]; found {
		return errors2.NewItemAlreadyExistsError("User domain role already exists")
	}

	r.userDomainRoles[domainKey][key] = role

	return nil
}

func (r *UserDomainRoleRepositoryMock) DeleteUserDomainRole(role *model.UserDomainRole) error {
	domain, err := r.domainRepository.FetchDomainById(role.DomainId)
	if err != nil {
		return err
	}

	domainKey := domain.Fqdn
	if _, found := r.userDomainRoles[domainKey]; !found {
		return errors2.NewItemNotFoundError("User domain role not found")
	}

	key := fmt.Sprintf("%d-%d", role.UserId, role.DomainRoleId)
	if _, found := r.userDomainRoles[domainKey][key]; !found {
		return errors2.NewItemNotFoundError("User domain role not found")
	}

	delete(r.userDomainRoles[domainKey], key)

	return nil
}

func (r *UserDomainRoleRepositoryMock) FindDomainRolesForUserAndDomain(user *model.User, domain *model.Domain) ([]*model.UserDomainRole, error) {
	roles := make([]*model.UserDomainRole, 0)

	if _, found := r.userDomainRoles[domain.Fqdn]; !found {
		return roles, nil
	}

	for _, role := range r.userDomainRoles[domain.Fqdn] {
		if role.UserId == user.ID {
			roles = append(roles, role)
		}
	}

	return roles, nil
}

func (r *UserDomainRoleRepositoryMock) FindDomainRolesForUser(user *model.User) ([]*model.UserDomainRole, error) {
	roles := make([]*model.UserDomainRole, 0)

	for _, domainRoles := range r.userDomainRoles {
		for _, role := range domainRoles {
			if role.UserId == user.ID {
				roles = append(roles, role)
			}
		}
	}

	return roles, nil
}

func ProvideUserDomainRoleRepositoryMock(container *basics.InjectionContainer) (*UserDomainRoleRepositoryMock, error) {
	if container.DomainRepository == nil {
		return nil, fmt.Errorf("could not provide user domain role repository mock: domain repository could not be resolved")
	}

	return &UserDomainRoleRepositoryMock{
		domainRepository: container.DomainRepository,
		userDomainRoles:  make(map[string]map[string]*model.UserDomainRole),
	}, nil
}
