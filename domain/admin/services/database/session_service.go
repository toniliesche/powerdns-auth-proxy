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

package database

import (
	"powerdns-auth-proxy/domain/shared/basics"
	interfaces2 "powerdns-auth-proxy/domain/shared/database/repository/interfaces"
	"powerdns-auth-proxy/domain/shared/mappers"
	"powerdns-auth-proxy/domain/shared/model/management"
)

type SessionService struct {
	sessionRepository interfaces2.TokenSessionRepositoryInterface
	userRepository    interfaces2.UserRepositoryInterface
	mapper            *mappers.SessionMapper
}

func (s SessionService) ListUserSessions(userId uint) ([]*management.SessionMinimal, error) {
	dbSessions, err := s.sessionRepository.FindTokenSessionsByUserId(userId)
	if err != nil {
		return nil, err
	}

	return s.mapper.MapDatabaseToMinimalDtoList(dbSessions), nil
}

func (s SessionService) LogoutUserSession(userSessions []string) error {
	for _, sessionID := range userSessions {
		err := s.sessionRepository.DeleteTokenSession(sessionID)
		if err != nil {
			return err
		}
	}

	return nil
}

func ProvideSessionService(container *basics.InjectionContainer) (*SessionService, error) {
	if container.SessionRepository == nil {
		return nil, basics.NewMissingDependencyError("could not provide session service: session repository could not be resolved")
	}

	if container.UserRepository == nil {
		return nil, basics.NewMissingDependencyError("could not provide session service: user repository could not be resolved")
	}

	return &SessionService{
		sessionRepository: container.SessionRepository,
		userRepository:    container.UserRepository,
		mapper:            &mappers.SessionMapper{},
	}, nil
}
