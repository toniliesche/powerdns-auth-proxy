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

package sessions

import (
	"github.com/gin-gonic/gin"
	interfaces2 "powerdns-auth-proxy/domain/admin/services/interfaces"
	"powerdns-auth-proxy/domain/shared/basics"
	"powerdns-auth-proxy/domain/shared/interfaces"
	"powerdns-auth-proxy/domain/shared/mappers"
)

type SessionsController struct {
	interfaces.BaseControllerInterface
	sessionService interfaces2.SessionServiceInterface
	userService    interfaces2.UserServiceInterface
	mapper         *mappers.SessionMapper
}

func (c *SessionsController) ConfigureGroupRoutes(router *gin.RouterGroup) {
	router.GET("/v1/users/:userId/sessions", c.listUserSession)
	router.POST("/v1/users/:userId/sessions/logout", c.logoutUserSession)
}

func NewSessionsController(container *basics.InjectionContainer) (*SessionsController, error) {
	if container == nil {
		return nil, basics.NewMissingDependencyError("could not provide sessions controller: passed injection container is nil")
	}

	if container.BaseControllerAdminAPI == nil {
		return nil, basics.NewMissingDependencyError("could not provide sessions controller: base controller could not be resolved")
	}

	if container.SessionService == nil {
		return nil, basics.NewMissingDependencyError("could not provide sessions controller: session service could not be resolved")
	}

	if container.UserService == nil {
		return nil, basics.NewMissingDependencyError("could not provide sessions controller: user service could not be resolved")
	}

	return &SessionsController{
		BaseControllerInterface: container.BaseControllerAdminAPI,
		sessionService:          container.SessionService,
		userService:             container.UserService,
		mapper:                  &mappers.SessionMapper{},
	}, nil
}
