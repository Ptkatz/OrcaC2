package rpce

import (
	"bytes"
	"encoding/binary"
)

type RequestReq struct { //poorly named, I'm aware
	CommonHead   CommonHead
	AllocHint    uint32
	PContextID   uint16
	Opnum        uint16
	StubData     []byte
	AuthVerifier *AuthVerifier
}

func NewRequestReq(callID uint32, ctxID uint16, opNum uint16, data []byte, auth *AuthVerifier) RequestReq {
	r := RequestReq{}
	//todo, don't hard code ptype
	r.CommonHead = NewCommonHeader(0, 0x03, callID)
	r.PContextID = ctxID
	r.Opnum = opNum
	r.StubData = make([]byte, len(data))
	copy(r.StubData, data)
	r.AuthVerifier = auth

	r.AllocHint = 0 //idfk, I guess this should be the length of the data segment?

	r.UpdateLengths()

	return r
}

//UpdateLengths sets the commonheader.fraglength and .authlength values based on the current state of the object.
func (r *RequestReq) UpdateLengths() {
	r.CommonHead.FragLength = 24 //length of common header
	r.CommonHead.FragLength += uint16(len(r.StubData))

	if r.AuthVerifier != nil {
		r.AuthVerifier.UpdatePadding(len(r.StubData) % 4)
		r.CommonHead.FragLength += uint16(r.AuthVerifier.SizeOf())
		r.CommonHead.AuthLength = uint16(len(r.AuthVerifier.AuthValue))
	}
}

func (r RequestReq) AuthBytes() []byte {

	buff := bytes.Buffer{}
	binary.Write(&buff, binary.LittleEndian, r.CommonHead)
	binary.Write(&buff, binary.LittleEndian, r.AllocHint)
	binary.Write(&buff, binary.LittleEndian, r.PContextID)
	binary.Write(&buff, binary.LittleEndian, r.Opnum)
	buff.Write(r.StubData)
	if r.AuthVerifier != nil {
		//this is a bit gross - the value is where the verifier gets put, and doesn't get included in the hashing scheme (obviously)
		//but the length needs to be included in the common headers. Doing this makes sure if a value is set in the auth object, it's not returned as part of the request object.
		buff.Write(r.AuthVerifier.Bytes()[:r.AuthVerifier.SizeOf()-r.AuthVerifier.ValueSize()])
	}

	return buff.Bytes()
}

func (r RequestReq) Bytes() []byte {
	buff := bytes.Buffer{}

	binary.Write(&buff, binary.LittleEndian, r.CommonHead)
	binary.Write(&buff, binary.LittleEndian, r.AllocHint)
	binary.Write(&buff, binary.LittleEndian, r.PContextID)
	binary.Write(&buff, binary.LittleEndian, r.Opnum)
	buff.Write(r.StubData)
	if r.AuthVerifier != nil {
		buff.Write(r.AuthVerifier.Bytes())
	}
	return buff.Bytes()
}
