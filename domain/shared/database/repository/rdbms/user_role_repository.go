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

type UserRoleRepository struct {
	database *gorm.DB
}

func (r *UserRoleRepository) SaveNewUserRole(role *model.UserRole) error {
	return r.database.Create(role).Error
}

func (r *UserRoleRepository) DeleteUserRole(role *model.UserRole) error {
	return r.database.Where("user_id = ? AND role_id = ?", role.UserID, role.RoleID).Delete(&model.UserRole{}).Error
}

func (r *UserRoleRepository) FindRolesForUser(user *model.User) ([]*model.UserRole, error) {
	var userRoles []*model.UserRole
	result := r.database.
		Preload("Role").
		Where("user_id = ?", user.ID).
		Model(userRoles).
		Find(&userRoles)

	if result.Error != nil {
		return nil, result.Error
	}

	return userRoles, nil
}

func ProvideUserRoleRepository(container *basics.InjectionContainer) (*UserRoleRepository, error) {
	if container.DB == nil {
		return nil, basics.NewMissingDependencyError("could not provide user role repository: database client could not be resolved")
	}

	return &UserRoleRepository{database: container.DB}, nil
}
