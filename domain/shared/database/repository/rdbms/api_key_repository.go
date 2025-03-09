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
	"powerdns-auth-proxy/domain/shared/database/repository/errors"
)

type APIKeyRepository struct {
	database *gorm.DB
}

func (r *APIKeyRepository) FindAPIKeysByUser(id uint) ([]*model.APIKey, error) {
	var apiKeys []*model.APIKey

	result := r.database.
		Preload("User").
		Where("user_id = ?", id).
		Find(&apiKeys)

	if result.Error != nil {
		return nil, result.Error
	}

	return apiKeys, nil
}

func (r *APIKeyRepository) SaveNewAPIKey(apiKey *model.APIKey) error {
	return r.database.Create(apiKey).Error
}

func (r *APIKeyRepository) CheckExistenceByAPIKey(key string) bool {
	var count int64
	r.database.Model(&model.APIKey{}).Where("api_key = ?", key).Count(&count)

	return count > 0
}

func (r *APIKeyRepository) FetchAPIKey(key string) (*model.APIKey, error) {
	var apiKey model.APIKey

	result := r.database.
		Model(&apiKey).
		Where("api_key = ?", key).
		Scan(&apiKey)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.NewItemNotFoundError("api key not found")
	}

	return &apiKey, nil
}

func (r *APIKeyRepository) DeleteAPIKey(key string) error {
	return r.database.Delete(&model.APIKey{}, "api_key = ?", key).Error
}

func ProvideAPIKeyRepository(container *basics.InjectionContainer) (*APIKeyRepository, error) {
	if container.DB == nil {
		return nil, basics.NewMissingDependencyError("could not provide api key repository: database client could not be resolved")
	}

	return &APIKeyRepository{database: container.DB}, nil
}
