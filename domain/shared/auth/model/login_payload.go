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

package model

import (
	"fmt"
	"powerdns-auth-proxy/domain/shared/http/errors"
	"strings"
)

type LoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (p *LoginPayload) Verify() errors.HTTPError {
	if strings.TrimSpace(p.Username) == "" {
		return errors.NewBadRequestError(fmt.Errorf("username is required"))
	}

	if strings.TrimSpace(p.Password) == "" {
		return errors.NewBadRequestError(fmt.Errorf("password is required"))
	}

	return nil
}
