package generateopt

import (
	"bytes"
)

const xorKey = 24

func ReplaceBytes(data, sBytes, dBytes []byte) []byte {
	dBytesArr := make([]byte, len(sBytes))
	copy(dBytesArr, dBytes)
	data = bytes.Replace(data, sBytes, dBytesArr, -1)
	return data
}

func DoXor(sBytes []byte) []byte {
	for i, _ := range sBytes {
		sBytes[i] = sBytes[i] ^ xorKey
	}
	return sBytes
}
