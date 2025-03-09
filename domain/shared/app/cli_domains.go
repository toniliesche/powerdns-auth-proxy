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
	"powerdns-auth-proxy/domain/admin/services/interfaces"
	"powerdns-auth-proxy/domain/shared/model/management"
	"powerdns-auth-proxy/domain/shared/setup"
	"strconv"
)

func RunCliDomainList(context *cli.Context) error {
	domainService, err := getDomainService(context)
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "FQDN", "Last update"})

	domainList, err := domainService.ListDomains()
	if err != nil {
		return fmt.Errorf("could not list domains: %w", err)
	}

	for _, domain := range domainList {
		table.Append([]string{strconv.Itoa(int(domain.ID)), domain.FQDN, domain.UpdatedAt})
	}

	table.Render()

	return nil
}

func RunCliDomainAdd(context *cli.Context) error {
	if context.Args().Len() == 0 {
		return fmt.Errorf("missing arguments: <fqdn>")
	}

	domainService, err := getDomainService(context)
	if err != nil {
		return err
	}

	payload := &management.DomainCreatePayload{FQDN: context.Args().Get(0)}
	if _, err = domainService.CreateDomain(payload); err != nil {
		return fmt.Errorf("could not create domain: %w", err)
	}

	fmt.Printf("domain successfully created: %s\n", context.Args().Get(0))

	return nil
}

func RunCliDomainDelete(context *cli.Context) error {
	if context.Args().Len() == 0 {
		return fmt.Errorf("missing arguments: <fqdn>")
	}

	domainService, err := getDomainService(context)
	if err != nil {
		return err
	}

	if err = domainService.DeleteDomain(context.Args().Get(0)); err != nil {
		return fmt.Errorf("could not delete domain: %w", err)
	}

	fmt.Printf("domain successfully deleted: %s\n", context.Args().Get(0))

	return nil
}

func getDomainService(context *cli.Context) (interfaces.DomainServiceInterface, error) {
	container, err := setup.InitContainerCli(context)
	if err != nil {
		return nil, err
	}

	return container.DomainService, nil
}
