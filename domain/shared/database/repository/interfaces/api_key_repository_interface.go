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

package interfaces

import "powerdns-auth-proxy/domain/shared/database/model"

type ApiKeyRepositoryInterface interface {
	SaveNewApiKey(apiKey *model.ApiKey) error
	CheckExistenceByApiKey(key string) bool
	FetchApiKey(key string) (*model.ApiKey, error)
	FetchApiKeyByIdentifier(identifier string) (*model.ApiKey, error)
	FindApiKeysByUserId(userId uint) ([]*model.ApiKey, error)
	DeleteApiKey(key string) error
}
