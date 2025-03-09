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
	"github.com/urfave/cli/v2"
	"powerdns-auth-proxy/domain/admin/services/interfaces"
	"powerdns-auth-proxy/domain/shared/setup"
)

func RunCliDomainRoleList(context *cli.Context) error {
	if context.Args().Len() == 0 {
		return fmt.Errorf("missing arguments: <username> <fqdn>")
	}

	if context.Args().Len() == 1 {
		return fmt.Errorf("missing argument: <fqdn>")
	}

	roleService, err := getDomainRoleService(context)
	if err != nil {
		return err
	}

	_, err = roleService.ListDomainRolesForUser(context.Args().Get(1), context.Args().Get(0))

	return err
}

func RunCliDomainRoleAdd(context *cli.Context) error {
	if context.Args().Len() == 0 {
		return fmt.Errorf("missing arguments: <username> <fqdn> <role>")
	}

	if context.Args().Len() == 1 {
		return fmt.Errorf("missing arguments: <fqdn> <role>")
	}

	if context.Args().Len() == 2 {
		return fmt.Errorf("missing argument: <role>")
	}

	roleService, err := getDomainRoleService(context)
	if err != nil {
		return err
	}

	err = roleService.GrantDomainRoleToUser(context.Args().Get(1), context.Args().Get(0), context.Args().Get(2))
	if err != nil {
		return fmt.Errorf("could not add role to user: %w", err)
	}

	fmt.Printf("domain role successfully added to user: %s [domain: %s, role: %s]\n", context.Args().Get(0), context.Args().Get(1), context.Args().Get(2))

	return nil
}

func RunCliDomainRoleRemove(context *cli.Context) error {
	if context.Args().Len() == 0 {
		return fmt.Errorf("missing arguments: <username> <fqdn> <role>")
	}

	if context.Args().Len() == 1 {
		return fmt.Errorf("missing arguments: <fqdn> <role>")
	}

	if context.Args().Len() == 2 {
		return fmt.Errorf("missing argument: <role>")
	}

	roleService, err := getDomainRoleService(context)
	if err != nil {
		return err
	}

	err = roleService.RevokeDomainRoleFromUser(context.Args().Get(1), context.Args().Get(0), context.Args().Get(2))
	if err != nil {
		return fmt.Errorf("could not delete role from user: %w", err)
	}

	fmt.Printf("domain role successfully removed user: %s [domain: %s, role: %s]\n", context.Args().Get(0), context.Args().Get(1), context.Args().Get(2))

	return nil
}

func getDomainRoleService(context *cli.Context) (interfaces.UserDomainRoleServiceInterface, error) {
	container, err := setup.InitContainerCli(context)
	if err != nil {
		return nil, err
	}

	return container.UserDomainRoleService, nil
}
