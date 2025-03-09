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

package user_roles

import (
	"github.com/gin-gonic/gin"
	interfaces2 "powerdns-auth-proxy/domain/admin/services/interfaces"
	"powerdns-auth-proxy/domain/shared/basics"
	"powerdns-auth-proxy/domain/shared/interfaces"
)

type UserRolesController struct {
	interfaces.BaseControllerInterface
	roleService interfaces2.UserRoleServiceInterface
	userService interfaces2.UserServiceInterface
}

func (c *UserRolesController) ConfigureGroupRoutes(router *gin.RouterGroup) {
	router.GET("/v1/users/:userId/roles", c.listUserRoles)
	router.POST("/v1/users/:userId/roles/grant", c.grantUserRoles)
	router.POST("/v1/users/:userId/roles/revoke", c.revokeUserRoles)
}

func ProvideUserRolesController(container *basics.InjectionContainer) (*UserRolesController, error) {
	if container.BaseController == nil {
		return nil, basics.NewMissingDependencyError("user roles controller could not be created: base controller could not be resolved")
	}

	if container.UserRoleService == nil {
		return nil, basics.NewMissingDependencyError("user roles controller could not be created: user role service could not be resolved")
	}

	return &UserRolesController{container.BaseController, container.UserRoleService, container.UserService}, nil
}
