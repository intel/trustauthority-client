/*
 *   Copyright (c) 2022 Intel Corporation
 *   All rights reserved.
 *   SPDX-License-Identifier: BSD-3-Clause
 */

package constants

const (
	MaxKeyLen              = 256
	PemBlockTypePrivateKey = "PRIVATE KEY"
	PemBlockTypePubliceKey = "PUBLIC KEY"
	RSAKeyBitLength        = 3072
	CLIShortDescription    = "Amber Attestation Client for TDX"
)

// Command Names
const (
	CreateKeyPairCmd = "create-key-pair"
	DecryptCmd       = "decrypt"
	QuoteCmd         = "quote"
	TokenCmd         = "token"
	RootCmd          = "amber-cli"
	VersionCmd       = "version"
	VerifyCmd        = "verify"
)

// Options Names
const (
	PrivateKeyPathOption = "key-path"
	PublicKeyPathOption  = "pub-path"
	PrivateKeyOption     = "key"
	PolicyIdsOption      = "policy-ids"
	InputOption          = "in"
	UserDataOption       = "user-data"
	NonceOption          = "nonce"
	VerifyTokenOption    = "token-path"
)

const (
	AmberApiKeyEnv = "AMBER_API_KEY"
	AmberUrlEnv    = "AMBER_URL"
)
