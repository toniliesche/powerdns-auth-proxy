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

package mappers

import (
	"powerdns-auth-proxy/domain/shared/database/model"
	"powerdns-auth-proxy/domain/shared/model/management"
)

type SessionMapper struct {
}

func (m *SessionMapper) MapDatabaseToMinimalDto(session *model.TokenSession) *management.SessionMinimal {
	return &management.SessionMinimal{
		ID:             session.ID,
		SessionId:      session.SessionId,
		SequenceNumber: session.SequenceNumber,
		ExpiresAt:      session.ExpiresAt.Format("2006-01-02 15:04:05"),
	}
}

func (m *SessionMapper) MapDatabaseToMinimalDtoList(sessions []*model.TokenSession) []*management.SessionMinimal {
	dtos := make([]*management.SessionMinimal, 0)

	for _, session := range sessions {
		dtos = append(dtos, m.MapDatabaseToMinimalDto(session))
	}

	return dtos
}
