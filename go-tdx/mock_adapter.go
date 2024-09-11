//go:build test

/*
 *   Copyright (c) 2022-2024 Intel Corporation
 *   All rights reserved.
 *   SPDX-License-Identifier: BSD-3-Clause
 */
package tdx

import (
	"github.com/intel/trustauthority-client/go-connector"
)

type mockAdapter struct {
	uData       []byte
	EvLogParser EventLogParser
}

func NewTdxAdapter(udata []byte, evLogParser EventLogParser) (connector.EvidenceAdapter, error) {
	return &mockAdapter{
		uData:       udata,
		EvLogParser: evLogParser,
	}, nil
}

func NewAzureTdxAdapter(udata []byte) (connector.EvidenceAdapter, error) {
	return &mockAdapter{
		uData: udata,
	}, nil
}

func (adapter *mockAdapter) CollectEvidence(nonce []byte) (*connector.Evidence, error) {

	return &connector.Evidence{
		Type:     1,
		Evidence: nil,
		UserData: nil,
		EventLog: nil,
	}, nil
}
