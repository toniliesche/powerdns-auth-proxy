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

import (
	"bytes"
	"fmt"
	"gopkg.in/fifo.v0"
	"io"
	"net/http"
	"powerdns-auth-proxy/domain/shared/basics"
	"powerdns-auth-proxy/domain/shared/interfaces"
)

type ForwardServiceMock struct {
	responseQueue *fifo.Queue[*http.Response]
}

func (s *ForwardServiceMock) ForwardRequest(request *http.Request) (*http.Response, error) {
	body := io.NopCloser(
		bytes.NewBufferString(fmt.Sprintf("{\"status\":200,\"method\":\"%s\",\"path\":\"%s\"}", request.Method, request.URL.Path)),
	)

	var response *http.Response
	var err error
	if s.responseQueue.Len() == 0 {
		response = &http.Response{
			StatusCode: 200,
			Body:       body,
		}
	} else {
		response, err = s.responseQueue.Dequeue()
		if err != nil {
			return nil, err
		}
	}

	response.Request = request

	return response, nil
}

func (s *ForwardServiceMock) PushResponse(response *http.Response) {
	err := s.responseQueue.Enqueue(response)
	if err != nil {
		return
	}
}

func ProvideForwardServiceMock(container *basics.InjectionContainer) (interfaces.ForwardServiceInterface, error) {
	return &ForwardServiceMock{
		responseQueue: fifo.New[*http.Response](10),
	}, nil
}
