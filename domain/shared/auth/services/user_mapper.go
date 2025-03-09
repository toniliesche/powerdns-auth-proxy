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
	"github.com/golang-jwt/jwt/v5"
	"powerdns-auth-proxy/domain/shared/auth/model"
	dbmodel "powerdns-auth-proxy/domain/shared/database/model"
)

type UserMapper struct {
}

func (m *UserMapper) FromDatabase(user *dbmodel.User) (*model.User, error) {
	roles := make([]string, 0, len(user.UserRoles))

	for _, role := range user.UserRoles {
		roles = append(roles, role.Role.Name)
	}

	domainRoles := make([]*model.UserDomainRole, 0, len(user.UserDomainRoles))

	for _, domainRole := range user.UserDomainRoles {
		domainRoles = append(domainRoles, &model.UserDomainRole{
			Domain: domainRole.Domain.FQDN,
			Role:   domainRole.DomainRole.Name,
		})
	}

	return &model.User{
		Username:    user.Username,
		Roles:       roles,
		DomainRoles: domainRoles,
	}, nil
}

func (m *UserMapper) FromTokenClaims(claims jwt.MapClaims) (*model.User, error) {
	username := claims["sub"].(string)

	usrClaim := claims["usr"].([]interface{})
	userRoles := make([]string, 0, len(claims["usr"].([]interface{})))
	for _, roleName := range usrClaim {
		userRoles = append(userRoles, roleName.(string))
	}

	udrClaim := claims["udr"].(map[string]interface{})
	udr := make(map[string][]string)
	for domain, roles := range udrClaim {
		domainRoles := make([]string, 0, len(roles.([]interface{})))
		for _, roleName := range roles.([]interface{}) {
			domainRoles = append(domainRoles, roleName.(string))
		}
		udr[domain] = domainRoles
	}

	var userDomainRoles []*model.UserDomainRole
	for domain, roles := range udr {
		for _, role := range roles {
			userDomainRoles = append(userDomainRoles, &model.UserDomainRole{
				Domain: domain,
				Role:   role,
			})
		}
	}

	return &model.User{
		Username:    username,
		Roles:       userRoles,
		DomainRoles: userDomainRoles,
	}, nil
}
