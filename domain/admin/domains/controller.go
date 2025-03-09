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

package domains

import (
	"github.com/gin-gonic/gin"
	admininterfaces "powerdns-auth-proxy/domain/admin/services/interfaces"
	"powerdns-auth-proxy/domain/shared/basics"
	"powerdns-auth-proxy/domain/shared/interfaces"
)

type DomainController struct {
	interfaces.BaseControllerInterface
	domainService admininterfaces.DomainServiceInterface
}

func (c *DomainController) ConfigureGroupRoutes(router *gin.RouterGroup) {
	router.GET("/v1/domains", c.listDomains)
	router.POST("/v1/domains", c.createDomain)
	router.GET("/v1/domains/details/:id", c.getDomain)
	router.DELETE("/v1/domains/:id", c.deleteDomain)
}

func ProvideDomainController(container *basics.InjectionContainer) (*DomainController, error) {
	if container.BaseController == nil {
		return nil, basics.NewMissingDependencyError("domains controller could not be created: base controller could not be resolved")
	}

	if container.DomainService == nil {
		return nil, basics.NewMissingDependencyError("domains controller could not be created: domain service could not be resolved")
	}

	return &DomainController{container.BaseController, container.DomainService}, nil
}
