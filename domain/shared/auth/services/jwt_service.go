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
	"crypto"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"powerdns-auth-proxy/domain/shared/auth/model"
	"powerdns-auth-proxy/domain/shared/basics"
	dbmodel "powerdns-auth-proxy/domain/shared/database/model"
	databaseerrors "powerdns-auth-proxy/domain/shared/database/repository/errors"
	"powerdns-auth-proxy/domain/shared/database/repository/interfaces"
	httperrors "powerdns-auth-proxy/domain/shared/http/errors"
	sharedinterfaces "powerdns-auth-proxy/domain/shared/interfaces"
	"time"
)

type JWTService struct {
	loginService      sharedinterfaces.LoginServiceInterface
	secretKey         crypto.PrivateKey
	publicKey         *ecdsa.PublicKey
	sessionRepository interfaces.TokenSessionRepositoryInterface
	audience          string
	issuer            string
}

func (s *JWTService) createRefreshToken(user *dbmodel.User, session *dbmodel.TokenSession) (*jwt.Token, error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodES256,
		jwt.MapClaims{
			"aud": s.audience,
			"exp": session.ExpiresAt.Unix(),
			"iss": s.issuer,
			"jti": session.SessionID,
			"nbf": session.UpdatedAt.Add(time.Second * -30).Unix(),
			"sid": session.SequenceNumber,
			"sub": user.Username,
			"tti": "refresh",
		})

	return refreshToken, nil
}

func (s *JWTService) createAccessToken(user *dbmodel.User, session *dbmodel.TokenSession) (*jwt.Token, error) {
	roles := make([]string, 0, len(user.UserRoles))
	for _, dbRole := range user.UserRoles {
		if dbRole.Role.Name == "" {
			continue
		}

		roles = append(roles, dbRole.Role.Name)
	}

	domainRoles := make(map[string][]string)
	for _, dbDomainRole := range user.UserDomainRoles {
		if dbDomainRole.DomainRole.Name == "" {
			continue
		}

		if _, ok := domainRoles[dbDomainRole.Domain.Fqdn]; !ok {
			domainRoles[dbDomainRole.Domain.Fqdn] = make([]string, 0)
		}

		domainRoles[dbDomainRole.Domain.Fqdn] = append(domainRoles[dbDomainRole.Domain.Fqdn], dbDomainRole.DomainRole.Name)
	}

	accessTokenId := uuid.New()
	accessToken := jwt.NewWithClaims(jwt.SigningMethodES256,
		jwt.MapClaims{
			"aud": s.audience,
			"exp": session.UpdatedAt.Add(time.Minute * 30).Unix(),
			"iid": session.SessionID,
			"iss": s.issuer,
			"jti": accessTokenId.String(),
			"nbf": session.UpdatedAt.Add(time.Second * -30).Unix(),
			"sid": session.SequenceNumber,
			"sub": user.Username,
			"tti": "access",
			"udr": domainRoles,
			"usr": roles,
		})

	return accessToken, nil
}

func (s *JWTService) Login(context *gin.Context) (*model.TokenResponse, httperrors.HTTPError) {
	payload := model.LoginPayload{}
	err := json.NewDecoder(context.Request.Body).Decode(&payload)
	if err != nil {
		return nil, httperrors.NewBadRequestError(fmt.Errorf("failed json decoding: %v", err))
	}

	if err := payload.Verify(); err != nil {
		return nil, err
	}

	var dbUser *dbmodel.User
	if dbUser, err = s.loginService.LoginByCredentials(payload.Username, payload.Password); err != nil {
		return nil, httperrors.NewUnauthorizedError(err)
	}

	refreshTokenId := uuid.New()

	now := time.Now()
	dbTokenSession := &dbmodel.TokenSession{
		ExpiresAt:      now.Add(time.Hour * 24 * 30),
		SessionID:      refreshTokenId.String(),
		SequenceNumber: 1,
		UserID:         dbUser.ID,
	}

	if err = s.sessionRepository.SaveNewTokenSession(dbTokenSession); err != nil {
		return nil, httperrors.NewInternalServerErrorError(fmt.Errorf("failed to save token session: %v", err))
	}

	refreshToken, err := s.createRefreshToken(dbUser, dbTokenSession)
	accessToken, err := s.createAccessToken(dbUser, dbTokenSession)

	return s.createResponse(accessToken, refreshToken, dbTokenSession)
}

