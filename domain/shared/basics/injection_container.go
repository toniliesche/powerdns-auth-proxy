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

package basics

import (
	"gorm.io/gorm"
	"powerdns-auth-proxy/domain/admin/services/interfaces"
	"powerdns-auth-proxy/domain/shared/config"
	interfaces2 "powerdns-auth-proxy/domain/shared/database/repository/interfaces"
	sharedinterfaces "powerdns-auth-proxy/domain/shared/interfaces"
)

type InjectionContainer struct {
	Config             *config.Config
	DB                 *gorm.DB
	AdminAuthenticator sharedinterfaces.RequestAuthenticatorInterface
	ApiKeyRepository   interfaces2.ApiKeyRepositoryInterface
	ApiKeyService      sharedinterfaces.ApiKeyServiceInterface
	Authenticator      sharedinterfaces.RequestAuthenticatorInterface
	AuthService        sharedinterfaces.AuthenticationServiceInterface

	AdminDomainController          sharedinterfaces.ControllerInterface
	AdminIndexController           sharedinterfaces.ControllerInterface
	AdminJwtController             sharedinterfaces.ControllerInterface
	AdminSessionsController        sharedinterfaces.ControllerInterface
	AdminUserDomainRolesController sharedinterfaces.ControllerInterface
	AdminUserRolesController       sharedinterfaces.ControllerInterface
	AdminUserController            sharedinterfaces.ControllerInterface
	ApiController                  sharedinterfaces.ControllerInterface
	AutoprimariesController        sharedinterfaces.ControllerInterface
	BaseControllerAdminAPI         sharedinterfaces.BaseControllerInterface
	BaseControllerPowerDNS         sharedinterfaces.BaseControllerInterface
	CacheController                sharedinterfaces.ControllerInterface
	CryptokeysController           sharedinterfaces.ControllerInterface
	DomainRepository               interfaces2.DomainRepositoryInterface
	DomainRoleRepository           interfaces2.DomainRoleRepositoryInterface
	DomainService                  interfaces.DomainServiceInterface
	ForwardService                 sharedinterfaces.ForwardServiceInterface
	JWTService                     sharedinterfaces.JWTServiceInterface
	LoginService                   sharedinterfaces.LoginServiceInterface
	MetadataController             sharedinterfaces.ControllerInterface
	ResponseWriterAdminAPI         sharedinterfaces.ResponseWriterInterface
	ResponseWriterPowerDNS         sharedinterfaces.ResponseWriterInterface
	RoleRepository                 interfaces2.RoleRepositoryInterface
	RoleService                    interfaces.RoleServiceInterface
	SearchController               sharedinterfaces.ControllerInterface
	ServersController              sharedinterfaces.ControllerInterface
	SessionRepository              interfaces2.TokenSessionRepositoryInterface
	SessionService                 interfaces.SessionServiceInterface
	StatisticsController           sharedinterfaces.ControllerInterface
	TsigkeysController             sharedinterfaces.ControllerInterface
	UserDomainRoleRepository       interfaces2.UserDomainRoleRepositoryInterface
	UserDomainRoleService          interfaces.UserDomainRoleServiceInterface
	UserRepository                 interfaces2.UserRepositoryInterface
	UserRoleRepository             interfaces2.UserRoleRepositoryInterface
	UserRoleService                interfaces.UserRoleServiceInterface
	UserService                    interfaces.UserServiceInterface
	ZonesController                sharedinterfaces.ControllerInterface
}

func (c *InjectionContainer) GetPowerDnsControllers() []sharedinterfaces.ControllerInterface {
	return []sharedinterfaces.ControllerInterface{
		c.AutoprimariesController,
		c.CacheController,
		c.CryptokeysController,
		c.MetadataController,
		c.SearchController,
		c.ServersController,
		c.StatisticsController,
		c.TsigkeysController,
		c.ZonesController,
	}
}

func (c *InjectionContainer) GetSystemControllers() []sharedinterfaces.ControllerInterface {
	return []sharedinterfaces.ControllerInterface{
		c.AdminIndexController,
		c.AdminJwtController,
	}
}

func (c *InjectionContainer) GetAdminControllers() []sharedinterfaces.ControllerInterface {
	return []sharedinterfaces.ControllerInterface{
		c.AdminDomainController,
		c.AdminSessionsController,
		c.AdminUserDomainRolesController,
		c.AdminUserRolesController,
		c.AdminUserController,
	}
}
