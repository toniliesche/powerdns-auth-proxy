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
	"fmt"
	"github.com/gin-gonic/gin"
	"powerdns-auth-proxy/domain/shared/auth/model"
	"powerdns-auth-proxy/domain/shared/basics"
	"powerdns-auth-proxy/domain/shared/interfaces"
)

type ApiKeyRequestAuthenticator struct {
	userMapper   *UserMapper
	loginService interfaces.LoginServiceInterface
}

func (a *ApiKeyRequestAuthenticator) Authenticate(context *gin.Context) (*model.User, error) {
	request := context.Request
	apiKey := request.Header.Get("X-API_KEY")

	if apiKey == "" {
		return nil, fmt.Errorf("missing api key header")
	}

	dbUser, err := a.loginService.LoginByApiKey(apiKey)
	if err != nil {
		return nil, err
	}

	return a.userMapper.FromDatabase(dbUser)
}

func ProvideApiKeyRequestAuthenticator(container *basics.InjectionContainer) (interfaces.RequestAuthenticatorInterface, error) {
	if container.LoginService == nil {
		return nil, basics.NewMissingDependencyError("could not provide api key request authenticator: login service could not be resolved")
	}

	return &ApiKeyRequestAuthenticator{loginService: container.LoginService, userMapper: &UserMapper{}}, nil
}
