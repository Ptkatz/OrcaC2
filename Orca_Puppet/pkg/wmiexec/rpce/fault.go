package rpce

import (
	"bytes"
	"encoding/binary"
)

//12.6.4.7

type FaultResp struct {
	CommonHead   CommonHead
	AllocHint    uint32
	PcontID      uint16
	CancelCount  byte
	Reserved     byte
	Status       uint32
	Reserved2    [4]byte
	StubData     []byte
	AuthVerifier AuthVerifier
}

func ParseFault(b []byte) FaultResp {
	r := FaultResp{}
	br := bytes.NewReader(b)
	binary.Read(br, binary.LittleEndian, &r.CommonHead)
	binary.Read(br, binary.LittleEndian, &r.AllocHint)
	binary.Read(br, binary.LittleEndian, &r.PcontID)
	binary.Read(br, binary.LittleEndian, &r.CancelCount)
	binary.Read(br, binary.LittleEndian, &r.Reserved)
	binary.Read(br, binary.LittleEndian, &r.Status)
	binary.Read(br, binary.LittleEndian, &r.Reserved2)

	//todo: marshal stub data here

	if r.CommonHead.AuthLength > 0 {
		binary.Read(br, binary.LittleEndian, &r.AuthVerifier.AuthType)
		binary.Read(br, binary.LittleEndian, &r.AuthVerifier.AuthLevel)
		binary.Read(br, binary.LittleEndian, &r.AuthVerifier.AuthPadLength)
		binary.Read(br, binary.LittleEndian, &r.AuthVerifier.Reserved)
		binary.Read(br, binary.LittleEndian, &r.AuthVerifier.ContextID)
		r.AuthVerifier.AuthValue = make([]byte, r.CommonHead.AuthLength)
		br.Read(r.AuthVerifier.AuthValue)
	}
	return r
}

func (pf FaultResp) StatusString() string {
	return statusmap[pf.Status]
}

const (
	AccessDenied uint32 = 5
)

var statusmap = map[uint32]string{
	5:          "nca_s_fault_access_denied",
	0x1c00001b: "nca_s_fault_remote_no_memory",
	0x1c01000b: "nca_proto_error (The RPC client or server protocol has been violated.)",
	0x1c010003: "nca_unk_if (The server does not export the requested interface.)",
}
