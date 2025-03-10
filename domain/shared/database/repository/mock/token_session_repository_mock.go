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

type TokenSessionRepositoryMock struct {
	tokenSessions map[string]*model.TokenSession
	counter       uint
}

func (r *TokenSessionRepositoryMock) SaveNewTokenSession(tokenSession *model.TokenSession) error {
	if _, found := r.tokenSessions[tokenSession.SessionID]; found {
		return errors2.NewItemAlreadyExistsError("Token session already exists")
	}

	tokenSession.ID = r.counter
	r.tokenSessions[tokenSession.SessionID] = tokenSession
	r.counter++

	return nil
}

func (r *TokenSessionRepositoryMock) FindTokenSessionsByUserId(userId uint) ([]*model.TokenSession, error) {
	tokenSessions := make([]*model.TokenSession, 0)

	for _, tokenSession := range r.tokenSessions {
		if tokenSession.UserID == userId {
			tokenSessions = append(tokenSessions, tokenSession)
		}
	}

	return tokenSessions, nil
}

func (r *TokenSessionRepositoryMock) FetchTokenSession(sessionID string) (*model.TokenSession, error) {
	if tokenSession, found := r.tokenSessions[sessionID]; found {
		return tokenSession, nil
	}

	return nil, errors2.NewItemNotFoundError("Token session not found")
}

func (r *TokenSessionRepositoryMock) DeleteTokenSession(sessionID string) error {
	if _, found := r.tokenSessions[sessionID]; !found {
		return errors2.NewItemNotFoundError("Token session not found")
	}

	delete(r.tokenSessions, sessionID)

	return nil
}

func (r *TokenSessionRepositoryMock) UpdateTokenSession(session *model.TokenSession) error {
	return nil
}

func ProvideTokenSessionRepositoryMock() (*TokenSessionRepositoryMock, error) {
	return &TokenSessionRepositoryMock{
		tokenSessions: make(map[string]*model.TokenSession),
		counter:       1,
	}, nil
}
