package rpce

import (
	"bytes"
	"encoding/binary"
)

//https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-rpce/a6b7b03c-4ac5-4c25-8c52-f2bec872ac97

type Auth3Req struct {
	CommonHead CommonHead
	Pad        [4]byte
	SecTrailer AuthVerifier
}

func NewAuth3Req(callID uint32, authLevel AuthLevel, authData []byte) Auth3Req {
	r := Auth3Req{}
	r.CommonHead = NewCommonHeader(0x10, 0x03, callID)
	r.CommonHead.FragLength = uint16(len(authData) + 28)
	r.CommonHead.AuthLength = uint16(len(authData))
	r.SecTrailer.AuthType = RPC_C_AUTHN_WINNT
	r.SecTrailer.AuthLevel = authLevel
	r.SecTrailer.AuthValue = make([]byte, len(authData))
	copy(r.SecTrailer.AuthValue, authData)

	return r
}

func (a Auth3Req) Bytes() []byte {
	buff := bytes.Buffer{}
	binary.Write(&buff, binary.LittleEndian, a.CommonHead)
	binary.Write(&buff, binary.LittleEndian, a.Pad)
	binary.Write(&buff, binary.LittleEndian, a.SecTrailer.Bytes())

	return buff.Bytes()
}
