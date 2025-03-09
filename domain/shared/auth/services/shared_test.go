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

package services_test

import (
	"powerdns-auth-proxy/domain/shared/basics"
	"powerdns-auth-proxy/domain/shared/database/model"
	"powerdns-auth-proxy/domain/shared/setup"
)

func getContainer(authenticationType string) *basics.InjectionContainer {
	container, _ := setup.InitContainerTest(&setup.TestConfig{AuthenticationType: authenticationType, EnableMockRepositories: true})

	container.UserRepository.SaveNewUser(&model.User{
		Username: "user",
		Password: "$2a$10$0yI64XAHi8q2SVZ.yuGpYeEu2Ufsyz0VgjvBRy5TknrQ9j4glrRQq",
		APIKeys: []*model.APIKey{
			{
				APIKey: "testapikey",
			},
		},
	})

	return container
}
