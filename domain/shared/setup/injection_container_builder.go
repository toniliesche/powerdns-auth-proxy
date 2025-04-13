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

package setup

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"powerdns-auth-proxy/domain/admin/domains"
	"powerdns-auth-proxy/domain/admin/index"
	"powerdns-auth-proxy/domain/admin/jwt"
	database2 "powerdns-auth-proxy/domain/admin/services/database"
	"powerdns-auth-proxy/domain/admin/sessions"
	"powerdns-auth-proxy/domain/admin/user_domain_roles"
	"powerdns-auth-proxy/domain/admin/user_roles"
	"powerdns-auth-proxy/domain/admin/users"
	"powerdns-auth-proxy/domain/powerdns/api"
	"powerdns-auth-proxy/domain/powerdns/autoprimaries"
	"powerdns-auth-proxy/domain/powerdns/cache"
	"powerdns-auth-proxy/domain/powerdns/cryptokeys"
	"powerdns-auth-proxy/domain/powerdns/metadata"
	"powerdns-auth-proxy/domain/powerdns/search"
	"powerdns-auth-proxy/domain/powerdns/servers"
	"powerdns-auth-proxy/domain/powerdns/statistics"
	"powerdns-auth-proxy/domain/powerdns/tsigkeys"
	"powerdns-auth-proxy/domain/powerdns/zones"
	"powerdns-auth-proxy/domain/shared/auth"
	"powerdns-auth-proxy/domain/shared/auth/services"
	"powerdns-auth-proxy/domain/shared/basics"
	"powerdns-auth-proxy/domain/shared/config"
	"powerdns-auth-proxy/domain/shared/controller"
	"powerdns-auth-proxy/domain/shared/database"
	"powerdns-auth-proxy/domain/shared/database/repository/mock"
	"powerdns-auth-proxy/domain/shared/database/repository/rdbms"
	"powerdns-auth-proxy/domain/shared/http"
	"powerdns-auth-proxy/domain/shared/interfaces"
)

type InjectionContainerBuilder struct {
}

func (b *InjectionContainerBuilder) BuildTest(config *TestConfig) (*basics.InjectionContainer, error) {
	container := &basics.InjectionContainer{
		Config: b.provideMockConfig(config),
	}

	var err error
	if config.EnableMockRepositories {
		if err = b.initMockRepositories(container); err != nil {
			return nil, err
		}
	} else {
		if container.DB, err = database.ProvideTestDatabaseClient(); err != nil {
			return nil, fmt.Errorf("could not initialize database client: %w", err)
		}

		if err = b.initDbRepositories(container); err != nil {
			return nil, err
		}

		if config.RunDatabaseMigrations {
			migrationService, err := database.ProvideMigrationService(container)
			if err != nil {
				return nil, fmt.Errorf("could not initialize migration service: %w", err)
			}

			if err = migrationService.RunMigrations(); err != nil {
				return nil, fmt.Errorf("could not run migrations: %w", err)
			}
		}
	}

	if err = b.initDataServices(container); err != nil {
		return nil, err
	}

	if err = b.initHttpBasicsTest(container, config); err != nil {
		return nil, err
	}

	if err = b.initSystemControllers(container); err != nil {
		return nil, err
	}

	if err = b.initPowerdnsControllers(container); err != nil {
		return nil, err
	}

	if err = b.initAdminControllers(container); err != nil {
		return nil, err
	}

	return container, nil
}

func (b *InjectionContainerBuilder) Build(cfg *config.Config) (*basics.InjectionContainer, error) {
	container := &basics.InjectionContainer{
		Config: cfg,
	}

	var err error
	if container.DB, err = database.ProvideDatabaseClient(cfg); err != nil {
		return nil, err
	}

	if err = b.initDbRepositories(container); err != nil {
		return nil, err
	}

	if err = b.initDataServices(container); err != nil {
		return nil, err
	}

	if err = b.initHttpBasics(container); err != nil {
		return nil, err
	}

	if err = b.initSystemControllers(container); err != nil {
		return nil, err
	}

	if err = b.initPowerdnsControllers(container); err != nil {
		return nil, err
	}

	if err = b.initAdminControllers(container); err != nil {
		return nil, err
	}

	return container, nil
}

