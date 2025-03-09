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

type BasicAuthRequestAuthenticator struct {
	loginService interfaces.LoginServiceInterface
	userMapper   *UserMapper
}

func (a *BasicAuthRequestAuthenticator) Authenticate(context *gin.Context) (*model.User, error) {
	request := context.Request
	username, password, ok := request.BasicAuth()

	if !ok {
		return nil, fmt.Errorf("missing basic auth header")
	}

	dbUser, err := a.loginService.LoginByCredentials(username, password)
	if err != nil {
		return nil, err
	}

	return a.userMapper.FromDatabase(dbUser)
}

func ProvideBasicAuthenticator(container *basics.InjectionContainer) (interfaces.RequestAuthenticatorInterface, error) {
	if container.LoginService == nil {
		return nil, basics.NewMissingDependencyError("could not provide basic auth authenticator: login service could not be resolved")
	}

	return &BasicAuthRequestAuthenticator{loginService: container.LoginService, userMapper: &UserMapper{}}, nil
}
