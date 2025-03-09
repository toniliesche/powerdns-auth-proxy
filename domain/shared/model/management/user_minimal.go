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

package management

import "powerdns-auth-proxy/domain/shared/database/model"

type UserMinimal struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	UpdatedAt string `json:"updated_at"`
}

func MapUserMinimalFromDb(user *model.User) *UserMinimal {
	return &UserMinimal{
		ID:        user.ID,
		Username:  user.Username,
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func MapUsersMinimalFromDbList(users []*model.User) []*UserMinimal {
	minimalUsers := make([]*UserMinimal, 0)
	for _, user := range users {
		minimalUsers = append(minimalUsers, MapUserMinimalFromDb(user))
	}
	return minimalUsers
}