func (b *InjectionContainerBuilder) initHttpBasics(container *basics.InjectionContainer) error {
	var err error
	if container.ForwardService, err = http.ProvideForwardService(container); err != nil {
		return err
	}

	if container.ResponseWriterAdminAPI, err = http.ProvideResponseWriterAdminAPI(container); err != nil {
		return err
	}

	if container.ResponseWriterPowerDNS, err = http.ProvideResponseWriterPowerDNS(container); err != nil {
		return err
	}

	if container.LoginService, err = services.ProvideLoginService(container); err != nil {
		return err
	}

	if container.Authenticator, err = b.initAuthenticator(container); err != nil {
		return err
	}

	if container.AdminAuthenticator, err = services.ProvideJWTRequestAuthenticator(container); err != nil {
		return fmt.Errorf("could not initialize admin authenticator: %w", err)
	}

	if container.AuthService, err = auth.ProvideAuthenticationService(container); err != nil {
		return err
	}

	if container.JWTService, err = services.ProvideJWTService(container); err != nil {
		return err
	}

	if container.BaseControllerAdminAPI, err = controller.ProvideControllerAdminAPI(container); err != nil {
		return err
	}

	if container.BaseControllerPowerDNS, err = controller.ProvideControllerPowerDNS(container); err != nil {
		return err
	}

	if container.ApiController, err = api.ProvideApiController(container); err != nil {
		return err
	}

	return nil
}

func (b *InjectionContainerBuilder) initHttpBasicsTest(container *basics.InjectionContainer, config *TestConfig) error {
	var err error

	if config.EnableMockForwardService {
		if container.ForwardService, err = http.ProvideForwardServiceMock(container); err != nil {
			return err
		}
	} else {
		if container.ForwardService, err = http.ProvideForwardService(container); err != nil {
			return err
		}
	}

	if container.ResponseWriterAdminAPI, err = http.ProvideResponseWriterAdminAPI(container); err != nil {
		return err
	}

	if container.ResponseWriterPowerDNS, err = http.ProvideResponseWriterPowerDNS(container); err != nil {
		return err
	}

	if container.LoginService, err = services.ProvideLoginService(container); err != nil {
		return err
	}

	if container.Authenticator, err = b.initAuthenticator(container); err != nil {
		return err
	}

	if container.AdminAuthenticator, err = services.ProvideJWTRequestAuthenticator(container); err != nil {
		return fmt.Errorf("could not initialize admin authenticator: %w", err)
	}

	if config.EnableMockAuthentication {
		container.AuthService = &auth.AuthenticationServiceMock{}
	} else {
		if container.AuthService, err = auth.ProvideAuthenticationService(container); err != nil {
			return err
		}
	}

	if config.EnableMockJwtService {
		container.JWTService = &services.JWTServiceMock{}
	} else {
		if container.JWTService, err = services.ProvideJWTService(container); err != nil {
			return err
		}
	}

	if container.BaseControllerAdminAPI, err = controller.ProvideControllerAdminAPI(container); err != nil {
		return err
	}

	if container.BaseControllerPowerDNS, err = controller.ProvideControllerPowerDNS(container); err != nil {
		return err
	}

	if container.ApiController, err = api.ProvideApiController(container); err != nil {
		return err
	}

	return nil
}

