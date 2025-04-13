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

package zones_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	nethttp "net/http"
	"powerdns-auth-proxy/domain/powerdns/zones"
	"powerdns-auth-proxy/domain/shared/http"
	"powerdns-auth-proxy/domain/shared/setup"
	"powerdns-auth-proxy/domain/test"
	"testing"
)

func TestConfigureRoutes(t *testing.T) {
	router := gin.New()
	controller, err := getController()
	if !assert.NoError(t, err, "could not initialize TestConfigureRoutes") {
		return
	}
	controller.ConfigureGroupRoutes(router.Group("/api"))

	if !assert.Equal(t, 10, len(router.Routes()), "There should be 10 route configured") {
		return
	}
}

func TestCallListZones(t *testing.T) {
	controller, forwardService, err := getControllerAndForwardService()
	if !assert.NoError(t, err, fmt.Sprintf("could not initialize TestCallListZones: %s", err)) {
		return
	}

	responseContent := []*zones.Zone{
		{
			ID:               "example.com",
			Name:             "example.com",
			Type:             "",
			URL:              "/api/v1/servers/localhost/zones/example.com.",
			Kind:             "Native",
			RRSets:           nil,
			Serial:           2025041301,
			NotifiedSerial:   0,
			EditedSerial:     2025041301,
			Masters:          []string{},
			DNSSec:           false,
			NSEC3Param:       "",
			NSEC3Narrow:      false,
			Presigned:        false,
			SOAEdit:          "",
			SOAEditApi:       "",
			ApiRectify:       false,
			Zone:             "",
			Catalog:          "",
			Account:          "",
			Nameservers:      nil,
			MasterTSIGKeyIDs: nil,
			SlaveTSIGKeyIDs:  nil,
		},
	}

	responseJson, _ := json.Marshal(responseContent)

	body := io.NopCloser(
		bytes.NewBuffer(responseJson),
	)

	response := &nethttp.Response{
		StatusCode:    200,
		Body:          body,
		ContentLength: int64(len(responseJson)),
		Header:        nethttp.Header{},
	}

	forwardService.PushResponse(response)

	test.RunRequestSimple(t, controller, "/api", "/api/v1/servers/localhost/zones", "GET", 200)
}

func TestCallCreateZone(t *testing.T) {
	controller, err := getController()
	if !assert.NoError(t, err, fmt.Sprintf("could not initialize TestCallCreateZone: %s", err)) {
		return
	}

	test.RunRequest(t, controller, "/api", "/api/v1/servers/localhost/zones", "POST", 200)
}

func TestCallGetZone(t *testing.T) {
	controller, err := getController()
	if !assert.NoError(t, err, fmt.Sprintf("could not initialize TestCallGetZone: %s", err)) {
		return
	}

	test.RunRequest(t, controller, "/api", "/api/v1/servers/localhost/zones/example.com", "GET", 200)
}

func TestCallDeleteZone(t *testing.T) {
	controller, err := getController()
	if !assert.NoError(t, err, fmt.Sprintf("could not initialize TestCallDeleteZone: %s", err)) {
		return
	}

	test.RunRequest(t, controller, "/api", "/api/v1/servers/localhost/zones/example.com", "DELETE", 200)
}

func TestCallUpdateZone(t *testing.T) {
	controller, err := getController()
	if !assert.NoError(t, err, fmt.Sprintf("could not initialize TestCallUpdateZone: %s", err)) {
		return
	}

	test.RunRequest(t, controller, "/api", "/api/v1/servers/localhost/zones/example.com", "PUT", 200)
}

func TestCallAxfrRetrieve(t *testing.T) {
	controller, err := getController()
	if !assert.NoError(t, err, fmt.Sprintf("could not initialize TestCallAxfrRetrieve: %s", err)) {
		return
	}

	test.RunRequest(t, controller, "/api", "/api/v1/servers/localhost/zones/example.com/axfr-retrieve", "PUT", 200)
}

func TestCallNotifySlaves(t *testing.T) {
	controller, err := getController()
	if !assert.NoError(t, err, fmt.Sprintf("could not initialize TestCallNotifySlaves: %s", err)) {
		return
	}

	test.RunRequest(t, controller, "/api", "/api/v1/servers/localhost/zones/example.com/notify", "PUT", 200)
}

func TestCallExportZone(t *testing.T) {
	controller, err := getController()
	if !assert.NoError(t, err, fmt.Sprintf("could not initialize TestCallExportZone: %s", err)) {
		return
	}

	test.RunRequest(t, controller, "/api", "/api/v1/servers/localhost/zones/example.com/export", "GET", 200)
}

func TestCallUpdateRRSet(t *testing.T) {
	controller, err := getController()

	if !assert.NoError(t, err, fmt.Sprintf("could not initialize TestCallUpdateRRSet: %s", err)) {
		return
	}

	payload := &zones.RequestUpdateRRSet{
		RRSets: []*zones.RRSet{
			{
				Name:       "dev.example.com",
				Type:       "A",
				TTL:        3600,
				ChangeType: "REPLACE",
				Records: []zones.Record{
					{
						Content:  "127.0.0.1",
						Disabled: false,
					},
				},
			},
		},
	}

	test.RunRequest(t, controller, "/api", "/api/v1/servers/localhost/zones/example.com", "PATCH", 200, payload)
}

func TestCallRectifyZone(t *testing.T) {
	controller, err := getController()
	if !assert.NoError(t, err, fmt.Sprintf("could not initialize TestCallRectifyZone: %s", err)) {
		return
	}

	test.RunRequest(t, controller, "/api", "/api/v1/servers/localhost/zones/example.com/rectify", "PUT", 200)
}

func getController() (*zones.ZonesController, error) {
	container, err := setup.InitContainerTest(&setup.TestConfig{EnableMockAuthentication: true, EnableMockForwardService: true})
	if err != nil {
		return nil, err
	}

	return container.ZonesController.(*zones.ZonesController), nil
}

func getControllerAndForwardService() (*zones.ZonesController, *http.ForwardServiceMock, error) {
	container, err := setup.InitContainerTest(&setup.TestConfig{EnableMockAuthentication: true, EnableMockForwardService: true})
	if err != nil {
		return nil, nil, err
	}

	return container.ZonesController.(*zones.ZonesController), container.ForwardService.(*http.ForwardServiceMock), nil
}
