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

package test

import "strconv"

type Registry struct {
	Data map[string]string
}

func (r *Registry) Set(key, value string) {
	r.Data[key] = value
}

func (r *Registry) Get(key string) string {
	return r.Data[key]
}

func (r *Registry) GetUint(s string) uint {
	id := r.Data[s]

	intVal, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return 0
	}

	return uint(intVal)
}

func NewRegistry() *Registry {
	return &Registry{Data: make(map[string]string)}
}
