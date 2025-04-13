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

package rdbms

import (
	"gorm.io/gorm"
	"powerdns-auth-proxy/domain/shared/basics"
	"powerdns-auth-proxy/domain/shared/database/model"
	"powerdns-auth-proxy/domain/shared/database/repository/errors"
)

type TokenSessionRepository struct {
	database *gorm.DB
}

func (r *TokenSessionRepository) SaveNewTokenSession(tokenSession *model.TokenSession) error {
	return r.database.Create(tokenSession).Error
}

func (r *TokenSessionRepository) FindTokenSessionsByUserId(userId uint) ([]*model.TokenSession, error) {
	var tokenSessions []*model.TokenSession

	result := r.database.
		Model(&model.TokenSession{}).
		Where("user_id = ?", userId).
		Find(&tokenSessions)

	if result.Error != nil {
		return nil, result.Error
	}

	return tokenSessions, nil
}

func (r *TokenSessionRepository) FetchTokenSession(sessionID string) (*model.TokenSession, error) {
	var tokenSession model.TokenSession
	result := r.database.
		Model(&tokenSession).
		Where("session_id = ?", sessionID).
		Scan(&tokenSession)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.NewItemNotFoundError("token session not found")
	}

	return &tokenSession, nil
}

func (r *TokenSessionRepository) DeleteTokenSession(sessionID string) error {
	return r.database.Where("session_id = ?", sessionID).Delete(&model.TokenSession{}).Error
}

func (r *TokenSessionRepository) UpdateTokenSession(session *model.TokenSession) error {
	return r.database.Save(session).Error
}

func NewTokenSessionRepository(container *basics.InjectionContainer) (*TokenSessionRepository, error) {
	if container == nil {
		return nil, basics.NewMissingDependencyError("could not provide token session repository: passed injection container is nil")
	}

	if container.DB == nil {
		return nil, basics.NewMissingDependencyError("could not provide token session repository: database client could not be resolved")
	}

	return &TokenSessionRepository{database: container.DB}, nil
}
