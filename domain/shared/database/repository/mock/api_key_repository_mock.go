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

type APIKeyRepositoryMock struct {
	keys    map[string]*model.APIKey
	counter uint
}

func (r *APIKeyRepositoryMock) FindAPIKeysByUser(id uint) ([]*model.APIKey, error) {
	var apiKeys []*model.APIKey

	for _, apiKey := range r.keys {
		if apiKey.User.ID == id {
			apiKeys = append(apiKeys, apiKey)
		}
	}

	return apiKeys, nil
}

func (r *APIKeyRepositoryMock) SaveNewAPIKey(apiKey *model.APIKey) error {
	if r.CheckExistenceByAPIKey(apiKey.APIKey) {
		return errors2.NewItemAlreadyExistsError("api key already exists")
	}

	apiKey.ID = r.counter
	r.keys[apiKey.APIKey] = apiKey
	r.counter++

	return nil
}

func (r *APIKeyRepositoryMock) CheckExistenceByAPIKey(key string) bool {
	_, found := r.keys[key]

	return found
}

func (r *APIKeyRepositoryMock) FetchAPIKey(key string) (*model.APIKey, error) {
	apiKey, found := r.keys[key]
	if !found {
		return nil, errors2.NewItemNotFoundError("api key not found")
	}

	return apiKey, nil
}

func (r *APIKeyRepositoryMock) DeleteAPIKey(key string) error {
	if !r.CheckExistenceByAPIKey(key) {
		return errors2.NewItemNotFoundError("api key not found")
	}

	delete(r.keys, key)

	return nil
}

func ProvideAPIKeyRepositoryMock() (*APIKeyRepositoryMock, error) {
	return &APIKeyRepositoryMock{
		keys:    make(map[string]*model.APIKey),
		counter: 1,
	}, nil
}
