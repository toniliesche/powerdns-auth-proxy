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

package database_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"powerdns-auth-proxy/domain/shared/database"
	"powerdns-auth-proxy/domain/shared/setup"
	"testing"
)

func TestRunMigrations(t *testing.T) {
	container, err := setup.InitContainerTest(&setup.TestConfig{})
	if !assert.NoError(t, err, "failed to init container") {
		return
	}

	if container.DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{}); !assert.NoError(t, err, "failed to open database") {
		return
	}

	migrationsService, err := database.ProvideMigrationService(container)
	if !assert.NoError(t, err, fmt.Sprintf("failed to provide migration service %s", err)) {
		return
	}

	if err = migrationsService.RunMigrations(); !assert.NoError(t, err, "failed to run migrations") {
		return
	}
}
