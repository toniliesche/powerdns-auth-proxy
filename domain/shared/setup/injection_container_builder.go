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
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
	"powerdns-auth-proxy/domain/admin/domains"
	"powerdns-auth-proxy/domain/admin/index"
	"powerdns-auth-proxy/domain/admin/jwt"
	dbservices "powerdns-auth-proxy/domain/admin/services/database"
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
	"powerdns-auth-proxy/domain/shared/log"
)

type InjectionContainerBuilder struct {
	logger *zerolog.Logger
}

func (b *InjectionContainerBuilder) BuildTest(config *TestConfig) (*basics.InjectionContainer, error) {
	container := &basics.InjectionContainer{
		Config: b.provideMockConfig(config),
	}

	logger := zerolog.Nop()
	container.Logger = &logger

	var err error
	if config.EnableMockRepositories {
		if err = b.initMockRepositories(container); err != nil {
			return nil, err
		}
	} else {
		if container.DB, err = database.NewTestDatabaseClient(container); err != nil {
			return nil, fmt.Errorf("could not initialize database client: %w", err)
		}

		if err = b.initDbRepositories(container); err != nil {
			return nil, err
		}

		if config.RunDatabaseMigrations {
			migrationService, err := database.NewMigrationService(container)
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

	if err = b.initHttpBasics(container, config); err != nil {
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

	b.logger.Debug().
		Msg("initializing logger")

	var err error
	if err = b.initLogger(container); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize logger")

		return nil, err
	}

	b.logger.Debug().
		Msg("initializing database client")

	if container.DB, err = database.NewDatabaseClient(container, cfg); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize database client")

		return nil, err
	}

	b.logger.Debug().
		Msg("initializing database repositories")

	if err = b.initDbRepositories(container); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize database repositories")

		return nil, err
	}

	b.logger.Debug().
		Msg("initializing data services")

	if err = b.initDataServices(container); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize data services")

		return nil, err
	}

	b.logger.Debug().
		Msg("initializing http basics")

	if err = b.initHttpBasics(container, nil); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize http basics")

		return nil, err
	}

	b.logger.Debug().
		Msg("initializing system controllers")

	if err = b.initSystemControllers(container); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize system controllers")

		return nil, err
	}

	b.logger.Debug().
		Msg("initializing powerdns controllers")

	if err = b.initPowerdnsControllers(container); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize powerdns controllers")

		return nil, err
	}

	b.logger.Debug().
		Msg("initializing admin controllers")

	if err = b.initAdminControllers(container); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize admin controllers")

		return nil, err
	}

	return container, nil
}

func (b *InjectionContainerBuilder) initLogger(container *basics.InjectionContainer) error {
	logger, err := log.NewLogger(container)

	if err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize logger")

		return fmt.Errorf("could not initialize logger: %w", err)
	}

	if logger == nil {
		b.logger.Error().
			Msg("logger is nil")

		return fmt.Errorf("logger is nil")
	}

	container.Logger = logger
	b.logger = logger

	return nil
}

