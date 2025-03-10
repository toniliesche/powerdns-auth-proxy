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
	"fmt"
	"gorm.io/gorm"
	"powerdns-auth-proxy/domain/shared/basics"
	"powerdns-auth-proxy/domain/shared/config"
	"powerdns-auth-proxy/domain/shared/config/importer"
	"powerdns-auth-proxy/domain/shared/database/model"
	interfaces2 "powerdns-auth-proxy/domain/shared/database/repository/interfaces"
	"powerdns-auth-proxy/domain/shared/security"
)

type MigrationService struct {
	database                 *gorm.DB
	domainRepository         interfaces2.DomainRepositoryInterface
	domainRoleRepository     interfaces2.DomainRoleRepositoryInterface
	roleRepository           interfaces2.RoleRepositoryInterface
	userDomainRoleRepository interfaces2.UserDomainRoleRepositoryInterface
	userRepository           interfaces2.UserRepositoryInterface
	userRoleRepository       interfaces2.UserRoleRepositoryInterface
}

func (s *MigrationService) RunMigrations() error {
	var err error

	if err = s.database.AutoMigrate(&model.Domain{}); err != nil {
		return err
	}

	if err = s.database.AutoMigrate(&model.DomainRole{}); err != nil {
		return err
	}

	if err = s.database.AutoMigrate(&model.Role{}); err != nil {
		return err
	}

	if err = s.database.AutoMigrate(&model.User{}); err != nil {
		return err
	}

	if err = s.database.AutoMigrate(&model.ApiKey{}); err != nil {
		return err
	}

	if err = s.database.AutoMigrate(&model.TokenSession{}); err != nil {
		return err
	}

	if err = s.database.AutoMigrate(&model.UserRole{}); err != nil {
		return err
	}

	return s.database.AutoMigrate(&model.UserDomainRole{})
}

func (s *MigrationService) RunInit(config *config.DBInitConfig) error {
	for _, role := range config.Roles {
		var err error
		if s.roleRepository.CheckExistenceByName(role.Name) {
			continue
		}

		newRole := &model.Role{
			Name: role.Name,
		}

		if err = s.roleRepository.SaveNewRole(newRole); err != nil {
			return fmt.Errorf("failed to save role %s: %w", role.Name, err)
		}
	}

	for _, domainRole := range config.DomainRoles {
		var err error
		if s.domainRoleRepository.CheckExistenceByName(domainRole.Name) {
			continue
		}

		newDomainRole := &model.DomainRole{
			Name: domainRole.Name,
		}

		if err = s.domainRoleRepository.SaveNewDomainRole(newDomainRole); err != nil {
			return fmt.Errorf("failed to save domain role %s: %w", domainRole.Name, err)
		}
	}

	return nil
}

func (s *MigrationService) RunImport(config *config.DBImportConfig) error {
	var err error
	if err = s.RunDomainImport(config); err != nil {
		return err
	}

	for _, user := range config.Users {
		if err = s.RunUserImport(user); err != nil {
			return err
		}

		var dbUser *model.User
		if dbUser, err = s.userRepository.FetchUser(user.Username); err != nil {
			return fmt.Errorf("failed to fetch user %s: %w", user.Username, err)
		}

		if err = s.RunUserRoleImport(dbUser, user); err != nil {
			return err
		}

		if err = s.RunUserDomainRoleImport(dbUser, user); err != nil {
			return err
		}
	}

	return nil
}

func (s *MigrationService) RunDomainImport(config *config.DBImportConfig) error {
	for _, domain := range config.Domains {
		var err error

		if domain.Deleted {
			if err = s.domainRepository.DeleteDomain(domain.Fqdn); err != nil {
				return fmt.Errorf("failed to delete domain %s: %w", domain.Fqdn, err)
			}

			continue
		}

		if s.domainRepository.CheckExistenceByFqdn(domain.Fqdn) {
			continue
		}

		newDomain := &model.Domain{
			Fqdn: domain.Fqdn,
		}

		if err = s.domainRepository.SaveNewDomain(newDomain); err != nil {
			return fmt.Errorf("failed to save domain %s: %w", domain.Fqdn, err)
		}
	}

	return nil
}

func (s *MigrationService) RunUserImport(user *importer.User) error {
	var err error
	if user.Deleted {
		if s.userRepository.CheckExistenceByUsername(user.Username) {
			if err = s.userRepository.DeleteUser(user.Username); err != nil {
				return fmt.Errorf("failed to delete user %s: %w", user.Username, err)
			}
		}

		return nil
	}

	if !s.userRepository.CheckExistenceByUsername(user.Username) {
		var hashedPassword string
		if hashedPassword, err = security.HashPassword(user.Password); err != nil {
			return fmt.Errorf("failed to hash password for user %s: %w", user.Username, err)
		}

		newUser := &model.User{
			Username: user.Username,
			Password: hashedPassword,
		}

		if err = s.userRepository.SaveNewUser(newUser); err != nil {
			return fmt.Errorf("failed to save user %s: %w", user.Username, err)
		}
	}

	return nil
}

