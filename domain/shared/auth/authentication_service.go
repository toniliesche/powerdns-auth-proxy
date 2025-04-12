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

package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	authmodel "powerdns-auth-proxy/domain/shared/auth/model"
	"powerdns-auth-proxy/domain/shared/basics"
	"powerdns-auth-proxy/domain/shared/http/errors"
	"powerdns-auth-proxy/domain/shared/interfaces"
	"strings"
)

type AuthenticationService struct {
	ruleSetProvider        *RuleSetProvider
	authenticator          interfaces.RequestAuthenticatorInterface
	adminAuthenticator     interfaces.RequestAuthenticatorInterface
	responseWriterAdminAPI interfaces.ResponseWriterInterface
	responseWriterPowerDNS interfaces.ResponseWriterInterface
}

func (s *AuthenticationService) Authentication(admin bool) gin.HandlerFunc {
	var authenticator interfaces.RequestAuthenticatorInterface
	if admin {
		authenticator = s.adminAuthenticator
	} else {
		authenticator = s.authenticator
	}

	return func(context *gin.Context) {
		user, err := authenticator.Authenticate(context)

		if err != nil {
			s.responseWriterPowerDNS.HandleError(context, errors.NewUnauthorizedError(err))
			context.Abort()
			return
		}

		context.Set("user", user)
		context.Next()
	}
}

func (s *AuthenticationService) AdminAuthentication() gin.HandlerFunc {
	return func(context *gin.Context) {
		if !s.CheckAccessOnResource(context, Admin, "") {
			s.responseWriterAdminAPI.HandleError(context, errors.NewForbiddenError(fmt.Errorf("you are not authorized to access this resource")))
			context.Abort()
			return
		}

		context.Next()
	}
}

func (s *AuthenticationService) CheckAccessOnResource(context *gin.Context, ruleSetName string, resource string) bool {
	userParam, found := context.Get("user")
	if !found {
		return false
	}

	user := userParam.(*authmodel.User)

	var ruleSet *RuleSet
	var err error

	if ruleSet, err = s.ruleSetProvider.GetRuleSet(ruleSetName); err != nil {
		return false
	}

	resource, _ = strings.CutSuffix(resource, ".")
	for _, rule := range ruleSet.rules {
		if rule.CheckAccessOnResource(user, resource) {
			return true
		}
	}

	return false
}

func ProvideAuthenticationService(container *basics.InjectionContainer) (*AuthenticationService, error) {
	if container.Authenticator == nil {
		return nil, fmt.Errorf("could not provide authentication service: authenticator could not be resolved")
	}

	if container.AdminAuthenticator == nil {
		return nil, fmt.Errorf("could not provide authentication service: admin authenticator could not be resolved")
	}

	if container.ResponseWriterAdminAPI == nil {
		return nil, fmt.Errorf("could not provide authentication service: response writer could not be resolved")
	}

	if container.ResponseWriterPowerDNS == nil {
		return nil, fmt.Errorf("could not provide authentication service: response writer could not be resolved")
	}

	return &AuthenticationService{
		authenticator:          container.Authenticator,
		adminAuthenticator:     container.AdminAuthenticator,
		responseWriterAdminAPI: container.ResponseWriterAdminAPI,
		responseWriterPowerDNS: container.ResponseWriterPowerDNS,
		ruleSetProvider:        GetRuleSetProvider(),
	}, nil
}
