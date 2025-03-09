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

package services

import (
	"powerdns-auth-proxy/domain/shared/basics"
	"powerdns-auth-proxy/domain/shared/database/model"
	interfaces2 "powerdns-auth-proxy/domain/shared/database/repository/interfaces"
	"powerdns-auth-proxy/domain/shared/security"
)

type APIKeyService struct {
	apiKeyRepository interfaces2.APIKeyRepositoryInterface
	userRepository   interfaces2.UserRepositoryInterface
}

func (s *APIKeyService) Create(username string) (*model.APIKey, error) {
	user, err := s.userRepository.FetchUser(username)
	if err != nil {
		return nil, err
	}

	apiKey := security.GenerateRandomString(32, security.CharsetAlphaNumeric)
	hash, err := security.HashPassword(apiKey)
	if err != nil {
		return nil, err
	}

	dbApiKey := &model.APIKey{
		Identifier: security.GenerateRandomString(12, security.CharsetAlphaNumeric),
		APIKey:     hash,
		UserID:     user.ID,
		User:       user,
	}

	if err := s.apiKeyRepository.SaveNewAPIKey(dbApiKey); err != nil {
		return nil, err
	}

	return dbApiKey, nil
}

func (s *APIKeyService) Delete(apiKeyId string) error {
	return s.apiKeyRepository.DeleteAPIKey(apiKeyId)
}

func (s *APIKeyService) List(username string) ([]*model.APIKey, error) {
	user, err := s.userRepository.FetchUser(username)
	if err != nil {
		return nil, err
	}

	return s.apiKeyRepository.FindAPIKeysByUser(user.ID)
}

func ProvideAPIKeyService(container *basics.InjectionContainer) (*APIKeyService, error) {
	if container.APIKeyRepository == nil {
		return nil, basics.NewMissingDependencyError("could not provide api key service: api key repository could not be resolved")
	}

	if container.UserRepository == nil {
		return nil, basics.NewMissingDependencyError("could not provide api key service: user repository could not be resolved")
	}

	return &APIKeyService{apiKeyRepository: container.APIKeyRepository, userRepository: container.UserRepository}, nil
}
