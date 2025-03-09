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
	"strings"
)

type SessionLogoutPayload struct {
	SessionIDs []string `json:"session_ids"`
}

func (p *SessionLogoutPayload) Verify() errors.HTTPError {
	if len(p.SessionIDs) == 0 {
		return errors.NewBadRequestError(fmt.Errorf("session_ids is required"))
	}

	for _, sessionID := range p.SessionIDs {
		if strings.TrimSpace(sessionID) == "" {
			return errors.NewBadRequestError(fmt.Errorf("session_id cannot be empty"))
		}
	}

	return nil
}
