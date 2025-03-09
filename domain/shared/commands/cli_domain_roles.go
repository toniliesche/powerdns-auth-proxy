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

func getCliDomainRoleCommand() *cli.Command {
	return &cli.Command{
		Name:  "domain-role",
		Usage: "manage domain roles",
		Subcommands: []*cli.Command{
			getCliDomainRoleListCommand(),
			getCliDomainRoleAddCommand(),
			getCliDomainRoleRemoveCommand(),
		},
	}
}

func getCliDomainRoleListCommand() *cli.Command {
	return &cli.Command{
		Name:      "list",
		Usage:     "list all domain roles of a user for a domain",
		Action:    app.RunCliDomainRoleList,
		ArgsUsage: " <username> <fqdn>",
	}
}

func getCliDomainRoleAddCommand() *cli.Command {
	return &cli.Command{
		Name:      "add",
		Usage:     "add a domain role to a user",
		Action:    app.RunCliDomainRoleAdd,
		ArgsUsage: " <username> <fqdn> <role>",
	}
}

func getCliDomainRoleRemoveCommand() *cli.Command {
	return &cli.Command{
		Name:      "remove",
		Usage:     "remove a domain role from a user",
		Action:    app.RunCliDomainRoleRemove,
		ArgsUsage: " <username> <fdqn> <role>",
	}
}
