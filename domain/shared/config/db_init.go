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
)

type DBDomainRole struct {
	Name string `yaml:"name"`
}

func (r *DBDomainRole) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("name is required")

	}

	return nil
}

type DBRole struct {
	Name string `yaml:"name"`
}

func (r *DBRole) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("name is required")

	}

	return nil
}

type DBInitConfig struct {
	DomainRoles []*DBDomainRole `yaml:"domain_roles"`
	Roles       []*DBRole       `yaml:"roles"`
}

func (c *DBInitConfig) Validate() error {
	var err error
	if c.DomainRoles != nil {
		for number, domainRole := range c.DomainRoles {
			if err = domainRole.Validate(); err != nil {
				return fmt.Errorf("invalid domain role config [domain role %d]: %w", number+1, err)
			}
		}
	}

	if c.Roles != nil {
		for number, role := range c.Roles {
			if err = role.Validate(); err != nil {
				return fmt.Errorf("invalid role config [role %d]: %w", number+1, err)
			}
		}
	}

	return nil
}

func NewInitConfig(context *cli.Context) (*DBInitConfig, error) {
	initFile := context.String("init-file")
	if initFile == "" {
		return nil, nil
	}

	var err error
	var configContent []byte
	if configContent, err = os.ReadFile(initFile); err != nil {
		return nil, fmt.Errorf("could not read init file: %w", err)
	}

	var config *DBInitConfig
	if err = yaml.Unmarshal(configContent, &config); err != nil {
		return nil, fmt.Errorf("could not unmarshal init file: %w", err)
	}

	if err = config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid init config: %w", err)
	}

	return config, nil
}