func (s *JWTService) Refresh(context *gin.Context) (*model.TokenResponse, httperrors.HTTPError) {
	payload := model.RefreshPayload{}
	err := json.NewDecoder(context.Request.Body).Decode(&payload)
	if err != nil {
		return nil, httperrors.NewBadRequestError(fmt.Errorf("failed json decoding: %v", err))
	}

	if err := payload.Verify(); err != nil {
		return nil, err
	}

	var claims jwt.MapClaims
	if claims, err = s.parseToken(payload.RefreshToken); err != nil {
		return nil, httperrors.NewUnauthorizedError(err)
	}

	var dbTokenSession *dbmodel.TokenSession
	if dbTokenSession, err = s.sessionRepository.FetchTokenSession(claims["jti"].(string)); err != nil {
		if errors.As(err, &databaseerrors.ItemNotFoundError{}) {
			return nil, httperrors.NewUnauthorizedError(fmt.Errorf("failed to fetch token session: %v", err))
		}

		return nil, httperrors.NewInternalServerErrorError(fmt.Errorf("failed to fetch token session: %v", err))
	}

	if float64(dbTokenSession.SequenceNumber) > claims["sid"].(float64) {
		return nil, httperrors.NewUnauthorizedError(fmt.Errorf("refresh token is outdated"))
	}

	var dbUser *dbmodel.User
	if dbUser, err = s.loginService.LoginByUsername(claims["sub"].(string)); err != nil {
		return nil, httperrors.NewUnauthorizedError(err)
	}

	if err = s.sessionRepository.UpdateTokenSession(dbTokenSession); err != nil {
		return nil, httperrors.NewInternalServerErrorError(fmt.Errorf("failed to update token session: %v", err))
	}

	now := time.Now()
	dbTokenSession.SequenceNumber++
	dbTokenSession.ExpiresAt = now.Add(time.Hour * 24 * 30)

	refreshToken, err := s.createRefreshToken(dbUser, dbTokenSession)
	accessToken, err := s.createAccessToken(dbUser, dbTokenSession)

	err = s.sessionRepository.UpdateTokenSession(dbTokenSession)
	if err != nil {
		return nil, httperrors.NewInternalServerErrorError(fmt.Errorf("failed to update token session: %v", err))
	}

	return s.createResponse(accessToken, refreshToken, dbTokenSession)
}

func (s *JWTService) createResponse(accessToken *jwt.Token, refreshToken *jwt.Token, session *dbmodel.TokenSession) (*model.TokenResponse, httperrors.HTTPError) {
	accessTokenString, err := accessToken.SignedString(s.secretKey)
	if err != nil {
		return nil, httperrors.NewInternalServerErrorError(fmt.Errorf("failed to sign access token: %v", err))
	}

	refreshTokenString, err := refreshToken.SignedString(s.secretKey)
	if err != nil {
		return nil, httperrors.NewInternalServerErrorError(fmt.Errorf("failed to sign refresh token: %v", err))
	}

	return &model.TokenResponse{AccessToken: accessTokenString, RefreshToken: refreshTokenString, SessionID: session.SessionID, SequenceNumber: session.SequenceNumber}, nil
}

func (s *JWTService) parseToken(tokenString string) (jwt.MapClaims, error) {
	var refreshToken *jwt.Token
	var err error

	if refreshToken, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return s.publicKey, nil
	}); err != nil {
		return nil, httperrors.NewUnauthorizedError(fmt.Errorf("failed to parse refresh token: %v", err))
	}

	if !refreshToken.Valid {
		return nil, httperrors.NewUnauthorizedError(fmt.Errorf("refresh token is invalid"))
	}

	claims := refreshToken.Claims.(jwt.MapClaims)
	if claims["tti"].(string) != "refresh" {
		return nil, httperrors.NewUnauthorizedError(fmt.Errorf("invalid refresh token type"))
	}

	if claims["aud"].(string) != s.audience {
		return nil, httperrors.NewUnauthorizedError(fmt.Errorf("invalid audience"))
	}

	return claims, nil
}

func ProvideJWTService(container *basics.InjectionContainer) (*JWTService, error) {
	if container.Config == nil {
		return nil, basics.NewMissingDependencyError("could not provide jwt service: config could not be resolved")
	}

	if container.Config.JWT == nil {
		return nil, basics.NewMissingDependencyError("could not provide jwt service: jwt config could not be resolved")
	}

	jwtConfig := container.Config.JWT
	if err := jwtConfig.IsValid(); err != nil {
		return nil, basics.NewMissingDependencyError(fmt.Sprintf("could not provide jwt service: jwt config is invalid: %v", err))
	}

	parsedPublicKey, err := jwt.ParseECPublicKeyFromPEM([]byte(jwtConfig.PublicKey))
	if err != nil {
		return nil, basics.NewMissingDependencyError(fmt.Sprintf("could not provide jwt service: failed to parse public key: %v", err))
	}

	parsedSecretKey, err := jwt.ParseECPrivateKeyFromPEM([]byte(jwtConfig.SecretKey))
	if err != nil {
		return nil, basics.NewMissingDependencyError(fmt.Sprintf("could not provide jwt service: failed to parse secret key: %v", err))
	}

	return &JWTService{loginService: container.LoginService, sessionRepository: container.SessionRepository, publicKey: parsedPublicKey, secretKey: parsedSecretKey, audience: jwtConfig.Audience, issuer: jwtConfig.Issuer}, nil
}
