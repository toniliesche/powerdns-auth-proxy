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

package rdbms

import (
	"gorm.io/gorm"
	"powerdns-auth-proxy/domain/shared/basics"
	"powerdns-auth-proxy/domain/shared/database/model"
)

type UserDomainRoleRepository struct {
	database *gorm.DB
}

func (r *UserDomainRoleRepository) SaveNewUserDomainRole(role *model.UserDomainRole) error {
	return r.database.Create(role).Error
}

func (r *UserDomainRoleRepository) DeleteUserDomainRole(role *model.UserDomainRole) error {
	return r.database.Where("user_id = ? AND domain_id = ? AND domain_role_id = ?", role.UserId, role.DomainId, role.DomainRoleId).Delete(&model.UserDomainRole{}).Error
}

func (r *UserDomainRoleRepository) FindDomainRolesForUserAndDomain(user *model.User, domain *model.Domain) ([]*model.UserDomainRole, error) {
	var userDomainRoles []*model.UserDomainRole
	result := r.database.
		Preload("Domain").
		Preload("DomainRole").
		Where("user_id = ? AND domain_id = ?", user.ID, domain.ID).
		Model(userDomainRoles).
		Find(&userDomainRoles)

	if result.Error != nil {
		return nil, result.Error
	}

	return userDomainRoles, nil
}

func (r *UserDomainRoleRepository) FindDomainRolesForUser(user *model.User) ([]*model.UserDomainRole, error) {
	var userDomainRoles []*model.UserDomainRole
	result := r.database.
		Preload("Domain").
		Preload("DomainRole").
		Where("user_id = ?", user.ID).
		Model(userDomainRoles).
		Find(&userDomainRoles)

	if result.Error != nil {
		return nil, result.Error
	}

	return userDomainRoles, nil
}

func NewUserDomainRoleRepository(container *basics.InjectionContainer) (*UserDomainRoleRepository, error) {
	if container == nil {
		return nil, basics.NewMissingDependencyError("could not provide user domain role repository: passed injection container is nil")
	}

	if container.DB == nil {
		return nil, basics.NewMissingDependencyError("could not provide user domain role repository: database client could not be resolved")
	}

	return &UserDomainRoleRepository{database: container.DB}, nil
}
