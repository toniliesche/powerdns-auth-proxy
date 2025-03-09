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

import nethttp "net/http"

type StatusMessageMapper struct {
	StatusCodeMapper *StatusCodeMapper
}

func (m *StatusMessageMapper) MapMessage(err error) string {
	switch m.StatusCodeMapper.MapStatusCode(err) {
	case nethttp.StatusOK:
		return "ok"
	case nethttp.StatusNotFound:
		return "not found"
	case nethttp.StatusBadRequest:
		return "bad request"
	case nethttp.StatusUnauthorized:
		return "unauthorized"
	case nethttp.StatusForbidden:
		return "forbidden"
	case nethttp.StatusConflict:
		return "conflict"
	default:
		return "internal error"
	}
}
