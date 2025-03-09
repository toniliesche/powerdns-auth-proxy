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

type TokenSessionRepositoryInterface interface {
	SaveNewTokenSession(tokenSession *model.TokenSession) error
	FindTokenSessionsByUser(id uint) ([]*model.TokenSession, error)
	FetchTokenSession(sessionID string) (*model.TokenSession, error)
	DeleteTokenSession(sessionID string) error
	UpdateTokenSession(session *model.TokenSession) error
}
