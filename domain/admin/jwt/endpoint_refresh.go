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

package jwt

import "github.com/gin-gonic/gin"

func (c *JWTController) refresh(context *gin.Context) {
	tokenResponse, err := c.jwtService.Refresh(context)
	if err != nil {
		c.HandleError(context, err)
		return
	}

	c.HandleData(context, tokenResponse)
}
