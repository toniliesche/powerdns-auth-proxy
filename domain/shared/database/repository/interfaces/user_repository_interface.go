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

package interfaces

import "powerdns-auth-proxy/domain/shared/database/model"

type UserRepositoryInterface interface {
	SaveNewUser(user *model.User) error
	UpdateUser(user *model.User) error
	CheckExistenceByUsername(username string) bool
	FetchUser(username string) (*model.User, error)
	FetchUserByAPIKey(apiKey string) (*model.User, error)
	FetchUserByID(id uint) (*model.User, error)
	DeleteUser(username string) error
	DeleteUserByID(id uint) error
	FindAll(advanced bool) ([]*model.User, error)
}