func (s *MigrationService) RunUserRoleImport(dbUser *model.User, user *importer.User) error {
	var err error
	if user.HasGlobalRoles() {
		for _, role := range user.UserRoles {
			if role.Deleted {
				if dbUser.HasGlobalRole(role.Role) {
					var userRole *model.UserRole
					if userRole, err = dbUser.GetGlobalRole(role.Role); err != nil {
						return fmt.Errorf("failed to fetch user role %s: %w", role.Role, err)
					}

					if err = s.userRoleRepository.DeleteUserRole(userRole); err != nil {
						return fmt.Errorf("failed to delete user role %s: %w", role.Role, err)
					}
				}

				continue
			}

			if dbUser.HasGlobalRole(role.Role) {
				continue
			}

			var dbRole *model.Role
			if dbRole, err = s.roleRepository.FetchRoleByName(role.Role); err != nil {
				return fmt.Errorf("failed to fetch role %s: %w", role.Role, err)
			}

			userRole := &model.UserRole{
				UserID: dbUser.ID,
				User:   dbUser,
				RoleID: dbRole.ID,
				Role:   dbRole,
			}

			if err = s.userRoleRepository.SaveNewUserRole(userRole); err != nil {
				return fmt.Errorf("failed to save user role %s: %w", role.Role, err)
			}

		}
	}

	return nil
}

func (s *MigrationService) RunUserDomainRoleImport(dbUser *model.User, user *importer.User) error {
	var err error
	if user.HasDomainRoles() {
		for _, domainRole := range user.UserDomainRoles {
			if domainRole.Deleted {
				if dbUser.HasDomainRole(domainRole.Domain, domainRole.Role) {
					var userDomainRole *model.UserDomainRole
					if userDomainRole, err = dbUser.GetDomainRole(domainRole.Domain, domainRole.Role); err != nil {
						return fmt.Errorf("failed to fetch user domain role %s for %s: %w", domainRole.Role, domainRole.Domain, err)
					}

					if err = s.userDomainRoleRepository.DeleteUserDomainRole(userDomainRole); err != nil {
						return fmt.Errorf("failed to delete user domain role %s for %s: %w", domainRole.Role, domainRole.Domain, err)
					}
				}

				continue
			}

			if dbUser.HasDomainRole(domainRole.Domain, domainRole.Role) {
				continue
			}

			var dbDomain *model.Domain
			if dbDomain, err = s.domainRepository.FetchDomainByFqdn(domainRole.Domain); err != nil {
				return fmt.Errorf("failed to fetch domain %s: %w", domainRole.Domain, err)
			}

			var dbDomainRole *model.DomainRole
			if dbDomainRole, err = s.domainRoleRepository.FetchDomainRoleByName(domainRole.Role); err != nil {
				return fmt.Errorf("failed to fetch domain role %s: %w", domainRole.Role, err)
			}

			userDomainRole := &model.UserDomainRole{
				UserID:       dbUser.ID,
				User:         dbUser,
				DomainID:     dbDomain.ID,
				Domain:       dbDomain,
				DomainRoleID: dbDomainRole.ID,
				DomainRole:   dbDomainRole,
			}

			if err = s.userDomainRoleRepository.SaveNewUserDomainRole(userDomainRole); err != nil {
				return fmt.Errorf("failed to save user domain role %s: %w", domainRole.Role, err)
			}
		}
	}

	return nil
}

func ProvideMigrationService(container *basics.InjectionContainer) (*MigrationService, error) {
	if container.DB == nil {
		return nil, fmt.Errorf("database client could not be resolved")
	}

	if container.DomainRepository == nil {
		return nil, fmt.Errorf("domain repository could not be resolved")
	}

	if container.DomainRoleRepository == nil {
		return nil, fmt.Errorf("domain role repository could not be resolved")
	}

	if container.RoleRepository == nil {
		return nil, fmt.Errorf("role repository could not be resolved")
	}

	if container.UserDomainRoleRepository == nil {
		return nil, fmt.Errorf("user domain role repository could not be resolved")
	}

	if container.UserRepository == nil {
		return nil, fmt.Errorf("user repository could not be resolved")
	}

	if container.UserRoleRepository == nil {
		return nil, fmt.Errorf("user role repository could not be resolved")
	}

	return &MigrationService{
		database:                 container.DB,
		domainRepository:         container.DomainRepository,
		domainRoleRepository:     container.DomainRoleRepository,
		roleRepository:           container.RoleRepository,
		userDomainRoleRepository: container.UserDomainRoleRepository,
		userRepository:           container.UserRepository,
		userRoleRepository:       container.UserRoleRepository,
	}, nil
}