func (b *InjectionContainerBuilder) initDbRepositories(container *basics.InjectionContainer) error {
	var err error
	if container.DomainRepository, err = rdbms.ProvideDomainRepository(container); err != nil {
		return err
	}

	if container.SessionRepository, err = rdbms.ProvideTokenSessionRepository(container); err != nil {
		return err
	}

	if container.ApiKeyRepository, err = rdbms.ProvideApiKeyRepository(container); err != nil {
		return err
	}

	if container.UserRepository, err = rdbms.ProvideUserRepository(container); err != nil {
		return err
	}

	if container.RoleRepository, err = rdbms.ProvideRoleRepository(container); err != nil {
		return err
	}

	if container.UserRoleRepository, err = rdbms.ProvideUserRoleRepository(container); err != nil {
		return err
	}

	if container.DomainRoleRepository, err = rdbms.ProvideDomainRoleRepository(container); err != nil {
		return err
	}

	if container.UserDomainRoleRepository, err = rdbms.ProvideUserDomainRoleRepository(container); err != nil {
		return err
	}

	return nil
}

func (b *InjectionContainerBuilder) initMockRepositories(container *basics.InjectionContainer) error {
	var err error
	if container.ApiKeyRepository, err = mock.ProvideApiKeyRepositoryMock(); err != nil {
		return err
	}

	if container.DomainRepository, err = mock.ProvideDomainRepositoryMock(); err != nil {
		return err
	}

	if container.DomainRoleRepository, err = mock.ProvideDomainRoleRepositoryMock(); err != nil {
		return err
	}

	if container.RoleRepository, err = mock.ProvideRoleRepositoryMock(); err != nil {
		return err
	}

	if container.SessionRepository, err = mock.ProvideTokenSessionRepositoryMock(); err != nil {
		return err
	}

	if container.UserDomainRoleRepository, err = mock.ProvideUserDomainRoleRepositoryMock(container); err != nil {
		return err
	}

	if container.UserRepository, err = mock.ProvideUserRepositoryMock(); err != nil {
		return err
	}

	if container.UserRoleRepository, err = mock.ProvideUserRoleRepositoryMock(); err != nil {
		return err
	}

	return nil
}

func (b *InjectionContainerBuilder) initDataServices(container *basics.InjectionContainer) error {
	var err error

	if container.ApiKeyService, err = services.ProvideApiKeyService(container); err != nil {
		return err
	}

	if container.DomainService, err = database2.ProvideDomainService(container); err != nil {
		return err
	}

	if container.RoleService, err = database2.ProvideRoleService(container); err != nil {
		return err
	}

	if container.SessionService, err = database2.ProvideSessionService(container); err != nil {
		return err
	}

	if container.UserService, err = database2.ProvideUserService(container); err != nil {
		return err
	}

	if container.UserRoleService, err = database2.ProvideUserRoleService(container); err != nil {
		return err
	}

	if container.UserDomainRoleService, err = database2.ProvideUserDomainRoleService(container); err != nil {
		return err
	}

	return nil
}

func (b *InjectionContainerBuilder) initAuthenticator(container *basics.InjectionContainer) (interfaces.RequestAuthenticatorInterface, error) {
	var service interfaces.RequestAuthenticatorInterface
	var err error

	switch container.Config.AuthType {
	case "api_key":
		service, err = services.ProvideApiKeyRequestAuthenticator(container)
	case "basic_auth":
		service, err = services.ProvideBasicAuthenticator(container)
	case "client_cert":
		service, err = services.ProvideClientCertRequestAuthenticator(container)
	case "jwt":
		service, err = services.ProvideJWTRequestAuthenticator(container)
	case "traefik_client_cert":
		service, err = services.ProvideTraefikClientCertRequestAuthenticator(container)
	default:
		service, err = nil, fmt.Errorf("unknown auth type: %s", container.Config.AuthType)
	}

	if err != nil {
		return nil, fmt.Errorf("could not initialize authenticator: %w", err)
	}

	return service, nil
}

