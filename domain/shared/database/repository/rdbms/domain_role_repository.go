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

type DomainRoleRepository struct {
	database *gorm.DB
}

func (r *DomainRoleRepository) SaveNewDomainRole(role *model.DomainRole) error {
	return r.database.Create(role).Error
}

func (r *DomainRoleRepository) CheckExistenceByName(name string) bool {
	var count int64
	r.database.Model(&model.DomainRole{}).Where("name = ?", name).Count(&count)

	return count > 0
}

func (r *DomainRoleRepository) FetchDomainRoleByName(name string) (*model.DomainRole, error) {
	var domainRole model.DomainRole

	result := r.database.
		Model(&domainRole).
		Where("name = ?", name).
		Scan(&domainRole)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.NewItemNotFoundError("domain name not found")
	}

	return &domainRole, nil
}

func (r *DomainRoleRepository) DeleteDomainRoleByName(name string) error {
	role, err := r.FetchDomainRoleByName(name)
	if err != nil {
		return err
	}

	newRoleName := fmt.Sprintf("%s#deleted-%d", role.Name, role.ID)
	err = r.database.Model(&model.DomainRole{}).Where("id = ?", role.ID).Update("name", newRoleName).Error
	if err != nil {
		return err
	}

	return r.database.Where("id = ?", role.ID).Delete(&model.DomainRole{}).Error
}

func NewDomainRoleRepository(container *basics.InjectionContainer) (*DomainRoleRepository, error) {
	if container == nil {
		return nil, basics.NewMissingDependencyError("could not provide domain role repository: passed injection container is nil")
	}

	if container.DB == nil {
		return nil, basics.NewMissingDependencyError("could not provide domain role repository: database client could not be resolved")
	}

	return &DomainRoleRepository{database: container.DB}, nil
}
