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

import "fmt"

type DomainRole struct {
	Domain string `json:"domain"`
	Role   string `json:"role"`
}

func (r *DomainRole) Verify() error {
	if r.Domain == "" {
		return fmt.Errorf("domain is required")
	}

	if r.Role == "" {
		return fmt.Errorf("role is required")
	}

	return nil
}
