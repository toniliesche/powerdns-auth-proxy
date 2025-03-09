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

package cryptokeys

type CryptoKey struct {
	Type       string   `json:"type"`
	ID         int      `json:"id"`
	KeyType    string   `json:"keytype"`
	Active     bool     `json:"active"`
	Published  bool     `json:"published"`
	DNSKey     string   `json:"dnskey"`
	DS         []string `json:"ds"`
	CDS        []string `json:"cds"`
	PrivateKey string   `json:"privatekey"`
	Algorithm  string   `json:"algorithm"`
	Bits       int      `json:"bits"`
}
