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

package mock

import (
	"powerdns-auth-proxy/domain/shared/basics"
	"powerdns-auth-proxy/domain/shared/database/model"
	errors2 "powerdns-auth-proxy/domain/shared/database/repository/errors"
)

type UserRepositoryMock struct {
	users   map[string]*model.User
	counter uint
}

func (r *UserRepositoryMock) SaveNewUser(user *model.User) error {
	if r.CheckExistenceByUsername(user.Username) {
		return errors2.NewItemAlreadyExistsError("user already exists")
	}

	user.ID = r.counter
	r.users[user.Username] = user
	r.counter++

	return nil
}

func (r *UserRepositoryMock) UpdateUser(user *model.User) error {
	return nil
}

func (r *UserRepositoryMock) CheckExistenceByUsername(username string) bool {
	_, found := r.users[username]

	return found
}

func (r *UserRepositoryMock) FetchUser(username string) (*model.User, error) {
	user, found := r.users[username]
	if !found {
		return nil, errors2.NewItemNotFoundError("user not found")
	}

	return user, nil
}

func (r *UserRepositoryMock) FetchUserByApiKey(apiKey string) (*model.User, error) {
	for _, user := range r.users {
		for _, key := range user.ApiKeys {
			if key.ApiKey == apiKey {
				return user, nil
			}
		}
	}

	return nil, errors2.NewItemNotFoundError("user not found")
}

func (r *UserRepositoryMock) FetchUserById(id uint) (*model.User, error) {
	for _, user := range r.users {
		if user.ID == id {
			return user, nil
		}
	}

	return nil, errors2.NewItemNotFoundError("user not found")
}

func (r *UserRepositoryMock) FindAll(advanced bool) ([]*model.User, error) {
	users := make([]*model.User, 0, len(r.users))

	for _, user := range r.users {
		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepositoryMock) DeleteUser(username string) error {
	if !r.CheckExistenceByUsername(username) {
		return errors2.NewItemNotFoundError("user not found")
	}

	delete(r.users, username)

	return nil
}

func (r *UserRepositoryMock) DeleteUserById(id uint) error {
	for username, user := range r.users {
		if user.ID == id {
			delete(r.users, username)
			return nil
		}
	}

	return errors2.NewItemNotFoundError("user not found")
}

func NewUserRepositoryMock(container *basics.InjectionContainer) (*UserRepositoryMock, error) {
	if container == nil {
		return nil, basics.NewMissingDependencyError("could not provide user repository mock: passed injection container is nil")
	}

	return &UserRepositoryMock{
		users:   make(map[string]*model.User),
		counter: 1,
	}, nil
}
