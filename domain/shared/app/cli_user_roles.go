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
	"powerdns-auth-proxy/domain/shared/setup"
)

func RunCliUserRoleList(context *cli.Context) error {
	if context.Args().Len() == 0 {
		return fmt.Errorf("missing argument: <username>")
	}

	roleService, err := getRoleService(context)
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Role"})

	roleList, err := roleService.ListRolesForUser(context.Args().Get(0))
	if err != nil {
		return fmt.Errorf("could not list roles for user: %w", err)
	}

	for _, role := range roleList {
		table.Append([]string{role.Role})
	}

	table.Render()

	return nil
}

func RunCliUserRoleAdd(context *cli.Context) error {
	if context.Args().Len() == 0 {
		return fmt.Errorf("missing arguments: <username> <role>")
	}

	if context.Args().Len() == 1 {
		return fmt.Errorf("missing argument: <role>")
	}

	roleService, err := getRoleService(context)
	if err != nil {
		return err
	}

	err = roleService.GrantRoleToUser(context.Args().Get(0), context.Args().Get(1))
	if err != nil {
		return fmt.Errorf("could not add role to user: %w", err)
	}

	fmt.Printf("role successfully added to user: %s [role: %s]\n", context.Args().Get(0), context.Args().Get(1))

	return nil
}

func RunCliUserRoleRemove(context *cli.Context) error {
	if context.Args().Len() == 0 {
		return fmt.Errorf("missing arguments: <username> <role>")
	}

	if context.Args().Len() == 1 {
		return fmt.Errorf("missing argument: <role>")
	}

	roleService, err := getRoleService(context)
	if err != nil {
		return err
	}

	err = roleService.RevokeRoleFromUser(context.Args().Get(0), context.Args().Get(1))
	if err != nil {
		return fmt.Errorf("could not remove role from user: %w", err)
	}

	fmt.Printf("role successfully removed from user: %s [role: %s]\n", context.Args().Get(0), context.Args().Get(1))

	return nil
}

func getRoleService(context *cli.Context) (interfaces.UserRoleServiceInterface, error) {
	container, err := setup.InitContainerCli(context)
	if err != nil {
		return nil, err
	}

	return container.UserRoleService, nil
}