func (b *InjectionContainerBuilder) initHttpBasics(container *basics.InjectionContainer, testCfg *TestConfig) error {
	var err error

	b.logger.Trace().
		Msg("initializing forward service")

	if testCfg != nil && testCfg.EnableMockForwardService {
		b.logger.Trace().
			Msg("initializing mock forward service")

		if container.ForwardService, err = http.NewForwardServiceMock(container); err != nil {
			b.logger.Error().
				Err(err).
				Msg("could not initialize mock forward service")

			return err
		}
	} else {
		b.logger.Trace().
			Msg("initializing regular forward service")

		if container.ForwardService, err = http.NewForwardService(container); err != nil {
			b.logger.Error().
				Err(err).
				Msg("could not initialize forward service")

			return err
		}
	}

	b.logger.Trace().
		Msg("initializing admin api response writer")

	if container.ResponseWriterAdminAPI, err = http.NewResponseWriterAdminAPI(container); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize admin api response writer")

		return err
	}

	b.logger.Trace().
		Msg("initializing powerdns response writer")

	if container.ResponseWriterPowerDNS, err = http.NewResponseWriterPowerDNS(container); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize powerdns response writer")

		return err
	}

	b.logger.Trace().
		Msg("initializing login service")

	if container.LoginService, err = services.NewLoginService(container); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize login service")

		return err
	}

	b.logger.Trace().
		Msg("initializing authenticator")

	if container.Authenticator, err = b.initAuthenticator(container); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize authenticator")

		return err
	}

	b.logger.Trace().
		Msg("initializing admin authenticator")

	if container.AdminAuthenticator, err = services.NewJWTRequestAuthenticator(container); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize admin authenticator")

		return fmt.Errorf("could not initialize admin authenticator: %w", err)
	}

	b.logger.Trace().
		Msg("initializing authentication service")

	if testCfg != nil && testCfg.EnableMockAuthentication {
		b.logger.Trace().
			Msg("initializing mock authentication service")

		container.AuthService = &auth.AuthenticationServiceMock{}
	} else {
		b.logger.Trace().
			Msg("initializing regular authentication service")

		if container.AuthService, err = auth.NewAuthenticationService(container); err != nil {
			b.logger.Error().
				Err(err).
				Msg("could not initialize authentication service")

			return err
		}
	}

	b.logger.Trace().
		Msg("initializing jwt service")

	if testCfg != nil && testCfg.EnableMockJwtService {
		b.logger.Trace().
			Msg("initializing mock jwt service")

		container.JWTService = &services.JWTServiceMock{}
	} else {
		b.logger.Trace().
			Msg("initializing regular jwt service")

		if container.JWTService, err = services.NewJWTService(container); err != nil {
			b.logger.Error().
				Err(err).
				Msg("could not initialize jwt service")

			return err
		}
	}

	b.logger.Trace().
		Msg("initializing admin api base controller")

	if container.BaseControllerAdminAPI, err = controller.NewBaseControllerAdminAPI(container); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize admin api base controller")

		return err
	}

	b.logger.Trace().
		Msg("initializing powerdns base controller")

	if container.BaseControllerPowerDNS, err = controller.NewBaseControllerPowerDNS(container); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize powerdns base controller")

		return err
	}

	b.logger.Trace().
		Msg("initializing api controller")

	if container.ApiController, err = api.NewApiController(container); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize api controller")

		return err
	}

	b.logger.Trace().
		Msg("initializing request id middleware")

	if container.RequestIDMiddleware, err = http.NewRequestIDMiddleware(container); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize request id middleware")

		return err
	}

	b.logger.Trace().
		Msg("initializing request log middleware")

	if container.RequestLogMiddleware, err = http.NewRequestLogMiddleware(container); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize request log middleware")

		return err
	}

	b.logger.Trace().
		Msg("initializing request content log middleware")

	if container.RequestContentLogMiddleware, err = http.NewRequestContentLogMiddleware(container); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize request content log middleware")

		return err
	}

	b.logger.Trace().
		Msg("initializing response content log middleware")

	if container.ResponseContentLogMiddleware, err = http.NewResponseLogMiddleware(container); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize response content log middleware")

		return err
	}

	return nil
}

func (b *InjectionContainerBuilder) initDbRepositories(container *basics.InjectionContainer) error {
	var err error

	b.logger.Trace().
		Msg("initializing domain database repository")

	if container.DomainRepository, err = rdbms.NewDomainRepository(container); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize domain database repository")

		return err
	}

	b.logger.Trace().
		Msg("initializing session database repository")

	if container.SessionRepository, err = rdbms.NewTokenSessionRepository(container); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize session database repository")

		return err
	}

	b.logger.Trace().
		Msg("initializing api key database repository")

	if container.ApiKeyRepository, err = rdbms.NewApiKeyRepository(container); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize api key database repository")

		return err
	}

	b.logger.Trace().
		Msg("initializing user database repository")

	if container.UserRepository, err = rdbms.NewUserRepository(container); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize user database repository")

		return err
	}

	b.logger.Trace().
		Msg("initializing role database repository")

	if container.RoleRepository, err = rdbms.NewRoleRepository(container); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize role database repository")

		return err
	}

	b.logger.Trace().
		Msg("initializing user role database repository")

	if container.UserRoleRepository, err = rdbms.NewUserRoleRepository(container); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize user role database repository")

		return err
	}

	b.logger.Trace().
		Msg("initializing domain role database repository")

	if container.DomainRoleRepository, err = rdbms.NewDomainRoleRepository(container); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize domain role database repository")

		return err
	}

	b.logger.Trace().
		Msg("initializing user domain role database repository")

	if container.UserDomainRoleRepository, err = rdbms.NewUserDomainRoleRepository(container); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize user domain role database repository")

		return err
	}

	return nil
}

