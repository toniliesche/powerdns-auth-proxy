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

package users

import (
	"github.com/gin-gonic/gin"
	"powerdns-auth-proxy/domain/shared/auth"
	"powerdns-auth-proxy/domain/shared/model/management"
)

func (c *UserController) updateUser(context *gin.Context) {
	payload := &management.UserUpdatePayload{}
	if !c.ParsePayload(context, payload) {
		return
	}

	var userId uint
	var ok bool
	if userId, ok = c.GetIntParameter(context, "id"); !ok {
		return
	}

	isAdmin, err := c.userService.CheckIsAdminUser(userId)
	if err != nil {
		c.HandleError(context, err)
		return
	}

	if isAdmin && !c.CheckAccess(context, auth.Superadmin, "") {
		c.ForbiddenError(context)
		return
	}

	err = c.userService.UpdateUserPasswordById(userId, payload.Password)
	if err != nil {
		c.HandleError(context, err)
		return
	}

	c.HandleSuccess(context)
}
