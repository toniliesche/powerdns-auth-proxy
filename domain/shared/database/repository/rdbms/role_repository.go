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
	"fmt"
	"gorm.io/gorm"
	"powerdns-auth-proxy/domain/shared/basics"
	"powerdns-auth-proxy/domain/shared/database/model"
	"powerdns-auth-proxy/domain/shared/database/repository/errors"
)

type RoleRepository struct {
	database *gorm.DB
}

func (r *RoleRepository) FetchAllRoles() ([]*model.Role, error) {
	var roles []*model.Role

	result := r.database.Find(&roles)
	if result.Error != nil {
		return nil, result.Error
	}

	return roles, nil
}

func (r *RoleRepository) SaveNewRole(role *model.Role) error {
	return r.database.Create(role).Error
}

func (r *RoleRepository) CheckExistenceByName(name string) bool {
	var count int64
	r.database.Model(&model.Role{}).Where("name = ?", name).Count(&count)

	return count > 0
}

func (r *RoleRepository) FetchRoleByName(name string) (*model.Role, error) {
	var role model.Role

	result := r.database.
		Model(&role).
		Where("name = ?", name).
		Scan(&role)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.NewItemNotFoundError("role not found")
	}

	return &role, nil
}

func (r *RoleRepository) DeleteRoleByName(roleName string) error {
	role, err := r.FetchRoleByName(roleName)
	if err != nil {
		return err
	}

	newRoleName := fmt.Sprintf("%s#deleted-%d", role.Name, role.ID)
	err = r.database.Model(&model.Role{}).Where("id = ?", role.ID).Update("name", newRoleName).Error
	if err != nil {
		return err
	}

	return r.database.Where("id = ?", role.ID).Delete(&model.Role{}).Error
}

func ProvideRoleRepository(container *basics.InjectionContainer) (*RoleRepository, error) {
	if container.DB == nil {
		return nil, basics.NewMissingDependencyError("could not provide role repository: database client could not be resolved")
	}

	return &RoleRepository{database: container.DB}, nil
}
