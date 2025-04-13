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

package app

import (
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
	"powerdns-auth-proxy/domain/shared/basics"
	"powerdns-auth-proxy/domain/shared/log"
	"powerdns-auth-proxy/domain/shared/setup"
)

func RunApi(context *cli.Context) error {
	logger := log.NewTempLogger()
	logger.Info().
		Msg("Booting api gateway")

	container, err := setup.InitContainerCli(context, logger)
	if err != nil {
		return err
	}

	logger = container.Logger

	router := gin.Default()
	router.Use(container.RequestIDMiddleware.Middleware())

	if container.Config.Debug {
		router.Use(container.RequestLogMiddleware.Middleware())
		router.Use(container.ResponseLogMiddleware.Middleware())
	}

	logger.Info().
		Msg("Configuring api gateway routes")

	configureRoutes(router, container)

	logger.Info().
		Msg("Starting api gateway")

	if err = router.Run("0.0.0.0:8080"); err != nil {
		logger.Error().
			Err(err).
			Msg("An error occurred while running the api gateway")

		return err
	}

	logger.Info().
		Msg("Api gateway has been stopped")

	return nil
}

func configureRoutes(router *gin.Engine, container *basics.InjectionContainer) {
	for _, systemController := range container.GetSystemControllers() {
		systemController.ConfigureEngineRoutes(router)
	}

	baseRoute := router.Group("/api")
	container.ApiController.ConfigureGroupRoutes(baseRoute)

	apiRoutes := router.Group("/api")
	apiRoutes.Use(container.AuthService.Authentication(false))
	for _, apiController := range container.GetPowerDnsControllers() {
		apiController.ConfigureGroupRoutes(apiRoutes)
	}

	adminRoutes := router.Group("/admin")
	adminRoutes.Use(container.AuthService.Authentication(true))
	adminRoutes.Use(container.AuthService.AdminAuthentication())

	for _, adminController := range container.GetAdminControllers() {
		adminController.ConfigureGroupRoutes(adminRoutes)
	}
}
