/*
 *   Copyright (c) 2022-2024 Intel Corporation
 *   All rights reserved.
 *   SPDX-License-Identifier: BSD-3-Clause
 */

package tpm

import (
	"crypto"

	"github.com/canonical/go-tpm2"
	"github.com/canonical/go-tpm2/mu"
	"github.com/pkg/errors"
)

func (tpm *trustedPlatformModule) ReadPublic(handle int) (crypto.PublicKey, []byte, []byte, error) {

	// verify the handle and load a resource context for the handle
	h := tpm2.Handle(handle)
	t := h.Type()
	if t != tpm2.HandleTypePersistent && t != tpm2.HandleTypeTransient {
		return nil, nil, nil, ErrInvalidHandle
	}

	if !tpm.ctx.DoesHandleExist(h) {
		return nil, nil, nil, ErrHandleDoesNotExist
	}

	handleContext, err := tpm.ctx.NewResourceContext(h)
	if err != nil {
		return nil, nil, nil, errors.Wrapf(err, "Failed to create resource context for handle 0x%x", handle)
	}

	public, _, qualifiedName, err := tpm.ctx.ReadPublic(handleContext, nil)
	if err != nil {
		return nil, nil, nil, err
	}

	tpmtPublicBytes, err := mu.MarshalToBytes(public)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "Failed to marshal tpmt public")
	}

	return public.Public(), tpmtPublicBytes, qualifiedName, nil
}

func (tpm *trustedPlatformModule) HandleExists(handle int) bool {
	return tpm.ctx.DoesHandleExist(tpm2.Handle(handle))
}
