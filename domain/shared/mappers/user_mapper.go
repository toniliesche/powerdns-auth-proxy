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

type UserMapper struct {
}

func (m *UserMapper) MapDatabaseToDto(user *model.User) *management.User {
	userRoles := make([]string, 0, len(user.UserRoles))
	for _, role := range user.UserRoles {
		userRoles = append(userRoles, role.Role.Name)
	}

	domainRoles := make(map[string][]string, len(user.UserDomainRoles))
	for _, domainRole := range user.UserDomainRoles {
		if _, ok := domainRoles[domainRole.Domain.Fqdn]; !ok {
			domainRoles[domainRole.Domain.Fqdn] = make([]string, 0)
		}

		domainRoles[domainRole.Domain.Fqdn] = append(domainRoles[domainRole.Domain.Fqdn], domainRole.DomainRole.Name)
	}

	return &management.User{
		ID:          user.ID,
		Username:    user.Username,
		CreatedAt:   user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   user.UpdatedAt.Format("2006-01-02 15:04:05"),
		UserRoles:   userRoles,
		DomainRoles: domainRoles,
	}
}

func (m *UserMapper) MapDatabaseToDtoList(users []*model.User) []*management.User {
	dtos := make([]*management.User, 0)

	for _, user := range users {
		dtos = append(dtos, m.MapDatabaseToDto(user))
	}

	return dtos
}

func (m *UserMapper) MapDatabaseToMinimalDto(user *model.User) *management.UserMinimal {
	return &management.UserMinimal{
		ID:       user.ID,
		Username: user.Username,
	}
}

func (m *UserMapper) MapDatabaseToMinimalDtoList(users []*model.User) []*management.UserMinimal {
	dtos := make([]*management.UserMinimal, 0)

	for _, user := range users {
		dtos = append(dtos, m.MapDatabaseToMinimalDto(user))
	}

	return dtos
}

func (m *UserMapper) MapUserCreatePayloadToDb(user *management.UserCreatePayload) *model.User {
	return &model.User{
		Username: user.Username,
		Password: user.Password,
	}
}
