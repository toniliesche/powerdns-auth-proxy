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

package zones

type Zone struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	Type             string   `json:"type"`
	URL              string   `json:"url"`
	Kind             string   `json:"kind"`
	RRSets           []*RRSet `json:"rrsets"`
	Serial           int      `json:"serial"`
	NotifiedSerial   int      `json:"notified_serial"`
	EditedSerial     int      `json:"edited_serial"`
	Masters          []string `json:"masters"`
	DNSSec           bool     `json:"dnssec"`
	NSEC3Param       string   `json:"nsec3param"`
	NSEC3Narrow      bool     `json:"nsec3narrow"`
	Presigned        bool     `json:"presigned"`
	SOAEdit          string   `json:"soa_edit"`
	SOAEditApi       string   `json:"soa_edit_api"`
	ApiRectify       bool     `json:"api_rectify"`
	Zone             string   `json:"zone"`
	Catalog          string   `json:"catalog"`
	Account          string   `json:"account"`
	Nameservers      []string `json:"nameservers"`
	MasterTSIGKeyIDs []string `json:"master_tsig_key_ids"`
	SlaveTSIGKeyIDs  []string `json:"slave_tsig_key_ids"`
}
