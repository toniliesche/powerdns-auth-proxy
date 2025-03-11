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
	"strings"
)

func RunCliUserList(context *cli.Context) error {
	userService, err := getUserService(context)
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Username", "Global Roles", "Domain Specific Roles", "Last Update"})

	userList, err := userService.ListUsersComplete()
	if err != nil {
		return fmt.Errorf("could not list users: %w", err)
	}

	for _, user := range userList {
		globalRoleList := make([]string, 0, len(user.UserRoles))
		for _, role := range user.UserRoles {
			if role == "" {
				continue
			}

			globalRoleList = append(globalRoleList, role)
		}

		domainRoleList := make([]string, 0, len(user.DomainRoles))
		for domain, domainRoles := range user.DomainRoles {
			if len(domainRoles) == 0 {
				continue
			}

			for _, domainRole := range domainRoles {
				if domainRole == "" {
					continue
				}

				domainRoleList = append(domainRoleList, fmt.Sprintf("%s[%s]", domain, domainRole))
			}
		}

		var globalRoles string
		if len(globalRoleList) == 0 {
			globalRoles = "none"
		} else {
			globalRoles = strings.Join(globalRoleList, ", ")
		}

		var domainRoles string
		if len(domainRoleList) == 0 {
			domainRoles = "none"
		} else {
			domainRoles = strings.Join(domainRoleList, ", ")
		}

		lastUpdate := user.UpdatedAt
		table.Append([]string{strconv.Itoa(int(user.ID)), user.Username, globalRoles, domainRoles, lastUpdate})
	}

	table.Render()

	return nil
}

func RunCliUserAdd(context *cli.Context) error {
	if context.Args().Len() == 0 {
		return fmt.Errorf("missing arguments: <username> <password>")
	}

	if context.Args().Len() == 1 {
		return fmt.Errorf("missing argument: <password>")
	}

	userService, err := getUserService(context)
	if err != nil {
		return err
	}

	payload := &management.UserCreatePayload{
		Username: context.Args().Get(0),
		Password: context.Args().Get(1),
	}

	if _, err = userService.CreateUser(payload); err != nil {
		return fmt.Errorf("could not create user: %w", err)
	}

	fmt.Printf("user successfully created: %s\n", context.Args().Get(0))

	return nil
}

func RunCliUserDelete(context *cli.Context) error {
	if context.Args().Len() == 0 {
		return fmt.Errorf("missing argument: <username>")
	}

	userService, err := getUserService(context)
	if err != nil {
		return err
	}

	if err = userService.DeleteUser(context.Args().Get(0)); err != nil {
		return fmt.Errorf("could not delete user: %w", err)
	}

	fmt.Printf("user successfully deleted: %s\n", context.Args().Get(0))

	return nil
}

func getUserService(context *cli.Context) (interfaces.UserServiceInterface, error) {
	container, err := setup.InitContainerCli(context)
	if err != nil {
		return nil, err
	}

	return container.UserService, nil
}
