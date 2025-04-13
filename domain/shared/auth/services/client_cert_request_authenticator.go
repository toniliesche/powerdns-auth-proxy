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
	"github.com/gin-gonic/gin"
	"powerdns-auth-proxy/domain/shared/auth/model"
	"powerdns-auth-proxy/domain/shared/basics"
	"powerdns-auth-proxy/domain/shared/interfaces"
)

type ClientCertRequestAuthenticator struct {
	loginService interfaces.LoginServiceInterface
	userMapper   *UserMapper
}

func (a *ClientCertRequestAuthenticator) Authenticate(context *gin.Context) (*model.User, error) {
	dbUser, err := a.loginService.LoginByUsername("")
	if err != nil {
		return nil, err
	}

	return a.userMapper.FromDatabase(dbUser)
}

func NewClientCertRequestAuthenticator(container *basics.InjectionContainer) (*ClientCertRequestAuthenticator, error) {
	if container == nil {
		return nil, basics.NewMissingDependencyError("could not provide client cert request authenticator: passed injection container is nil")
	}

	if container.LoginService == nil {
		return nil, basics.NewMissingDependencyError("could not provide client cert request authenticator: login service could not be resolved")
	}

	return &ClientCertRequestAuthenticator{loginService: container.LoginService, userMapper: &UserMapper{}}, nil
}
