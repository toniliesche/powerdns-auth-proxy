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

func getCliDomainCommand() *cli.Command {
	return &cli.Command{
		Name:  "domain",
		Usage: "manage domains of the api gateway",
		Subcommands: []*cli.Command{
			getCliDomainListCommand(),
			getCliDomainAddCommand(),
			getCliDomainDeleteCommand(),
		},
	}
}

func getCliDomainAddCommand() *cli.Command {
	return &cli.Command{
		Name:      "add",
		Usage:     "add a new domain",
		Action:    app.RunCliDomainAdd,
		ArgsUsage: " <fqdn>",
	}
}

func getCliDomainListCommand() *cli.Command {
	return &cli.Command{
		Name:   "list",
		Usage:  "list all domains",
		Action: app.RunCliDomainList,
	}
}

func getCliDomainDeleteCommand() *cli.Command {
	return &cli.Command{
		Name:      "delete",
		Usage:     "delete a domain",
		Action:    app.RunCliDomainDelete,
		ArgsUsage: " <fqdn>",
	}
}
