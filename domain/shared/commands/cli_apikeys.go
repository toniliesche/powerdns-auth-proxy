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

package commands

import (
	"github.com/urfave/cli/v2"
	"powerdns-auth-proxy/domain/shared/app"
)

func getCliApiKeyCommand() *cli.Command {
	return &cli.Command{
		Name:  "api-key",
		Usage: "manage api keys of the api gateway",
		Subcommands: []*cli.Command{
			getCliApiKeyListCommand(),
			getCliApiKeyAddCommand(),
			getCliApiKeyDeleteCommand(),
		},
	}
}

func getCliApiKeyAddCommand() *cli.Command {
	return &cli.Command{
		Name:      "add",
		Usage:     "generate a new api key for user",
		Action:    app.RunCliApiKeyAdd,
		ArgsUsage: " <username>",
	}
}

func getCliApiKeyListCommand() *cli.Command {
	return &cli.Command{
		Name:   "list",
		Usage:  "list all api keys",
		Action: app.RunCliApiKeyList,
	}
}

func getCliApiKeyDeleteCommand() *cli.Command {
	return &cli.Command{
		Name:      "delete",
		Usage:     "delete a api key",
		Action:    app.RunCliApiKeyDelete,
		ArgsUsage: " <id>",
	}
}
