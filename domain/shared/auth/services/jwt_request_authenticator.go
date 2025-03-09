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

package services

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"powerdns-auth-proxy/domain/shared/auth/model"
	"powerdns-auth-proxy/domain/shared/basics"
	"strings"
)

type JWTRequestAuthenticator struct {
	publicKey  *ecdsa.PublicKey
	audience   string
	userMapper *UserMapper
}

func (a *JWTRequestAuthenticator) Authenticate(context *gin.Context) (*model.User, error) {
	req := context.Request

	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		return nil, fmt.Errorf("missing authorization header")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, fmt.Errorf("invalid authorization header")
	}

	tokenString := parts[1]

	var accessToken *jwt.Token
	var err error
	if accessToken, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])

		}
		return a.publicKey, nil
	}); err != nil {
		return nil, err
	}

	if !accessToken.Valid {
		return nil, fmt.Errorf("access token is invalid")
	}

	claims := accessToken.Claims.(jwt.MapClaims)

	if claims["tti"].(string) != "access" {
		return nil, fmt.Errorf("invalid access token type")
	}

	if claims["aud"].(string) != a.audience {
		return nil, fmt.Errorf("invalid audience")
	}

	return a.userMapper.FromTokenClaims(claims)
}

func ProvideJWTRequestAuthenticator(container *basics.InjectionContainer) (*JWTRequestAuthenticator, error) {
	if container.Config == nil {
		return nil, basics.NewMissingDependencyError("could not provide jwt request authenticator: config could not be resolved")
	}

	if container.Config.JWT == nil {
		return nil, basics.NewMissingDependencyError("could not provide jwt request authenticator: jwt config could not be resolved")
	}

	jwtConfig := container.Config.JWT
	if err := jwtConfig.IsValid(); err != nil {
		return nil, err
	}

	parsedPublicKey, err := jwt.ParseECPublicKeyFromPEM([]byte(jwtConfig.PublicKey))
	if err != nil {
		return nil, err
	}

	return &JWTRequestAuthenticator{publicKey: parsedPublicKey, audience: jwtConfig.Audience, userMapper: &UserMapper{}}, nil
}
