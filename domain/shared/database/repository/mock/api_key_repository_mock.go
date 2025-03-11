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

type ApiKeyRepositoryMock struct {
	keys    map[string]*model.ApiKey
	counter uint
}

func (r *ApiKeyRepositoryMock) FindApiKeysByUserId(userId uint) ([]*model.ApiKey, error) {
	var apiKeys []*model.ApiKey

	for _, apiKey := range r.keys {
		if apiKey.User.ID == userId {
			apiKeys = append(apiKeys, apiKey)
		}
	}

	return apiKeys, nil
}

func (r *ApiKeyRepositoryMock) SaveNewApiKey(apiKey *model.ApiKey) error {
	if r.CheckExistenceByApiKey(apiKey.ApiKey) {
		return errors2.NewItemAlreadyExistsError("api key already exists")
	}

	apiKey.ID = r.counter
	r.keys[apiKey.ApiKey] = apiKey
	r.counter++

	return nil
}

func (r *ApiKeyRepositoryMock) CheckExistenceByApiKey(key string) bool {
	_, found := r.keys[key]

	return found
}

func (r *ApiKeyRepositoryMock) FetchApiKey(key string) (*model.ApiKey, error) {
	apiKey, found := r.keys[key]
	if !found {
		return nil, errors2.NewItemNotFoundError("api key not found")
	}

	return apiKey, nil
}

func (r *ApiKeyRepositoryMock) FetchApiKeyByIdentifier(identifier string) (*model.ApiKey, error) {
	for _, apiKey := range r.keys {
		if apiKey.Identifier == identifier {
			return apiKey, nil
		}
	}

	return nil, errors2.NewItemNotFoundError("api key not found")
}

func (r *ApiKeyRepositoryMock) DeleteApiKey(key string) error {
	if !r.CheckExistenceByApiKey(key) {
		return errors2.NewItemNotFoundError("api key not found")
	}

	delete(r.keys, key)

	return nil
}

func (r *ApiKeyRepositoryMock) DeleteApiKeyByIdentifier(identifier string) error {
	for _, apiKey := range r.keys {
		if apiKey.Identifier == identifier {
			delete(r.keys, apiKey.ApiKey)

			return nil
		}
	}

	return errors2.NewItemNotFoundError("api key not found")
}

func ProvideApiKeyRepositoryMock() (*ApiKeyRepositoryMock, error) {
	return &ApiKeyRepositoryMock{
		keys:    make(map[string]*model.ApiKey),
		counter: 1,
	}, nil
}
