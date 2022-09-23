// Copyright (c) 2016 Andreas Auernhammer. All rights reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package skein

import "Orca_Puppet/pkg/go-engine/crypto/cryptonight/inter/skein/threefish"

const (
	// The blocksize of Skein-512 in bytes.
	BlockSize = threefish.BlockSize512
)

// The different parameter types
const (
	// CfgKey is the config type for the Key.
	CfgKey uint64 = 0

	// CfgConfig is the config type for the configuration.
	CfgConfig uint64 = 4

	// CfgPersonal is the config type for the personalization.
	CfgPersonal uint64 = 8

	// CfgPublicKey is the config type for the public key.
	CfgPublicKey uint64 = 12

	// CfgKeyID is the config type for the key id.
	CfgKeyID uint64 = 16

	// CfgNonce is the config type for the nonce.
	CfgNonce uint64 = 20

	// CfgMessage is the config type for the message.
	CfgMessage uint64 = 48

	// CfgOutput is the config type for the output.
	CfgOutput uint64 = 63

	// FirstBlock is the first block flag
	FirstBlock uint64 = 1 << 62

	// FinalBlock is the final block flag
	FinalBlock uint64 = 1 << 63

	// The skein schema ID = S H A 3 1 0 0 0
	SchemaID uint64 = 0x133414853
)

// Precomputed chain values for Skein-512
var iv160 = [9]uint64{
	0x28B81A2AE013BD91, 0xC2F11668B5BDF78F, 0x1760D8F3F6A56F12, 0x4FB747588239904F,
	0x21EDE07F7EAF5056, 0xD908922E63ED70B8, 0xB8EC76FFECCB52FA, 0x01A47BB8A3F27A6E,
	0,
}

var iv256 = [9]uint64{
	0xCCD044A12FDB3E13, 0xE83590301A79A9EB, 0x55AEA0614F816E6F, 0x2A2767A4AE9B94DB,
	0xEC06025E74DD7683, 0xE7A436CDC4746251, 0xC36FBAF9393AD185, 0x3EEDBA1833EDFC13,
	0,
}

var iv384 = [9]uint64{
	0xA3F6C6BF3A75EF5F, 0xB0FEF9CCFD84FAA4, 0x9D77DD663D770CFE, 0xD798CBF3B468FDDA,
	0x1BC4A6668A0E4465, 0x7ED7D434E5807407, 0x548FC1ACD4EC44D6, 0x266E17546AA18FF8,
	0,
}

var iv512 = [9]uint64{
	0x4903ADFF749C51CE, 0x0D95DE399746DF03, 0x8FD1934127C79BCE, 0x9A255629FF352CB1,
	0x5DB62599DF6CA7B0, 0xEABE394CA9D5C3F4, 0x991112C71A75B523, 0xAE18A40B660FCC33,
	0,
}
