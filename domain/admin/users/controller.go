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

package users

import (
	"github.com/gin-gonic/gin"
	admininterfaces "powerdns-auth-proxy/domain/admin/services/interfaces"
	"powerdns-auth-proxy/domain/shared/basics"
	"powerdns-auth-proxy/domain/shared/interfaces"
)

type UserController struct {
	interfaces.BaseControllerInterface
	userService admininterfaces.UserServiceInterface
}

func (c *UserController) ConfigureGroupRoutes(router *gin.RouterGroup) {
	router.GET("/v1/users", c.listUsers)
	router.POST("/v1/users", c.createUser)
	router.GET("/v1/users/details/:id", c.getUser)
	router.PATCH("/v1/users/:id", c.updateUser)
	router.DELETE("/v1/users/:id", c.deleteUser)

}

func NewUserController(container *basics.InjectionContainer) (*UserController, error) {
	if container == nil {
		return nil, basics.NewMissingDependencyError("could not provide user controller: passed injection container is nil")
	}

	if container.BaseControllerAdminAPI == nil {
		return nil, basics.NewMissingDependencyError("could not provide user controller: base controller could not be resolved")
	}

	if container.UserService == nil {
		return nil, basics.NewMissingDependencyError("could not provide user controller: user service could not be resolved")
	}

	return &UserController{container.BaseControllerAdminAPI, container.UserService}, nil
}
