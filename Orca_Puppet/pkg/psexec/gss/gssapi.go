package gss

// GSS-API/SPNEGO支持 遵循RFC-2713、RFC-4178标准

import (
	"Orca_Puppet/pkg/psexec/encoder"
	"Orca_Puppet/pkg/psexec/ntlm"
	"encoding/asn1"
	"log"
	"strconv"
	"strings"
)

const SPNEGOOID = "1.3.6.1.5.5.2"

const GssStateAcceptCompleted = 0
const GssStateAcceptIncomplete = 1
const GssStateReject = 2
const GssStateRequestMic = 3

// https://docs.microsoft.com/zh-cn/openspecs/windows_protocols/ms-spng/8e71cf53-e867-4b79-b5b5-38c92be3d472

type NegTokenInitData struct {
	MechTypes    []asn1.ObjectIdentifier `asn1:"explicit,tag:0"`                    //身份验证列表，遵循RFC4178标准
	ReqFlags     asn1.BitString          `asn1:"explicit,optional,omitempty,tag:1"` //上下文标志，发送方忽略该标志
	MechToken    []byte                  `asn1:"explicit,optional,omitempty,tag:2"` //令牌，8位可选字符串
	MechTokenMIC []byte                  `asn1:"explicit,optional,omitempty,tag:3"` //消息完整性令牌
}

type NegTokenInit struct {
	OID  asn1.ObjectIdentifier //SPNEGO GSSAPI协商安全机制
	Data NegTokenInitData      `asn1:"explicit"` //NegTokenInit的扩展NegTokenInit2
}

type NegTokenResp struct {
	// RFC4178 4.2.2
	NegState      asn1.Enumerated       `asn1:"explicit,optional,omitempty,tag:0"`
	SupportedMech asn1.ObjectIdentifier `asn1:"explicit,optional,omitempty,tag:1"`
	ResponseToken []byte                `asn1:"explicit,optional,omitempty,tag:2"`
	MechListMIC   []byte                `asn1:"explicit,optional,omitempty,tag:3"`
}

// SPNEGO部分的asn1数据处理
// gsswrapped used to force ASN1 encoding to include explicit sequence tags
// Type does not fulfill the BinaryMarshallable interfce and is used only as a
// helper to marshal a NegTokenResp
type gsswrapped struct{ G interface{} }

func ObjectIDStrToInt(oid string) ([]int, error) {
	ret := []int{}
	tokens := strings.Split(oid, ".")
	for _, token := range tokens {
		i, err := strconv.Atoi(token)
		if err != nil {
			return nil, err
		}
		ret = append(ret, i)
	}
	return ret, nil
}

func NewNegTokenInit() (NegTokenInit, error) {
	oid, err := ObjectIDStrToInt(SPNEGOOID)
	if err != nil {
		return NegTokenInit{}, err
	}
	ntlmoid, err := ObjectIDStrToInt(ntlm.NTLMSSPMECHTYPEOID)
	if err != nil {
		return NegTokenInit{}, err
	}
	return NegTokenInit{
		OID: oid,
		Data: NegTokenInitData{
			MechTypes:    []asn1.ObjectIdentifier{ntlmoid},
			ReqFlags:     asn1.BitString{},
			MechToken:    []byte{},
			MechTokenMIC: []byte{},
		},
	}, nil
}

func NewNegTokenResp() (NegTokenResp, error) {
	return NegTokenResp{}, nil
}

func (n *NegTokenInit) MarshalBinary(meta *encoder.Metadata) ([]byte, error) {
	buf, err := asn1.Marshal(*n)
	if err != nil {
		log.Panicln(err)
		return nil, err
	}

	// When marshalling struct, asn1 uses 30 (sequence) tag by default.
	// Override to set 60 (application) to remain consistent with GSS/SMB
	buf[0] = 0x60
	return buf, nil
}

func (n *NegTokenInit) UnmarshalBinary(buf []byte, meta *encoder.Metadata) error {
	data := NegTokenInit{}
	if _, err := asn1.UnmarshalWithParams(buf, &data, "application"); err != nil {
		return err
	}
	*n = data
	return nil
}

func (r *NegTokenResp) MarshalBinary(meta *encoder.Metadata) ([]byte, error) {
	// Oddities in Go's ASN1 package vs SMB encoding mean we have to wrap our
	// struct in another struct to ensure proper tags and lengths are added
	// to encoded data
	wrapped := &gsswrapped{*r}
	return wrapped.MarshalBinary(meta)
}

func (r *NegTokenResp) UnmarshalBinary(buf []byte, meta *encoder.Metadata) error {
	data := NegTokenResp{}
	if _, err := asn1.UnmarshalWithParams(buf, &data, "explicit,tag:1"); err != nil {
		return err
	}
	*r = data
	return nil
}

func (g *gsswrapped) MarshalBinary(meta *encoder.Metadata) ([]byte, error) {
	buf, err := asn1.Marshal(*g)
	if err != nil {
		return nil, err
	}
	buf[0] = 0xa1
	return buf, nil
}
