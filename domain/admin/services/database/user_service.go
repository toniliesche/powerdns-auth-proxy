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

package database

import (
	"powerdns-auth-proxy/domain/shared/basics"
	"powerdns-auth-proxy/domain/shared/database/model"
	"powerdns-auth-proxy/domain/shared/database/repository/interfaces"
	"powerdns-auth-proxy/domain/shared/mappers"
	"powerdns-auth-proxy/domain/shared/model/management"
	"powerdns-auth-proxy/domain/shared/security"
)

type UserService struct {
	userRepository interfaces.UserRepositoryInterface
	mapper         *mappers.UserMapper
}

func (s *UserService) CreateUser(payload *management.UserCreatePayload) (uint, error) {
	hashedPassword, err := security.HashPassword(payload.Password)
	if err != nil {
		return 0, err
	}

	user := &model.User{
		Username: payload.Username,
		Password: hashedPassword,
	}

	err = s.userRepository.SaveNewUser(user)
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}

func (s *UserService) ListUsers() ([]*management.UserMinimal, error) {
	dbUsers, err := s.userRepository.FindAll(false)
	if err != nil {
		return nil, err
	}

	return s.mapper.MapDatabaseToMinimalDtoList(dbUsers), nil
}

func (s *UserService) ListUsersComplete() ([]*management.User, error) {
	dbUsers, err := s.userRepository.FindAll(true)
	if err != nil {
		return nil, err
	}

	return s.mapper.MapDatabaseToDtoList(dbUsers), nil
}

func (s *UserService) CheckIsAdminUser(userId uint) (bool, error) {
	user, err := s.userRepository.FetchUserById(userId)
	if err != nil {
		return false, err
	}

	return user.IsAdminUser(), nil
}

func (s *UserService) GetUser(username string) (*management.User, error) {
	dbUser, err := s.userRepository.FetchUser(username)
	if err != nil {
		return nil, err
	}

	return s.mapper.MapDatabaseToDto(dbUser), nil
}

func (s *UserService) GetUserById(userId uint) (*management.User, error) {
	dbUser, err := s.userRepository.FetchUserById(userId)
	if err != nil {
		return nil, err
	}

	return s.mapper.MapDatabaseToDto(dbUser), nil
}

func (s *UserService) DeleteUser(username string) error {
	return s.userRepository.DeleteUser(username)
}

func (s *UserService) DeleteUserById(userId uint) error {
	return s.userRepository.DeleteUserById(userId)
}

func (s *UserService) UpdateUserPasswordById(userId uint, password string) error {
	user, err := s.userRepository.FetchUserById(userId)
	if err != nil {
		return err
	}

	hashedPassword, err := security.HashPassword(password)
	if err != nil {
		return err
	}

	user.Password = hashedPassword

	return s.userRepository.UpdateUser(user)
}

func ProvideUserService(container *basics.InjectionContainer) (*UserService, error) {
	if container.UserRepository == nil {
		return nil, basics.NewMissingDependencyError("could not provide user service: user repository could not be resolved")
	}

	return &UserService{
		userRepository: container.UserRepository,
		mapper:         &mappers.UserMapper{},
	}, nil
}
