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
	"powerdns-auth-proxy/domain/shared/basics"
	"powerdns-auth-proxy/domain/shared/database/model"
	"powerdns-auth-proxy/domain/shared/database/repository/interfaces"
	errors2 "powerdns-auth-proxy/domain/shared/http/errors"
	"powerdns-auth-proxy/domain/shared/security"
)

type LoginService struct {
	userRepository interfaces.UserRepositoryInterface
}

func (s *LoginService) LoginByCredentials(username string, password string) (*model.User, errors2.HTTPError) {
	var dbUser *model.User
	var err error
	if dbUser, err = s.userRepository.FetchUser(username); err != nil {
		return nil, errors2.NewUnauthorizedError(err)
	}

	if !security.CheckPasswordHash(password, dbUser.Password) {
		return nil, errors2.NewUnauthorizedError(fmt.Errorf("invalid password"))
	}

	return dbUser, nil
}

func (s *LoginService) LoginByApiKey(apiKey string) (*model.User, error) {
	return s.userRepository.FetchUserByApiKey(apiKey)
}

func (s *LoginService) LoginByUsername(username string) (*model.User, error) {
	return s.userRepository.FetchUser(username)
}

func NewLoginService(container *basics.InjectionContainer) (*LoginService, error) {
	if container == nil {
		return nil, basics.NewMissingDependencyError("could not provide login service: passed injection container is nil")
	}
	if container.UserRepository == nil {
		return nil, basics.NewMissingDependencyError("could not provide login service: user repository could not be resolved")
	}

	return &LoginService{userRepository: container.UserRepository}, nil
}