func (b *InjectionContainerBuilder) initMockRepositories(container *basics.InjectionContainer) error {
	var err error

	b.logger.Trace().
		Msg("initializing mock api key repository")

	if container.ApiKeyRepository, err = mock.NewApiKeyRepositoryMock(container); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize mock api key repository")

		return err
	}

	b.logger.Trace().
		Msg("initializing mock domain repository")

	if container.DomainRepository, err = mock.NewDomainRepositoryMock(container); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize mock domain repository")

		return err
	}

	b.logger.Trace().
		Msg("initializing mock domain role repository")

	if container.DomainRoleRepository, err = mock.NewDomainRoleRepositoryMock(container); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize mock domain role repository")

		return err
	}

	b.logger.Trace().
		Msg("initializing mock role repository")

	if container.RoleRepository, err = mock.NewRoleRepositoryMock(container); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize mock role repository")

		return err
	}

	b.logger.Trace().
		Msg("initializing mock session repository")

	if container.SessionRepository, err = mock.NewTokenSessionRepositoryMock(container); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize mock session repository")

		return err
	}

	b.logger.Trace().
		Msg("initializing mock user domain role repository")

	if container.UserDomainRoleRepository, err = mock.NewUserDomainRoleRepositoryMock(container); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize mock user domain role repository")

		return err
	}

	b.logger.Trace().
		Msg("initializing mock user repository")

	if container.UserRepository, err = mock.NewUserRepositoryMock(container); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize mock user repository")

		return err
	}

	b.logger.Trace().
		Msg("initializing mock user role repository")

	if container.UserRoleRepository, err = mock.NewUserRoleRepositoryMock(container); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize mock user role repository")

		return err
	}

	return nil
}

func (b *InjectionContainerBuilder) initDataServices(container *basics.InjectionContainer) error {
	var err error

	b.logger.Trace().
		Msg("initializing api key service")

	if container.ApiKeyService, err = services.NewApiKeyService(container); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize api key service")

		return err
	}

	b.logger.Trace().
		Msg("initializing domain service")

	if container.DomainService, err = dbservices.NewDomainService(container); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize domain service")

		return err
	}

	b.logger.Trace().
		Msg("initializing role service")

	if container.RoleService, err = dbservices.NewRoleService(container); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize role service")

		return err
	}

	b.logger.Trace().
		Msg("initializing session service")

	if container.SessionService, err = dbservices.NewSessionService(container); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize session service")

		return err
	}

	b.logger.Trace().
		Msg("initializing user service")

	if container.UserService, err = dbservices.NewUserService(container); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize user service")

		return err
	}

	b.logger.Trace().
		Msg("initializing user role service")

	if container.UserRoleService, err = dbservices.NewUserRoleService(container); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize user role service")

		return err
	}

	b.logger.Trace().
		Msg("initializing user domain role service")

	if container.UserDomainRoleService, err = dbservices.NewUserDomainRoleService(container); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize user domain role service")

		return err
	}

	return nil
}

func (b *InjectionContainerBuilder) initAuthenticator(container *basics.InjectionContainer) (interfaces.RequestAuthenticatorInterface, error) {
	var service interfaces.RequestAuthenticatorInterface
	var err error

	b.logger.Trace().
		Msg("initializing authenticator")

	switch container.Config.AuthType {
	case "api_key":
		b.logger.Trace().
			Msg("initializing api key authenticator")

		service, err = services.NewApiKeyRequestAuthenticator(container)
	case "basic_auth":
		b.logger.Trace().
			Msg("initializing basic authenticator")

		service, err = services.NewBasicAuthenticator(container)
	case "client_cert":
		b.logger.Trace().
			Msg("initializing client cert authenticator")

		service, err = services.NewClientCertRequestAuthenticator(container)
	case "jwt":
		b.logger.Trace().
			Msg("initializing jwt authenticator")

		service, err = services.NewJWTRequestAuthenticator(container)
	case "traefik_client_cert":
		b.logger.Trace().
			Msg("initializing traefik client cert authenticator")

		service, err = services.NewTraefikClientCertRequestAuthenticator(container)
	default:
		b.logger.Error().
			Msgf("unknown auth type: %s", container.Config.AuthType)

		service, err = nil, fmt.Errorf("unknown auth type: %s", container.Config.AuthType)
	}

	if err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize authenticator")

		return nil, fmt.Errorf("could not initialize authenticator: %w", err)
	}

	return service, nil
}

