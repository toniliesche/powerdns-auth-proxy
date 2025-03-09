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

func getCliUserCommand() *cli.Command {
	return &cli.Command{
		Name:  "user",
		Usage: "manage users of the api gateway",
		Subcommands: []*cli.Command{
			getCliUserListCommand(),
			getCliUserAddCommand(),
			getCliUserDeleteCommand(),
		},
	}
}

func getCliUserAddCommand() *cli.Command {
	return &cli.Command{
		Name:      "add",
		Usage:     "add a new user",
		Action:    app.RunCliUserAdd,
		ArgsUsage: " <username> <password>",
	}
}

func getCliUserListCommand() *cli.Command {
	return &cli.Command{
		Name:   "list",
		Usage:  "list all users",
		Action: app.RunCliUserList,
	}
}

func getCliUserDeleteCommand() *cli.Command {
	return &cli.Command{
		Name:      "delete",
		Usage:     "delete a user",
		Action:    app.RunCliUserDelete,
		ArgsUsage: " <username>",
	}
}
