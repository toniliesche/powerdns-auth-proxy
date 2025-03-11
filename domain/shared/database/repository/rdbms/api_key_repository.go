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

type ApiKeyRepository struct {
	database *gorm.DB
}

func (r *ApiKeyRepository) FindApiKeysByUserId(userId uint) ([]*model.ApiKey, error) {
	var apiKeys []*model.ApiKey

	result := r.database.
		Preload("User").
		Where("user_id = ?", userId).
		Find(&apiKeys)

	if result.Error != nil {
		return nil, result.Error
	}

	return apiKeys, nil
}

func (r *ApiKeyRepository) SaveNewApiKey(apiKey *model.ApiKey) error {
	return r.database.Create(apiKey).Error
}

func (r *ApiKeyRepository) CheckExistenceByApiKey(key string) bool {
	var count int64
	r.database.Model(&model.ApiKey{}).Where("api_key = ?", key).Count(&count)

	return count > 0
}

func (r *ApiKeyRepository) FetchApiKey(key string) (*model.ApiKey, error) {
	var apiKey model.ApiKey

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

func (r *ApiKeyRepository) FetchApiKeyByIdentifier(identifier string) (*model.ApiKey, error) {
	var apiKey model.ApiKey

	result := r.database.
		Model(&apiKey).
		Where("identifier = ?", identifier).
		Scan(&apiKey)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.NewItemNotFoundError("api key not found")
	}

	return &apiKey, nil
}

func (r *ApiKeyRepository) DeleteApiKey(key string) error {
	apiKey, err := r.FetchApiKey(key)
	if err != nil {
		return err
	}

	newIdentifier := fmt.Sprintf("%s#deleted-%d", apiKey.Identifier, apiKey.ID)
	err = r.database.Model(&model.ApiKey{}).Where("id = ?", apiKey.ID).Update("identifier", newIdentifier).Error
	if err != nil {
		return err
	}

	return r.database.Where("id = ?", apiKey.ID).Delete(&model.ApiKey{}).Error
}

func (r *ApiKeyRepository) DeleteApiKeyByIdentifier(identifier string) error {
	apiKey, err := r.FetchApiKeyByIdentifier(identifier)
	if err != nil {
		return err
	}

	newIdentifier := fmt.Sprintf("%s#deleted-%d", apiKey.Identifier, apiKey.ID)
	err = r.database.Model(&model.ApiKey{}).Where("id = ?", apiKey.ID).Update("identifier", newIdentifier).Error
	if err != nil {
		return err
	}

	return r.database.Where("id = ?", apiKey.ID).Delete(&model.ApiKey{}).Error
}

func ProvideApiKeyRepository(container *basics.InjectionContainer) (*ApiKeyRepository, error) {
	if container.DB == nil {
		return nil, basics.NewMissingDependencyError("could not provide api key repository: database client could not be resolved")
	}

	return &ApiKeyRepository{database: container.DB}, nil
}