func (b *InjectionContainerBuilder) initPowerdnsControllers(container *basics.InjectionContainer) error {
	var err error

	b.logger.Trace().
		Msg("initializing powerdns autoprimaries controller")

	if container.AutoprimariesController, err = autoprimaries.NewAutoprimariesController(container); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize autoprimaries controller")

		return err
	}

	b.logger.Trace().
		Msg("initializing powerdns cache controller")

	if container.CacheController, err = cache.NewCacheController(container); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize cache controller")

		return err
	}

	b.logger.Trace().
		Msg("initializing powerdns cryptokeys controller")

	if container.CryptokeysController, err = cryptokeys.NewCryptokeysController(container); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize cryptokeys controller")

		return err
	}

	b.logger.Trace().
		Msg("initializing powerdns metadata controller")

	if container.MetadataController, err = metadata.NewMetadataController(container); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize metadata controller")

		return err
	}

	b.logger.Trace().
		Msg("initializing powerdns search controller")

	if container.SearchController, err = search.NewSearchController(container); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize search controller")

		return err
	}

	b.logger.Trace().
		Msg("initializing powerdns servers controller")

	if container.ServersController, err = servers.NewServersController(container); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize servers controller")

		return err
	}

	b.logger.Trace().
		Msg("initializing powerdns statistics controller")

	if container.StatisticsController, err = statistics.NewStatisticsController(container); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize statistics controller")

		return err
	}

	b.logger.Trace().
		Msg("initializing powerdns tsigkeys controller")

	if container.TsigkeysController, err = tsigkeys.NewTsigkeysController(container); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize tsigkeys controller")

		return err
	}

	b.logger.Trace().
		Msg("initializing powerdns zones controller")

	if container.ZonesController, err = zones.NewZonesController(container); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize zones controller")

		return err
	}

	return nil
}

func (b *InjectionContainerBuilder) initSystemControllers(container *basics.InjectionContainer) error {
	var err error

	b.logger.Trace().
		Msg("initializing system admin index controller")

	if container.AdminIndexController, err = index.NewIndexController(container); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize admin index controller")

		return err
	}

	b.logger.Trace().
		Msg("initializing system admin jwt controller")

	if container.AdminJwtController, err = jwt.NewJwtController(container); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize admin jwt controller")

		return err
	}

	return nil
}

func (b *InjectionContainerBuilder) initAdminControllers(container *basics.InjectionContainer) error {
	var err error

	b.logger.Trace().
		Msg("initializing admin api domain controller")

	if container.AdminDomainController, err = domains.NewDomainController(container); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize admin api domain controller")

		return err
	}

	b.logger.Trace().
		Msg("initializing admin api sessions controller")

	if container.AdminSessionsController, err = sessions.NewSessionsController(container); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize admin api sessions controller")

		return err
	}

	b.logger.Trace().
		Msg("initializing admin api roles controller")

	if container.AdminUserRolesController, err = user_roles.NewUserRolesController(container); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize admin api roles controller")

		return err
	}

	b.logger.Trace().
		Msg("initializing admin api user domain roles controller")

	if container.AdminUserDomainRolesController, err = user_domain_roles.NewUserDomainRolesController(container); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize admin api user domain roles controller")

		return err
	}

	b.logger.Trace().
		Msg("initializing admin api users controller")

	if container.AdminUserController, err = users.NewUserController(container); err != nil {
		b.logger.Error().
			Err(err).
			Msg("could not initialize admin api users controller")

		return err
	}

	return nil
}

func (b *InjectionContainerBuilder) provideMockConfig(cfg *TestConfig) *config.Config {
	var authType string

	b.logger.Trace().
		Msg("creating mock application config")

	if cfg.AuthenticationType != "" {
		authType = cfg.AuthenticationType
	} else {
		authType = "basic_auth"
	}

	return &config.Config{
		AuthType: authType,
		LogLevel: "debug",
		LogPath:  "/dev/null",
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

func InitContainerCli(context *cli.Context, logger *zerolog.Logger) (*basics.InjectionContainer, error) {
	var err error
	var cfg *config.Config

	logger.Debug().
		Msg("creating application config")

	if cfg, err = config.NewApplicationConfig(context); err != nil {
		logger.Error().
			Err(err).
			Msg("could not create application config")

		return nil, err
	}

	containerBuilder := &InjectionContainerBuilder{
		logger: logger,
	}

	logger.Debug().
		Msg("creating injection container")

	container, err := containerBuilder.Build(cfg)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("could not create injection container")

		return nil, fmt.Errorf("could not initialize container: %w", err)
	}

	logger = container.Logger

	logger.Debug().
		Msg("injection container created")

	return container, nil
}

func InitContainerTest(config *TestConfig) (*basics.InjectionContainer, error) {
	logger := zerolog.Nop()
	containerBuilder := &InjectionContainerBuilder{
		logger: &logger,
	}

	container, err := containerBuilder.BuildTest(config)
	if err != nil {
		return nil, fmt.Errorf("could not initialize container: %w", err)
	}

	return container, nil
}
