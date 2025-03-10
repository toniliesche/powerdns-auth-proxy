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

type ApiKeyService struct {
	apiKeyRepository interfaces2.ApiKeyRepositoryInterface
	userRepository   interfaces2.UserRepositoryInterface
}

func (s *ApiKeyService) Create(username string) (*model.ApiKey, error) {
	user, err := s.userRepository.FetchUser(username)
	if err != nil {
		return nil, err
	}

	apiKey := security.GenerateRandomString(32, security.CharsetAlphaNumeric)
	hash, err := security.HashPassword(apiKey)
	if err != nil {
		return nil, err
	}

	dbApiKey := &model.ApiKey{
		Identifier: security.GenerateRandomString(12, security.CharsetAlphaNumeric),
		ApiKey:     hash,
		UserID:     user.ID,
		User:       user,
	}

	if err := s.apiKeyRepository.SaveNewApiKey(dbApiKey); err != nil {
		return nil, err
	}

	return dbApiKey, nil
}

func (s *ApiKeyService) Delete(apiKeyId string) error {
	return s.apiKeyRepository.DeleteApiKey(apiKeyId)
}

func (s *ApiKeyService) List(username string) ([]*model.ApiKey, error) {
	user, err := s.userRepository.FetchUser(username)
	if err != nil {
		return nil, err
	}

	return s.apiKeyRepository.FindApiKeysByUserId(user.ID)
}

func ProvideApiKeyService(container *basics.InjectionContainer) (*ApiKeyService, error) {
	if container.ApiKeyRepository == nil {
		return nil, basics.NewMissingDependencyError("could not provide api key service: api key repository could not be resolved")
	}

	if container.UserRepository == nil {
		return nil, basics.NewMissingDependencyError("could not provide api key service: user repository could not be resolved")
	}

	return &ApiKeyService{apiKeyRepository: container.ApiKeyRepository, userRepository: container.UserRepository}, nil
}
