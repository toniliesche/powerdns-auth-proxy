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

package http

import (
	"fmt"
	"net/http"
	"net/url"
	"powerdns-auth-proxy/domain/shared/basics"
)

type ForwardService struct {
	basePath string
	proto    string
	apiKey   string
	client   *http.Client
}

func (s *ForwardService) ForwardRequest(request *http.Request) (*http.Response, error) {
	request.Header.Del("Authorization")
	request.Header.Del("X-Api-Key")
	request.Header.Add("X-Api-Key", s.apiKey)

	request.Host = s.basePath

	newURL, err := url.Parse(request.URL.String())
	if err != nil {
		return nil, err
	}
	newURL.Host = s.basePath
	newURL.Scheme = s.proto

	request.URL = newURL
	request.RequestURI = ""

	return s.client.Do(request)
}

func ProvideForwardService(container *basics.InjectionContainer) (*ForwardService, error) {
	if container.Config == nil {
		return nil, basics.NewMissingDependencyError("could not provide forward service: config could not be resolved")
	}

	if container.Config.PowerDNS == nil {
		return nil, basics.NewMissingDependencyError("could not provide forward service: powerdns config could not be resolved")
	}

	powerDNSConfig := container.Config.PowerDNS

	var proto string
	if powerDNSConfig.SSL {
		proto = "https"
	} else {
		proto = "http"
	}

	basePath := powerDNSConfig.Host
	if (powerDNSConfig.SSL && powerDNSConfig.Port != 443) || (!powerDNSConfig.SSL && powerDNSConfig.Port != 80) {
		basePath = fmt.Sprintf("%s:%d", basePath, powerDNSConfig.Port)
	}

	return &ForwardService{
		basePath: basePath,
		proto:    proto,
		apiKey:   powerDNSConfig.ApiKey,
		client:   &http.Client{},
	}, nil
}