func (b *InjectionContainerBuilder) initPowerdnsControllers(container *basics.InjectionContainer) error {
	var err error
	if container.AutoprimariesController, err = autoprimaries.ProvideAutoprimariesController(container); err != nil {
		return err
	}

	if container.CacheController, err = cache.ProvideCacheController(container); err != nil {
		return err
	}

	if container.CryptokeysController, err = cryptokeys.ProvideCryptokeysController(container); err != nil {
		return err
	}

	if container.MetadataController, err = metadata.ProvideMetadataController(container); err != nil {
		return err
	}

	if container.SearchController, err = search.ProvideSearchController(container); err != nil {
		return err
	}

	if container.ServersController, err = servers.ProvideServersController(container); err != nil {
		return err
	}

	if container.StatisticsController, err = statistics.ProvideStatisticsController(container); err != nil {
		return err
	}

	if container.TsigkeysController, err = tsigkeys.ProvideTsigkeysController(container); err != nil {
		return err
	}

	if container.ZonesController, err = zones.ProvideZonesController(container); err != nil {
		return err
	}

	return nil
}

func (b *InjectionContainerBuilder) initSystemControllers(container *basics.InjectionContainer) error {
	var err error
	if container.AdminIndexController, err = index.ProvideIndexController(container); err != nil {
		return err
	}

	if container.AdminJwtController, err = jwt.ProvideJwtController(container); err != nil {
		return err
	}

	return nil
}

func (b *InjectionContainerBuilder) initAdminControllers(container *basics.InjectionContainer) error {
	var err error
	if container.AdminDomainController, err = domains.ProvideDomainController(container); err != nil {
		return err
	}

	if container.AdminSessionsController, err = sessions.ProvideSessionsController(container); err != nil {
		return err
	}

	if container.AdminUserRolesController, err = user_roles.ProvideUserRolesController(container); err != nil {
		return err
	}

	if container.AdminUserDomainRolesController, err = user_domain_roles.ProvideUserDomainRolesController(container); err != nil {
		return err
	}

	if container.AdminUserController, err = users.ProvideUserController(container); err != nil {
		return err
	}

	return nil
}

func (b *InjectionContainerBuilder) provideMockConfig(testconfig *TestConfig) *config.Config {
	var authType string
	if testconfig.AuthenticationType != "" {
		authType = testconfig.AuthenticationType
	} else {
		authType = "basic_auth"
	}

	return &config.Config{
		AuthType: authType,
		JWT: &config.JWTConfig{
			Audience:  "test",
			Issuer:    "test",
			SecretKey: "-----BEGIN PRIVATE KEY-----\nMIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQg+eYLflI5egoj6Uru\n36AULMOEGXV+rGelSO9RyXKGQVKhRANCAAQJQ38vsy/x8WBTHOip9TDG07fHwsP8\naE1GPtIj5/2nY9WCkTajC/IDbaj/iC9hdQ1UXStT49AZfbf9IxXb6Fja\n-----END PRIVATE KEY-----",
			PublicKey: "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAECUN/L7Mv8fFgUxzoqfUwxtO3x8LD\n/GhNRj7SI+f9p2PVgpE2owvyA22o/4gvYXUNVF0rU+PQGX23/SMV2+hY2g==\n-----END PUBLIC KEY-----",
		},
		PowerDNS: &config.PowerDNSConfig{
			ApiKey: "test",
			Host:   "127.0.0.1",
			Port:   8081,
		},
	}
}

func InitContainerCli(context *cli.Context) (*basics.InjectionContainer, error) {
	var err error
	var cfg *config.Config

	if cfg, err = config.ProvideApplicationConfig(context); err != nil {
		return nil, err
	}

	containerBuilder := &InjectionContainerBuilder{}
	container, err := containerBuilder.Build(cfg)
	if err != nil {
		return nil, fmt.Errorf("could not initialize container: %w", err)
	}

	return container, nil
}

func InitContainerTest(config *TestConfig) (*basics.InjectionContainer, error) {
	containerBuilder := &InjectionContainerBuilder{}
	container, err := containerBuilder.BuildTest(config)
	if err != nil {
		return nil, fmt.Errorf("could not initialize container: %w", err)
	}

	return container, nil
}
