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

func getCliUserRoleCommand() *cli.Command {
	return &cli.Command{
		Name:  "user-role",
		Usage: "manage user roles",
		Subcommands: []*cli.Command{
			getCliUserRoleListCommand(),
			getCliUserRoleAddCommand(),
			getCliUserRoleRemoveCommand(),
		},
	}
}

func getCliUserRoleListCommand() *cli.Command {
	return &cli.Command{
		Name:      "list",
		Usage:     "list all roles of a user",
		Action:    app.RunCliUserRoleList,
		ArgsUsage: " <username>",
	}
}

func getCliUserRoleAddCommand() *cli.Command {
	return &cli.Command{
		Name:      "add",
		Usage:     "add a role to a user",
		Action:    app.RunCliUserRoleAdd,
		ArgsUsage: " <username> <role>",
	}
}

func getCliUserRoleRemoveCommand() *cli.Command {
	return &cli.Command{
		Name:      "remove",
		Usage:     "remove a role from a user",
		Action:    app.RunCliUserRoleRemove,
		ArgsUsage: " <username> <role>",
	}
}
