package rpce

import (
	"bytes"
	"encoding/binary"
)

//12.6.4.4

type BindAckResp struct {
	CommonHead   CommonHead
	MaxXmitFrag  uint16
	MaxRecvFrag  uint16
	AssocGroupID uint32
	SecAddr      PortAny
	Pad2         []byte //idk what happened to pad1 ?
	PResultList  PResultList
	AuthVerifier AuthVerifier
}

func ParseBindAck(b []byte) BindAckResp {
	r := BindAckResp{}
	br := bytes.NewReader(b)
	l := br.Len()
	binary.Read(br, binary.LittleEndian, &r.CommonHead)
	binary.Read(br, binary.LittleEndian, &r.MaxXmitFrag)
	binary.Read(br, binary.LittleEndian, &r.MaxRecvFrag)
	binary.Read(br, binary.LittleEndian, &r.AssocGroupID)
	binary.Read(br, binary.LittleEndian, &r.SecAddr.Length)
	r.SecAddr.PortSpec = make([]byte, r.SecAddr.Length)
	br.Read(r.SecAddr.PortSpec)
	r.Pad2 = make([]byte, (l-br.Len())%4) //align to 4 byte.. this seems jank
	br.Read(r.Pad2)
	binary.Read(br, binary.LittleEndian, &r.PResultList.NResults)
	binary.Read(br, binary.LittleEndian, &r.PResultList.Reserved)
	binary.Read(br, binary.LittleEndian, &r.PResultList.Reserved2)
	r.PResultList.PResults = []PResult{}
	for i := byte(0); i < r.PResultList.NResults; i++ {
		tmp := PResult{}
		binary.Read(br, binary.LittleEndian, &tmp)
		r.PResultList.PResults = append(r.PResultList.PResults, tmp)
	}
	//PLACEHOLDER FOR INEVITABLE PAD WRANGLING
	//binary.Read(br, binary.LittleEndian, &r.AuthVerifier.Align)
	//
	binary.Read(br, binary.LittleEndian, &r.AuthVerifier.AuthType)
	binary.Read(br, binary.LittleEndian, &r.AuthVerifier.AuthLevel)
	binary.Read(br, binary.LittleEndian, &r.AuthVerifier.AuthPadLength)
	binary.Read(br, binary.LittleEndian, &r.AuthVerifier.Reserved)
	binary.Read(br, binary.LittleEndian, &r.AuthVerifier.ContextID)
	r.AuthVerifier.AuthValue = make([]byte, r.CommonHead.AuthLength)
	br.Read(r.AuthVerifier.AuthValue)

	return r
}
