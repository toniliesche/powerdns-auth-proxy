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

package user_domain_roles

import (
	"github.com/gin-gonic/gin"
	interfaces2 "powerdns-auth-proxy/domain/admin/services/interfaces"
	"powerdns-auth-proxy/domain/shared/basics"
	"powerdns-auth-proxy/domain/shared/interfaces"
)

type UserDomainRolesController struct {
	interfaces.BaseControllerInterface
	userDomainRoleService interfaces2.UserDomainRoleServiceInterface
	userService           interfaces2.UserServiceInterface
}

func (c *UserDomainRolesController) ConfigureGroupRoutes(router *gin.RouterGroup) {
	router.GET("/v1/users/:userId/domain-roles", c.listDomainRoles)
	router.POST("/v1/users/:userId/domain-roles/grant", c.grantDomainRoles)
	router.POST("/v1/users/:userId/domain-roles/revoke", c.revokeDomainRoles)
}

func NewUserDomainRolesController(container *basics.InjectionContainer) (*UserDomainRolesController, error) {
	if container == nil {
		return nil, basics.NewMissingDependencyError("could not provide user domain roles controller: passed injection container is nil")
	}

	if container.BaseControllerAdminAPI == nil {
		return nil, basics.NewMissingDependencyError("could not provide user domain roles controller: base controller could not be resolved")
	}

	if container.UserDomainRoleService == nil {
		return nil, basics.NewMissingDependencyError("could not provide user domain roles controller: domain role service could not be resolved")
	}

	return &UserDomainRolesController{container.BaseControllerAdminAPI, container.UserDomainRoleService, container.UserService}, nil
}
