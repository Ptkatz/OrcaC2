package rpce

import (
	"bytes"
	"encoding/binary"
)

type AlterContextReq struct {
	CommonHead              CommonHead
	MaxXmitFrag             uint16
	MaxRecvFrag             uint16
	AssocGroup              uint32
	PresentationContextList *PContextList
	AuthVerifier            *AuthVerifier //Optional
}

func NewAlterContextReq(callID, assgroup uint32, ctxList *PContextList, auth *AuthVerifier) AlterContextReq {
	r := AlterContextReq{}

	//todo, replace hard coded here with enum for bind
	r.CommonHead = NewCommonHeader(0x0e, 0x03, callID)
	r.MaxRecvFrag = 0x10b8 //minimum of 1432 (0x598) idk why, this is hard coded :grimmace:
	r.MaxXmitFrag = 0x10b8
	r.AssocGroup = assgroup
	r.PresentationContextList = ctxList
	r.AuthVerifier = auth

	r.CommonHead.FragLength = 24
	if r.PresentationContextList != nil {
		r.CommonHead.FragLength += uint16(r.PresentationContextList.SizeOf())
	}
	if r.AuthVerifier != nil {
		r.CommonHead.FragLength += uint16(r.AuthVerifier.SizeOf())
		r.CommonHead.AuthLength = uint16(len(r.AuthVerifier.AuthValue))
	}

	return r
}

func (b AlterContextReq) Bytes() []byte {
	buff := bytes.Buffer{}
	binary.Write(&buff, binary.LittleEndian, b.CommonHead)
	binary.Write(&buff, binary.LittleEndian, b.MaxXmitFrag)
	binary.Write(&buff, binary.LittleEndian, b.MaxRecvFrag)
	binary.Write(&buff, binary.LittleEndian, b.AssocGroup)
	if b.PresentationContextList != nil {
		binary.Write(&buff, binary.LittleEndian, b.PresentationContextList.Bytes())
	}
	if b.AuthVerifier != nil {
		binary.Write(&buff, binary.LittleEndian, b.AuthVerifier.Bytes())
	}
	return buff.Bytes()
}
