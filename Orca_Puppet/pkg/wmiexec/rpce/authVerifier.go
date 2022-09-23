package rpce

import (
	"bytes"
	"encoding/binary"
)

// 13.2.6.1

type AuthVerifier struct {
	Align         []byte // must be 4-byte aligned!!
	AuthType      SecurityProviders
	AuthLevel     AuthLevel
	AuthPadLength byte
	Reserved      byte
	ContextID     uint32
	AuthValue     []byte
}

func (a AuthVerifier) SizeOf() int {
	return len(a.Bytes())
}

func (a AuthVerifier) ValueSize() int {
	return len(a.AuthValue)
}

func (a AuthVerifier) Bytes() []byte {
	buff := bytes.Buffer{}

	binary.Write(&buff, binary.LittleEndian, a.Align)
	binary.Write(&buff, binary.LittleEndian, a.AuthType)
	binary.Write(&buff, binary.LittleEndian, a.AuthLevel)
	binary.Write(&buff, binary.LittleEndian, a.AuthPadLength)
	binary.Write(&buff, binary.LittleEndian, a.Reserved)
	binary.Write(&buff, binary.LittleEndian, a.ContextID)
	binary.Write(&buff, binary.LittleEndian, a.AuthValue)

	return buff.Bytes()
}

func NewAuthVerifier(authType SecurityProviders, authLevel AuthLevel, contextID uint32, value []byte) AuthVerifier {
	r := AuthVerifier{
		AuthType:  authType,
		AuthLevel: authLevel,
		ContextID: contextID,
		AuthValue: value,
	}

	return r
}

func (a *AuthVerifier) UpdatePadding(v int) {
	a.Align = make([]byte, v)
	a.AuthPadLength = byte(v)
}
