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

package config

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
	"os"
	"powerdns-auth-proxy/domain/shared/config/importer"
)

type DBImportConfig struct {
	Domains []*importer.Domain `yaml:"domains,omitempty"`
	Users   []*importer.User   `yaml:"shared,omitempty"`
}

func (c *DBImportConfig) Validate() error {
	var err error

	if c.Domains != nil {
		for number, domain := range c.Domains {
			if err = domain.Validate(); err != nil {
				return fmt.Errorf("invalid domain config [domain %d]: %w", number+1, err)
			}
		}
	}

	if c.Users != nil {
		for number, user := range c.Users {
			if err = user.Validate(); err != nil {
				return fmt.Errorf("invalid user config [user %d]: %w", number+1, err)
			}
		}
	}

	return nil
}

func ProvideImportConfig(context *cli.Context) (*DBImportConfig, error) {
	importFile := context.String("import-file")
	if importFile == "" {
		return nil, nil
	}

	var err error
	var configContent []byte
	if configContent, err = os.ReadFile(importFile); err != nil {
		if importFile == ImportConfig {
			return nil, nil
		}

		return nil, fmt.Errorf("could not read import file: %w", err)
	}

	var config *DBImportConfig
	if err = yaml.Unmarshal(configContent, &config); err != nil {
		return nil, fmt.Errorf("could not unmarshal import file: %w", err)
	}

	if err = config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid import config: %w", err)
	}

	return config, nil
}
