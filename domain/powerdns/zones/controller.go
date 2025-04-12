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

package zones

import (
	"github.com/gin-gonic/gin"
	"powerdns-auth-proxy/domain/shared/basics"
	"powerdns-auth-proxy/domain/shared/interfaces"
)

type ZonesController struct {
	interfaces.BaseControllerInterface
}

func (c *ZonesController) ConfigureGroupRoutes(router *gin.RouterGroup) {
	router.GET("/v1/servers/:server/zones", c.listZones)
	router.POST("/v1/servers/:server/zones", c.createZone)
	router.GET("/v1/servers/:server/zones/:zone", c.getZone)
	router.DELETE("/v1/servers/:server/zones/:zone", c.deleteZone)
	router.PUT("/v1/servers/:server/zones/:zone", c.updateZone)
	router.PUT("/v1/servers/:server/zones/:zone/axfr-retrieve", c.retrieveZoneFromMaster)
	router.PUT("/v1/servers/:server/zones/:zone/notify", c.sendNotifyToSlaves)
	router.GET("/v1/servers/:server/zones/:zone/export", c.exportZone)
	router.PATCH("/v1/servers/:server/zones/:zone", c.updateRRSet)
	router.PUT("/v1/servers/:server/zones/:zone/rectify", c.rectifyZone)
}

func ProvideZonesController(container *basics.InjectionContainer) (*ZonesController, error) {
	if container.BaseControllerPowerDNS == nil {
		return nil, basics.NewMissingDependencyError("zones controller could not be created: base controller could not be resolved")
	}

	return &ZonesController{container.BaseControllerPowerDNS}, nil
}
