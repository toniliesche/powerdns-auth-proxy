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

package management

import (
	"fmt"
	"powerdns-auth-proxy/domain/shared/http/errors"
)

type RoleCreatePayload struct {
	Name string `json:"name"`
}

func (p *RoleCreatePayload) Verify() errors.HTTPError {
	if p.Name == "" {
		return errors.NewBadRequestError(fmt.Errorf("name is required"))
	}

	return nil
}
