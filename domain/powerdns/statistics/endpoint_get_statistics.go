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

package statistics

import (
	"github.com/gin-gonic/gin"
	"powerdns-auth-proxy/domain/shared/auth"
)

func (c *StatisticsController) getStatistics(context *gin.Context) {
	if !c.CheckAccessOnResource(context, auth.Admin, "") {
		c.ForbiddenError(context)
		return
	}

	response, err := c.ForwardRequest(context.Request)
	if err != nil {
		c.HandleError(context, err)
		return
	}

	c.WriteResponse(context, response)
}
