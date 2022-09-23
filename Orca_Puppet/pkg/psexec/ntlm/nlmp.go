package ntlm

// https://docs.microsoft.com/zh-cn/openspecs/windows_protocols/ms-nlmp/
// 认证加密实现

import (
	"Orca_Puppet/pkg/psexec/encoder"
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"golang.org/x/crypto/md4"
	"hash"
	"strings"
)

// ssp安全签名
const NTLMSecSignature = "NTLMSSP\x00"

// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-nlmp/5e550938-91d4-459f-b67d-75d70009e3f3

// ntlm对象标识符
const NTLMSSPMECHTYPEOID = "1.3.6.1.4.1.311.2.2.10"

// ntlm协议头类型
const (
	NTLMNegotiate    = 0x00000001
	NTLMChallenge    = 0x00000002
	NTLMAuthenticate = 0x00000003
)

const (
	FlgNegUnicode uint32 = 1 << iota
	FlgNegOEM
	FlgNegRequestTarget
	FlgNegReserved10
	FlgNegSign
	FlgNegSeal
	FlgNegDatagram
	FlgNegLmKey
	FlgNegReserved9
	FlgNegNTLM
	FlgNegReserved8
	FlgNegAnonymous
	FlgNegOEMDomainSupplied
	FlgNegOEMWorkstationSupplied
	FlgNegReserved7
	FlgNegAlwaysSign
	FlgNegTargetTypeDomain
	FlgNegTargetTypeServer
	FlgNegReserved6
	FlgNegExtendedSessionSecurity
	FlgNegIdentify
	FlgNegReserved5
	FlgNegRequestNonNtSessionKey
	FlgNegTargetInfo
	FlgNegReserved4
	FlgNegVersion
	FlgNegReserved3
	FlgNegReserved2
	FlgNegReserved1
	FlgNeg128
	FlgNegKeyExch
	FlgNeg56
)

const (
	MsvAvEOL uint16 = iota
	MsvAvNbComputerName
	MsvAvNbDomainName
	MsvAvDnsComputerName
	MsvAvDnsDomainName
	MsvAvDnsTreeName
	MsvAvFlags
	MsvAvTimestamp
	MsvAvSingleHost
	MsvAvTargetName
	MsvChannelBindings
)

// https://docs.microsoft.com/zh-cn/openspecs/windows_protocols/ms-nlmp/464551a8-9fc4-428e-b3d3-bc5bfb2e73a5

// NTLMv1 认证
// Define NTOWFv1(Passwd, User, UserDom) as MD4(UNICODE(Passwd))
func NTOWFv1(pass string) []byte {
	hash := md4.New()
	hash.Write(encoder.ToUnicode(pass))
	return hash.Sum(nil)
}

// https://docs.microsoft.com/zh-cn/openspecs/windows_protocols/ms-nlmp/5e550938-91d4-459f-b67d-75d70009e3f3

// NTLMv2 认证
// Define NTOWFv2(Passwd, User, UserDom) as HMAC_MD5(
// MD4(UNICODE(Passwd)), UNICODE(ConcatenationOf( Uppercase(User),
// UserDom ) ) )
func NTOWFv2(password, user, userDomain string) []byte {
	h := hmac.New(md5.New, NTOWFv1(password))
	h.Write(encoder.ToUnicode(strings.ToUpper(user) + userDomain))
	return h.Sum(nil)
}

// NTLMv2 hash认证
func NTOWFv2Hash(hash, user, userDomain string) []byte {
	Hash, _ := hex.DecodeString(hash)
	hm := hmac.New(md5.New, Hash)
	hm.Write(encoder.ToUnicode(strings.ToUpper(user) + userDomain))
	return hm.Sum(nil)
}

// Define LMOWFv2(Passwd, User, UserDom) as NTOWFv2(Passwd, User,
// UserDom)
func LMOWFv2(password, user, userDomain string) []byte {
	return NTOWFv2(password, user, userDomain)
}

// 计算ntlmv2响应
// Set temp to ConcatenationOf(Responserversion, HiResponserversion,
//     Z(6), Time, ClientChallenge, Z(4), ServerName, Z(4))
// Set NTProofStr to HMAC_MD5(ResponseKeyNT,
//     ConcatenationOf(CHALLENGE_MESSAGE.ServerChallenge,temp))
// Set NtChallengeResponse to ConcatenationOf(NTProofStr, temp)
// Set LmChallengeResponse to ConcatenationOf(HMAC_MD5(ResponseKeyLM,
//     ConcatenationOf(CHALLENGE_MESSAGE.ServerChallenge, ClientChallenge)),
//     ClientChallenge )
func ComputeNTLMv2Response(h hash.Hash, clientChallenge, serverChallenge, timestamp, serverName []byte) (NTChallengeResponse, LMChallengeResponse, SessionBaseKey []byte) {
	temp := []byte{1, 1}
	temp = append(temp, 0, 0, 0, 0, 0, 0)
	temp = append(temp, timestamp...)
	temp = append(temp, clientChallenge...)
	temp = append(temp, 0, 0, 0, 0)
	temp = append(temp, serverName...)
	temp = append(temp, 0, 0, 0, 0)
	// 计算NT response
	h.Write(append(serverChallenge, temp...))
	hmacNT := h.Sum(nil)
	// 计算LM response
	h.Write(append(serverChallenge, clientChallenge...))
	hmacLM := h.Sum(nil)
	// 计算Session Key
	// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-nlmp/5e550938-91d4-459f-b67d-75d70009e3f3
	// Key = HMAC_MD5(NTLMv2 Hash, HMAC_MD5(NTLMv2 Hash, NTLMv2 Response + Challenge))
	// Set SessionBaseKey to HMAC_MD5(ResponseKeyNT, NTProofStr)
	h.Write(append(hmacNT, temp...))
	sessionBaseKey := h.Sum(nil)
	return append(hmacNT, temp...), append(hmacLM, clientChallenge...), sessionBaseKey
}

