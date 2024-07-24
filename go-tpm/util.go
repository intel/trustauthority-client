/*
 *   Copyright (c) 2022-2024 Intel Corporation
 *   All rights reserved.
 *   SPDX-License-Identifier: BSD-3-Clause
 */
package tpm

import (
	"crypto"
	"sort"
	"strconv"
	"strings"

	"github.com/canonical/go-tpm2"
	"github.com/pkg/errors"
)

// toTpm2PcrSelectionList takes in an array of PcrSelection structs and converts them to a
// tpm2.PCRSelectionList.  If no selection is provided, the defaultPcrSelections is used.
func toTpm2PcrSelectionList(selection ...PcrSelection) (tpm2.PCRSelectionList, error) {
	var selectedPcrs []PcrSelection
	var pcrSelectionList tpm2.PCRSelectionList

	if len(selection) == 0 {
		selectedPcrs = defaultPcrSelections
	} else {
		selectedPcrs = selection
		for _, s := range selectedPcrs {
			sort.Ints(s.Pcrs)
		}
	}

	for _, selected := range selectedPcrs {
		var algorithmId tpm2.HashAlgorithmId
		switch selected.Hash {
		case crypto.SHA1:
			algorithmId = tpm2.HashAlgorithmSHA1
		case crypto.SHA256:
			algorithmId = tpm2.HashAlgorithmSHA256
		case crypto.SHA384:
			algorithmId = tpm2.HashAlgorithmSHA384
		case crypto.SHA512:
			algorithmId = tpm2.HashAlgorithmSHA512
		default:
			return nil, errors.Errorf("Unsupported hash algorithm: %v", selected.Hash)
		}

		pcrSelectionList = append(pcrSelectionList, tpm2.PCRSelection{
			Hash:   algorithmId,
			Select: selected.Pcrs,
		})
	}

	return pcrSelectionList, nil
}

func parsePcrSelections(args string) ([]PcrSelection, error) {
	pcrSelections := []PcrSelection{}

	if args == "" {
		return pcrSelections, nil
	}

	// tpm2-tools like selection strings "sha1:1,2,3+sha256:1,2,3"
	selections := strings.Split(args, "+")
	for _, selection := range selections {
		pcrSelection := PcrSelection{}

		// Split the selection
		parts := strings.Split(selection, ":")
		if len(parts) != 2 {
			return nil, errors.New("invalid format")
		}

		hash := parts[0]
		switch hash {
		case "sha1":
			pcrSelection.Hash = crypto.SHA1
		case "sha256":
			pcrSelection.Hash = crypto.SHA256
		case "sha384":
			pcrSelection.Hash = crypto.SHA384
		case "sha512":
			pcrSelection.Hash = crypto.SHA512
		default:
			return nil, errors.Errorf("Invalid PCR hash %q", hash)
		}

		// Parse the array of pcr banks
		intsStr := parts[1]
		banks := strings.Split(intsStr, ",")
		for _, str := range banks {
			bank, err := strconv.Atoi(str)
			if err != nil {
				return nil, errors.Errorf("Failed to parse PCR bank %q", str)
			}
			if bank < 0 || bank > 23 {
				return nil, errors.Errorf("Bank %d out of range", bank)
			}
			pcrSelection.Pcrs = append(pcrSelection.Pcrs, bank)
		}

		pcrSelections = append(pcrSelections, pcrSelection)
	}

	return pcrSelections, nil
}