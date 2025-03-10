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
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v2"
	"os"
	"powerdns-auth-proxy/domain/shared/interfaces"
	"powerdns-auth-proxy/domain/shared/setup"
	"strconv"
)

func RunCliApiKeyAdd(context *cli.Context) error {
	if context.Args().Len() == 0 {
		return fmt.Errorf("missing arguments: <username>")
	}

	apiKeyService, err := getApiKeyService(context)
	if err != nil {
		return err
	}

	apiKey, err := apiKeyService.Create(context.Args().Get(0))
	if err != nil {
		return fmt.Errorf("could not generate api key: %w", err)
	}

	fmt.Printf("API key successfully generated for user: %s [user: %s, identifier: %s]\n", apiKey.ApiKey, context.Args().Get(0), apiKey.Identifier)

	return nil
}

func RunCliApiKeyList(context *cli.Context) error {
	if context.Args().Len() == 0 {
		return fmt.Errorf("missing arguments: <username>")
	}

	apiKeyService, err := getApiKeyService(context)
	if err != nil {
		return err
	}

	apiKeys, err := apiKeyService.List(context.Args().Get(0))
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Identifier", "Username", "Last used"})

	for _, apiKey := range apiKeys {
		table.Append([]string{
			strconv.Itoa(int(apiKey.ID)),
			apiKey.Identifier,
			apiKey.User.Username,
			apiKey.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	table.Render()

	return nil
}

func RunCliApiKeyDelete(context *cli.Context) error {
	if context.Args().Len() == 0 {
		return fmt.Errorf("missing arguments: <id>")
	}

	apiKeyService, err := getApiKeyService(context)
	if err != nil {
		return err
	}

	err = apiKeyService.Delete(context.Args().Get(0))
	if err != nil {
		return fmt.Errorf("could not delete api key: %w", err)
	}

	fmt.Printf("API key successfully deleted: %s\n", context.Args().Get(0))

	return nil
}

func getApiKeyService(context *cli.Context) (interfaces.ApiKeyServiceInterface, error) {
	container, err := setup.InitContainerCli(context)
	if err != nil {
		return nil, err
	}

	return container.ApiKeyService, nil
}
