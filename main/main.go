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

package main

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"powerdns-auth-proxy/domain/shared/commands"
)

func main() {
	gateway := &cli.App{
		Name:     "PowerDNS API Gateway",
		HelpName: "gateway",
		Usage:    "A gateway for the PowerDNS API to provide more specific permissions",
		Commands: []*cli.Command{
			commands.NewRunCommand(),
			commands.NewDatabaseMigrateCommand(),
			commands.NewCliCommand(),
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "config",
				Usage: "path to the configuration file",
				Value: "/etc/powerdns-auth-proxy/config.yaml",
			},
		},
	}

	if err := gateway.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
