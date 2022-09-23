package rpce

import (
	"bytes"
	"encoding/binary"
)

type Response struct {
	CommonHead   CommonHead
	AllocHint    uint32
	PContID      uint16
	CancelCount  byte
	Reserved     byte
	StubData     []byte
	AuthVerifier *AuthVerifier
}

func ParseResponse(b []byte) Response {
	r := Response{}
	br := bytes.NewReader(b)

	binary.Read(br, binary.LittleEndian, &r.CommonHead)
	binary.Read(br, binary.LittleEndian, &r.AllocHint)
	binary.Read(br, binary.LittleEndian, &r.PContID)
	binary.Read(br, binary.LittleEndian, &r.CancelCount)
	binary.Read(br, binary.LittleEndian, &r.Reserved)
	r.StubData = make([]byte, r.CommonHead.FragLength-24)
	br.Read(r.StubData)

	return r
}
