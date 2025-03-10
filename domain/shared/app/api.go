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
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
	"powerdns-auth-proxy/domain/shared/basics"
	"powerdns-auth-proxy/domain/shared/http"
	"powerdns-auth-proxy/domain/shared/setup"
)

func RunApi(context *cli.Context) error {
	fmt.Println("Starting the api gateway")

	container, err := setup.InitContainerCli(context)
	if err != nil {
		return err
	}

	fmt.Println("Container initialized")

	router := gin.Default()
	router.Use(http.ProvideRequestIDMiddleware("powerdns-auth-proxy").Middleware())

	fmt.Println("Configuring routes")

	configureRoutes(router, container)

	fmt.Println("Routes configured")

	fmt.Println("Running the api gateway")

	if err = router.Run("0.0.0.0:8080"); err != nil {
		return fmt.Errorf("an error occured while running the api gateway: %w", err)
	}

	fmt.Println("API gateway stopped")

	return nil
}

func configureRoutes(router *gin.Engine, container *basics.InjectionContainer) {
	for _, systemController := range container.GetSystemControllers() {
		systemController.ConfigureEngineRoutes(router)
	}

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