// 服务器响应检查
type AvPair struct {
	AvID  uint16
	AvLen uint16 `smb:"len:Value"`
	Value []byte
}

type AvPairSlice []AvPair

func (p AvPair) Size() uint64 {
	return uint64(binary.Size(p.AvID) + binary.Size(p.AvLen) + int(p.AvLen))
}

func (s *AvPairSlice) MarshalBinary(meta *encoder.Metadata) ([]byte, error) {
	var ret []byte
	w := bytes.NewBuffer(ret)
	for _, pair := range *s {
		buf, err := encoder.Marshal(pair)
		if err != nil {
			return nil, err
		}
		if err := binary.Write(w, binary.LittleEndian, buf); err != nil {
			return nil, err
		}
	}
	return w.Bytes(), nil
}

func (s *AvPairSlice) UnmarshalBinary(buf []byte, meta *encoder.Metadata) error {
	slice := []AvPair{}
	l, ok := meta.Lens[meta.CurrField]
	if !ok {
		return errors.New(fmt.Sprintf("Cannot unmarshal field '%s'. Missing length\n", meta.CurrField))
	}
	o, ok := meta.Offsets[meta.CurrField]
	if !ok {
		return errors.New(fmt.Sprintf("Cannot unmarshal field '%s'. Missing offset\n", meta.CurrField))
	}
	for i := l; i > 0; {
		var avPair AvPair
		err := encoder.Unmarshal(meta.ParentBuf[o:o+i], &avPair)
		if err != nil {
			return err
		}
		slice = append(slice, avPair)
		size := avPair.Size()
		o += size
		i -= size
	}
	*s = slice
	return nil
}

// 通用头
type Header struct {
	Signature   []byte `smb:"fixed:8"`
	MessageType uint32
}

type Negotiate struct {
	Header
	NegotiateFlags          uint32
	DomainNameLen           uint16 `smb:"len:DomainName"`
	DomainNameMaxLen        uint16 `smb:"len:DomainName"`
	DomainNameBufferOffset  uint32 `smb:"offset:DomainName"`
	WorkstationLen          uint16 `smb:"len:Workstation"`
	WorkstationMaxLen       uint16 `smb:"len:Workstation"`
	WorkstationBufferOffset uint32 `smb:"offset:Workstation"`
	DomainName              []byte
	Workstation             []byte
}

type Challenge struct {
	Header
	TargetNameLen          uint16 `smb:"len:TargetName"`
	TargetNameMaxLen       uint16 `smb:"len:TargetName"`
	TargetNameBufferOffset uint32 `smb:"offset:TargetName"`
	NegotiateFlags         uint32
	ServerChallenge        uint64
	Reserved               uint64
	TargetInfoLen          uint16 `smb:"len:TargetInfo"`
	TargetInfoMaxLen       uint16 `smb:"len:TargetInfo"`
	TargetInfoBufferOffset uint32 `smb:"offset:TargetInfo"`
	Version                uint64
	TargetName             []byte
	TargetInfo             *AvPairSlice
}

// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-nlmp/5e550938-91d4-459f-b67d-75d70009e3f3
//ntlm v2认证结构
type NTLMv2Authentication struct {
	Header
	LmChallengeResponseLen                uint16 `smb:"len:LmChallengeResponse"`
	LmChallengeResponseMaxLen             uint16 `smb:"len:LmChallengeResponse"`
	LmChallengeResponseBufferOffset       uint32 `smb:"offset:LmChallengeResponse"`
	NtChallengeResponseLen                uint16 `smb:"len:NtChallengeResponse"`
	NtChallengeResponseMaxLen             uint16 `smb:"len:NtChallengeResponse"`
	NtChallengResponseBufferOffset        uint32 `smb:"offset:NtChallengeResponse"`
	DomainNameLen                         uint16 `smb:"len:DomainName"`
	DomainNameMaxLen                      uint16 `smb:"len:DomainName"`
	DomainNameBufferOffset                uint32 `smb:"offset:DomainName"`
	UserNameLen                           uint16 `smb:"len:UserName"`
	UserNameMaxLen                        uint16 `smb:"len:UserName"`
	UserNameBufferOffset                  uint32 `smb:"offset:UserName"`
	WorkstationLen                        uint16 `smb:"len:Workstation"`
	WorkstationMaxLen                     uint16 `smb:"len:Workstation"`
	WorkstationBufferOffset               uint32 `smb:"offset:Workstation"`
	EncryptedRandomSessionKeyLen          uint16 `smb:"len:EncryptedRandomSessionKey"`
	EncryptedRandomSessionKeyMaxLen       uint16 `smb:"len:EncryptedRandomSessionKey"`
	EncryptedRandomSessionKeyBufferOffset uint32 `smb:"offset:EncryptedRandomSessionKey"`
	NegotiateFlags                        uint32
	DomainName                            []byte `smb:"unicode"`
	UserName                              []byte `smb:"unicode"`
	Workstation                           []byte `smb:"unicode"`
	EncryptedRandomSessionKey             []byte //16字节，会话加密密钥，可以为空
	LmChallengeResponse                   []byte //24字节，lm协商响应
	NtChallengeResponse                   []byte //24字节，nt协商响应
	MIC                                   []byte `smb:"fixed:16"` //16字节，会话完整性校验
}
