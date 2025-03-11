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

type Config struct {
	Debug      bool              `yaml:"debug" json:"debug"`
	AuthType   string            `yaml:"auth_type" json:"auth_type"`
	Database   string            `yaml:"database" json:"database"`
	LogPath    string            `yaml:"log_path" json:"log_path"`
	MySQL      *MySQLDBConfig    `yaml:"mysql,omitempty" json:"mysql,omitempty"`
	PostgreSQL *PostgreSQLConfig `yaml:"postgres,omitempty" json:"postgres,omitempty"`
	PowerDNS   *PowerDNSConfig   `yaml:"powerdns" json:"power_dns"`
	Sqlite     *SqliteConfig     `yaml:"sqlite,omitempty" json:"sqlite,omitempty"`
	JWT        *JWTConfig        `yaml:"jwt,omitempty" json:"jwt,omitempty"`
}

type PowerDNSConfig struct {
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
	SSL    bool   `yaml:"ssl"`
	ApiKey string `yaml:"api_key"`
}

func ProvideApplicationConfig(context *cli.Context) (*Config, error) {
	var err error
	var configContent []byte
	if configContent, err = os.ReadFile(context.String("config")); err != nil {
		return nil, fmt.Errorf("could not read config file: %w", err)
	}

	var config *Config
	if err := yaml.Unmarshal(configContent, &config); err != nil {
		return nil, fmt.Errorf("could not unmarshal config file: %w", err)
	}

	if config.JWT != nil {
		if !config.JWT.Initialized() {
			publicKeyBytes, err := os.ReadFile(config.JWT.PublicKeyPath)
			if err != nil {
				return nil, fmt.Errorf("could not read JWT public key: %w", err)
			}

			config.JWT.PublicKey = string(publicKeyBytes)

			secretKeyBytes, err := os.ReadFile(config.JWT.SecretKeyPath)
			if err != nil {
				return nil, fmt.Errorf("could not read JWT secret key: %w", err)
			}

			config.JWT.SecretKey = string(secretKeyBytes)
		}
	}

	return config, nil
}
