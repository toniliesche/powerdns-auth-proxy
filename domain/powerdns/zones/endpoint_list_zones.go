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

package zones

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"powerdns-auth-proxy/domain/shared/auth"
)

func (c *ZonesController) listZones(context *gin.Context) {
	response, err := c.ForwardRequest(context.Request)
	if err != nil {
		context.JSON(500, gin.H{"error": "failed to forward request"})
		return
	}
	responseBody := response.Body
	defer responseBody.Close()

	if response.StatusCode != 200 {
		c.WriteResponse(context, response)
		return
	}

	body, err := io.ReadAll(responseBody)
	if err != nil {
		context.JSON(500, gin.H{"error": "failed to read response body"})
		return
	}

	var zones []*Zone
	err = json.Unmarshal(body, &zones)
	if err != nil {
		context.JSON(500, gin.H{"error": "failed to unmarshal response body"})
		return
	}

	filteredZones := make([]*Zone, 0, len(zones))
	for _, zone := range zones {
		if c.CheckAccessOnResource(context, auth.Reader, zone.Name) {
			filteredZones = append(filteredZones, zone)
		}
	}

	if len(filteredZones) == 0 {
		context.JSON(200, gin.H{"message": "no zones found"})
		return
	}

	filteredBody, err := json.Marshal(filteredZones)
	if err != nil {
		context.JSON(500, gin.H{"error": "failed to marshal response body"})
		return
	}

	context.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", len(filteredBody)))
	response.Body = io.NopCloser(bytes.NewBuffer(filteredBody))
	c.WriteResponseModified(context, response)
}
