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
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"powerdns-auth-proxy/domain/shared/auth/model"
	errors2 "powerdns-auth-proxy/domain/shared/http/errors"
	"powerdns-auth-proxy/domain/shared/interfaces"
	"time"
)

type JWTServiceMock struct {
	secretKey string
}

func (s *JWTServiceMock) Login(context *gin.Context) (*model.TokenResponse, errors2.HTTPError) {
	return s.createToken()
}

func (s *JWTServiceMock) Refresh(context *gin.Context) (*model.TokenResponse, errors2.HTTPError) {
	return s.createToken()
}

func (s *JWTServiceMock) createToken() (*model.TokenResponse, errors2.HTTPError) {
	now := time.Now()
	refreshTokenId := uuid.New()
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"aud": "audience",
			"exp": now.Add(time.Hour * 24 * 30).Unix(),
			"iss": "issuer",
			"jti": refreshTokenId.String(),
			"nbf": now.Add(time.Second * -30).Unix(),
			"sid": 1,
			"sub": "subject",
			"tti": "refresh",
		})

	refreshTokenString, err := refreshToken.SignedString([]byte(s.secretKey))
	if err != nil {
		return nil, errors2.NewInternalServerErrorError(fmt.Errorf("failed to sign refreshToken: %v", err))
	}

	accessTokenId := uuid.New()
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"aud": "audience",
			"exp": now.Add(time.Minute * 30).Unix(),
			"iid": refreshTokenId.String(),
			"iss": "issuer",
			"jti": accessTokenId.String(),
			"nbf": now.Add(time.Second * -30).Unix(),
			"sid": 1,
			"sub": "username",
			"tti": "access",
			"udr": nil,
			"usr": nil,
		})

	accessTokenString, err := accessToken.SignedString([]byte(s.secretKey))
	if err != nil {
		return nil, errors2.NewInternalServerErrorError(fmt.Errorf("failed to sign accessToken: %v", err))
	}

	return &model.TokenResponse{AccessToken: accessTokenString, RefreshToken: refreshTokenString, SequenceNumber: 1}, nil
}

func GetJWTServiceMock() interfaces.JWTServiceInterface {
	return &JWTServiceMock{secretKey: "0123456789abcdef0123456789abcdef"}
}
