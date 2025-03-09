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

type RoleMapper struct{}

func (m *RoleMapper) MapDbToDto(role *model.Role) *management.Role {
	return &management.Role{
		ID:   role.ID,
		Name: role.Name,
	}
}

func (m *RoleMapper) MapDbToDtoList(roles []*model.Role) []*management.Role {
	roleList := make([]*management.Role, 0, len(roles))
	for _, role := range roles {
		roleList = append(roleList, m.MapDbToDto(role))
	}

	return roleList
}

func (m *RoleMapper) MapCreatePayloadToDb(payload *management.RoleCreatePayload) *model.Role {
	return &model.Role{
		Name: payload.Name,
	}
}
