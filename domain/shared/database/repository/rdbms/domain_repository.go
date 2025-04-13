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

type DomainRepository struct {
	database *gorm.DB
}

func (r *DomainRepository) SaveNewDomain(domain *model.Domain) error {
	return r.database.Create(domain).Error
}

func (r *DomainRepository) CheckExistenceByFqdn(fqdn string) bool {
	var count int64
	r.database.Model(&model.Domain{}).Where("fqdn = ?", fqdn).Count(&count)

	return count > 0
}

func (r *DomainRepository) FetchDomainById(id uint) (*model.Domain, error) {
	var domainModel model.Domain

	result := r.database.
		Model(&domainModel).
		Where("id = ?", id).
		Scan(&domainModel)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.NewItemNotFoundError(fmt.Sprintf("domain with id %d not found", id))
	}

	return &domainModel, nil
}

func (r *DomainRepository) FetchDomainByFqdn(fqdn string) (*model.Domain, error) {
	var domainModel model.Domain

	result := r.database.
		Model(&domainModel).
		Where("fqdn = ?", fqdn).
		Scan(&domainModel)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.NewItemNotFoundError(fmt.Sprintf("domain with fqdn %s not found", fqdn))
	}

	return &domainModel, nil
}

func (r *DomainRepository) FindAll() ([]*model.Domain, error) {
	var domains []*model.Domain

	result := r.database.
		Model(&model.Domain{}).
		Find(&domains)

	if result.Error != nil {
		return nil, result.Error
	}

	return domains, nil
}

func (r *DomainRepository) DeleteDomain(fqdn string) error {
	domain, err := r.FetchDomainByFqdn(fqdn)
	if err != nil {
		return err
	}

	newFqdn := fmt.Sprintf("%s#deleted-%d", domain.Fqdn, domain.ID)
	err = r.database.Model(&model.Domain{}).Where("id = ?", domain.ID).Update("fqdn", newFqdn).Error
	if err != nil {
		return err
	}

	return r.database.Where("id = ?", domain.ID).Delete(&model.Domain{}).Error
}

func (r *DomainRepository) DeleteDomainById(id uint) error {
	domain, err := r.FetchDomainById(id)
	if err != nil {
		return err
	}

	newFqdn := fmt.Sprintf("%s#deleted-%d", domain.Fqdn, domain.ID)
	err = r.database.Model(&model.Domain{}).Where("id = ?", domain.ID).Update("fqdn", newFqdn).Error
	if err != nil {
		return err
	}

	return r.database.Where("id = ?", domain.ID).Delete(&model.Domain{}).Error
}

func NewDomainRepository(container *basics.InjectionContainer) (*DomainRepository, error) {
	if container == nil {
		return nil, basics.NewMissingDependencyError("could not provide domain repository: passed injection container is nil")
	}

	if container.DB == nil {
		return nil, basics.NewMissingDependencyError("could not provide domain repository: database client could not be resolved")
	}

	return &DomainRepository{database: container.DB}, nil
}
