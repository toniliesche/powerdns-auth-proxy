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

package cryptokeys

import (
	"github.com/gin-gonic/gin"
	"powerdns-auth-proxy/domain/shared/basics"
	"powerdns-auth-proxy/domain/shared/interfaces"
)

type CryptokeysController struct {
	interfaces.BaseControllerInterface
}

func (c *CryptokeysController) ConfigureGroupRoutes(router *gin.RouterGroup) {
	router.GET("/v1/servers/:server/zones/:zone/cryptokeys", c.listCryptoKeys)
	router.POST("/v1/servers/:server/zones/:zone/cryptokeys", c.createCryptoKey)
	router.GET("/v1/servers/:server/zones/:zone/cryptokeys/:cryptokey", c.getCryptoKey)
	router.PUT("/v1/servers/:server/zones/:zone/cryptokeys/:cryptokey", c.toggleCryptoKeyActiveState)
	router.DELETE("/v1/servers/:server/zones/:zone/cryptokeys/:cryptokey", c.deleteCryptoKey)
}

func ProvideCryptokeysController(container *basics.InjectionContainer) (*CryptokeysController, error) {
	if container.BaseControllerPowerDNS == nil {
		return nil, basics.NewMissingDependencyError("cryptokeys controller could not be created: base controller could not be resolved")
	}

	return &CryptokeysController{container.BaseControllerPowerDNS}, nil
}
