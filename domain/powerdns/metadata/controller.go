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

package metadata

import (
	"github.com/gin-gonic/gin"
	"powerdns-auth-proxy/domain/shared/basics"
	"powerdns-auth-proxy/domain/shared/interfaces"
)

type MetadataController struct {
	interfaces.BaseControllerInterface
}

func (c *MetadataController) ConfigureGroupRoutes(router *gin.RouterGroup) {
	router.GET("/v1/servers/:server/zones/:zone/metadata", c.listMetadata)
	router.POST("/v1/servers/:server/zones/:zone/metadata", c.createMetadata)
	router.GET("/v1/servers/:server/zones/:zone/metadata/:metadata", c.getMetadata)
	router.PUT("/v1/servers/:server/zones/:zone/metadata/:metadata", c.updateMetadata)
	router.DELETE("/v1/servers/:server/zones/:zone/metadata/:metadata", c.deleteMetadata)
}

func NewMetadataController(container *basics.InjectionContainer) (*MetadataController, error) {
	if container == nil {
		return nil, basics.NewMissingDependencyError("could not provide metadata controller: passed injection container is nil")
	}

	if container.BaseControllerPowerDNS == nil {
		return nil, basics.NewMissingDependencyError("could not provide metadata controller: base controller could not be resolved")
	}

	return &MetadataController{container.BaseControllerPowerDNS}, nil
}
