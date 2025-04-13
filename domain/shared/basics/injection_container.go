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
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	serviceinterfaces "powerdns-auth-proxy/domain/admin/services/interfaces"
	"powerdns-auth-proxy/domain/shared/config"
	dbrepointerfaces "powerdns-auth-proxy/domain/shared/database/repository/interfaces"
	"powerdns-auth-proxy/domain/shared/http/interfaces"
	sharedinterfaces "powerdns-auth-proxy/domain/shared/interfaces"
)

type InjectionContainer struct {
	Config             *config.Config
	DB                 *gorm.DB
	AdminAuthenticator sharedinterfaces.RequestAuthenticatorInterface
	ApiKeyRepository   dbrepointerfaces.ApiKeyRepositoryInterface
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
	DomainRepository               dbrepointerfaces.DomainRepositoryInterface
	DomainRoleRepository           dbrepointerfaces.DomainRoleRepositoryInterface
	DomainService                  serviceinterfaces.DomainServiceInterface
	ForwardService                 sharedinterfaces.ForwardServiceInterface
	JWTService                     sharedinterfaces.JWTServiceInterface
	Logger                         *zerolog.Logger
	LoginService                   sharedinterfaces.LoginServiceInterface
	MetadataController             sharedinterfaces.ControllerInterface
	RequestContentLogMiddleware    interfaces.GinMiddleware
	RequestIDMiddleware            interfaces.GinMiddleware
	RequestLogMiddleware           interfaces.GinMiddleware
	ResponseContentLogMiddleware   interfaces.GinMiddleware
	ResponseWriterAdminAPI         sharedinterfaces.ResponseWriterInterface
	ResponseWriterPowerDNS         sharedinterfaces.ResponseWriterInterface
	RoleRepository                 dbrepointerfaces.RoleRepositoryInterface
	RoleService                    serviceinterfaces.RoleServiceInterface
	SearchController               sharedinterfaces.ControllerInterface
	ServersController              sharedinterfaces.ControllerInterface
	SessionRepository              dbrepointerfaces.TokenSessionRepositoryInterface
	SessionService                 serviceinterfaces.SessionServiceInterface
	StatisticsController           sharedinterfaces.ControllerInterface
	TsigkeysController             sharedinterfaces.ControllerInterface
	UserDomainRoleRepository       dbrepointerfaces.UserDomainRoleRepositoryInterface
	UserDomainRoleService          serviceinterfaces.UserDomainRoleServiceInterface
	UserRepository                 dbrepointerfaces.UserRepositoryInterface
	UserRoleRepository             dbrepointerfaces.UserRoleRepositoryInterface
	UserRoleService                serviceinterfaces.UserRoleServiceInterface
	UserService                    serviceinterfaces.UserServiceInterface
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
