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

package rdbms

import (
	"fmt"
	"gorm.io/gorm"
	"powerdns-auth-proxy/domain/shared/basics"
	"powerdns-auth-proxy/domain/shared/database/model"
	"powerdns-auth-proxy/domain/shared/database/repository/errors"
	"powerdns-auth-proxy/domain/shared/database/repository/interfaces"
)

type UserRepository struct {
	database         *gorm.DB
	apiKeyRepository interfaces.ApiKeyRepositoryInterface
}

func (r *UserRepository) SaveNewUser(user *model.User) error {
	return r.database.Create(user).Error
}

func (r *UserRepository) UpdateUser(user *model.User) error {
	return r.database.Save(user).Error
}

func (r *UserRepository) CheckExistenceByUsername(username string) bool {
	var count int64

	r.database.Model(&model.User{}).Where("username = ?", username).Count(&count)

	return count > 0
}

func (r *UserRepository) FetchUser(username string) (*model.User, error) {
	var user model.User
	result := r.database.
		Preload("UserRoles").
		Preload("UserRoles.Role").
		Preload("UserDomainRoles").
		Preload("UserDomainRoles.Domain").
		Preload("UserDomainRoles.DomainRole").
		Where("username = ?", username).
		Model(&user).
		Find(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.NewItemNotFoundError(fmt.Sprintf("user with username %s not found", username))
	}

	return &user, nil
}

func (r *UserRepository) FetchUserByApiKey(apiKey string) (*model.User, error) {
	var user model.User
	var dbApiKey *model.ApiKey
	var err error

	if dbApiKey, err = r.apiKeyRepository.FetchApiKey(apiKey); err != nil {
		return nil, err
	}

	result := r.database.
		Model(&user).
		Preload("UserRoles").
		Preload("UserRoles.Role").
		Preload("UserDomainRoles").
		Preload("UserDomainRoles.Domain").
		Preload("UserDomainRoles.DomainRole").
		Where("id = ?", dbApiKey.UserId).
		Find(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.NewItemNotFoundError(fmt.Sprintf("user with api key %s not found", apiKey))
	}

	return &user, nil
}

func (r *UserRepository) FetchUserById(id uint) (*model.User, error) {
	var user model.User
	result := r.database.
		Preload("UserRoles").
		Preload("UserRoles.Role").
		Preload("UserDomainRoles").
		Preload("UserDomainRoles.Domain").
		Preload("UserDomainRoles.DomainRole").
		Where("id = ?", id).
		Model(&user).
		Find(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.NewItemNotFoundError(fmt.Sprintf("user with id %d not found", id))
	}

	return &user, nil
}

func (r *UserRepository) FindAll(advanced bool) ([]*model.User, error) {
	if advanced {
		return r.findAllAdvanced()
	}

	return r.findAllBasic()
}

func (r *UserRepository) DeleteUser(username string) error {
	user, err := r.FetchUser(username)
	if err != nil {
		return err
	}

	newUsername := fmt.Sprintf("%s#deleted-%d", user.Username, user.ID)
	err = r.database.Model(&model.User{}).Where("id = ?", user.ID).Update("username", newUsername).Error
	if err != nil {
		return err
	}

	return r.database.Where("id = ?", user.ID).Delete(&model.User{}).Error
}

func (r *UserRepository) DeleteUserById(id uint) error {
	user, err := r.FetchUserById(id)
	if err != nil {
		return err
	}

	newUsername := fmt.Sprintf("%s#deleted-%d", user.Username, user.ID)
	err = r.database.Model(&model.User{}).Where("id = ?", id).Update("username", newUsername).Error
	if err != nil {
		return err
	}

	return r.database.Where("id = ?", user.ID).Delete(&model.User{}).Error
}

func (r *UserRepository) findAllAdvanced() ([]*model.User, error) {
	var users []*model.User

	result := r.database.
		Model(&model.User{}).
		Preload("UserRoles").
		Preload("UserRoles.Role").
		Preload("UserDomainRoles").
		Preload("UserDomainRoles.Domain").
		Preload("UserDomainRoles.DomainRole").
		Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

func (r *UserRepository) findAllBasic() ([]*model.User, error) {
	var users []*model.User

	result := r.database.
		Model(&model.User{}).
		Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

func NewUserRepository(container *basics.InjectionContainer) (*UserRepository, error) {
	if container.DB == nil {
		return nil, basics.NewMissingDependencyError("could not provide user repository: database client could not be resolved")
	}

	if container.ApiKeyRepository == nil {
		return nil, basics.NewMissingDependencyError("could not provide user repository: api key repository could not be resolved")
	}

	return &UserRepository{database: container.DB, apiKeyRepository: container.ApiKeyRepository}, nil
}
