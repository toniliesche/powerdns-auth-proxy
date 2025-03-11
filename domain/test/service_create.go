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

package test

import (
	"fmt"
	"powerdns-auth-proxy/domain/shared/basics"
)

func ServiceCreateApiKey(container *basics.InjectionContainer, registry *Registry) error {
	apiKey, err := container.ApiKeyService.Create(UserUsername)
	if err != nil {
		return err
	}

	registry.Set("apiKey", apiKey.ApiKey)
	registry.Set("apiKeyIdentifier", apiKey.Identifier)
	registry.Set("apiKeyId", fmt.Sprintf("%d", apiKey.ID))

	return nil
}
