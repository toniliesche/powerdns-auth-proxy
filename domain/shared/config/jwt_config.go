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

import "fmt"

type JWTConfig struct {
	Audience      string `yaml:"audience" json:"audience"`
	Issuer        string `yaml:"issuer" json:"issuer"`
	SecretKey     string `yaml:"secret_key" json:"secret_key"`
	SecretKeyPath string `yaml:"secret_key_path" json:"secret_key_path"`
	PublicKey     string `yaml:"public_key" json:"public_key"`
	PublicKeyPath string `yaml:"public_key_path" json:"public_key_path"`
}

func (c *JWTConfig) IsValid() error {
	if c.Audience == "" {
		return fmt.Errorf("JWT Config: audience is required")
	}

	if c.Issuer == "" {
		return fmt.Errorf("JWT Config: issuer is required")
	}

	if c.PublicKeyPath != "" && c.SecretKeyPath != "" {
		return nil
	}

	if c.PublicKey != "" && c.SecretKey != "" {
		return nil
	}

	return fmt.Errorf("JWT Config: public and secret key or their paths are required")
}

func (c *JWTConfig) Initialized() bool {
	return c.PublicKey != "" && c.SecretKey != ""
}
