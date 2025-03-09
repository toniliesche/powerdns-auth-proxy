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
	"github.com/urfave/cli/v2"
	"powerdns-auth-proxy/domain/shared/config"
	"powerdns-auth-proxy/domain/shared/database"
	"powerdns-auth-proxy/domain/shared/setup"
)

func RunDatabaseMigration(context *cli.Context) error {
	container, err := setup.InitContainerCli(context)
	if err != nil {
		return err
	}

	migrationService, err := database.ProvideMigrationService(container)
	if err != nil {
		return err
	}

	if err = migrationService.RunMigrations(); err != nil {
		return err
	}

	var initConfig *config.DBInitConfig
	if initConfig, err = config.ProvideInitConfig(context); err != nil {
		return err
	}

	if initConfig == nil {
		return nil
	}

	if err = migrationService.RunInit(initConfig); err != nil {
		return err
	}

	var importConfig *config.DBImportConfig
	if importConfig, err = config.ProvideImportConfig(context); err != nil {
		return err
	}

	if importConfig == nil {
		return nil
	}

	return migrationService.RunImport(importConfig)
}
