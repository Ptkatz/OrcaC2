package ntlm

// 此文件用于ntlm认证相关

import (
	"Orca_Puppet/pkg/psexec/encoder"
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"encoding/binary"
	"hash"
	"time"
)

// 协商版本
func NewNegotiate(domainName, workstation string) Negotiate {
	return Negotiate{
		Header: Header{
			Signature:   []byte(NTLMSecSignature),
			MessageType: NTLMNegotiate,
		},
		NegotiateFlags: FlgNeg56 |
			FlgNeg128 |
			FlgNegTargetInfo |
			FlgNegExtendedSessionSecurity |
			FlgNegOEMDomainSupplied |
			FlgNegNTLM |
			FlgNegRequestTarget |
			FlgNegUnicode,
		DomainNameLen:           0,
		DomainNameMaxLen:        0,
		DomainNameBufferOffset:  0,
		WorkstationLen:          0,
		WorkstationMaxLen:       0,
		WorkstationBufferOffset: 0,
		DomainName:              []byte(domainName),
		Workstation:             []byte(workstation),
	}
}

// 协商密钥
func NewChallenge() Challenge {
	return Challenge{
		Header: Header{
			Signature:   []byte(NTLMSecSignature),
			MessageType: NTLMChallenge,
		},
		TargetNameLen:          0,
		TargetNameMaxLen:       0,
		TargetNameBufferOffset: 0,
		NegotiateFlags: FlgNeg56 |
			FlgNeg128 |
			FlgNegVersion |
			FlgNegTargetInfo |
			FlgNegExtendedSessionSecurity |
			FlgNegTargetTypeServer |
			FlgNegNTLM |
			FlgNegRequestTarget |
			FlgNegUnicode,
		ServerChallenge:        0,
		Reserved:               0,
		TargetInfoLen:          0,
		TargetInfoMaxLen:       0,
		TargetInfoBufferOffset: 0,
		Version:                0,
		TargetName:             []byte{},
		TargetInfo:             new(AvPairSlice),
	}
}

func NewAuthenticatePass(domain, user, workstation, password string, c Challenge) NTLMv2Authentication {
	// 明文认证
	//nthash := NTOWFv2(password, user, domain)
	//lmhash := LMOWFv2(password, user, domain)
	h := hmac.New(md5.New, NTOWFv2(password, user, domain))
	//return newAuthenticate(h, domain, user, workstation, nthash, lmhash, c)
	return newAuthenticate(h, domain, user, workstation, c)
}

func NewAuthenticateHash(domain, user, workstation, hash string, c Challenge) NTLMv2Authentication {
	// hash认证
	//buf := make([]byte, len(hash)/2)
	//hex.Decode(buf, []byte(hash))
	h := hmac.New(md5.New, NTOWFv2Hash(hash, user, domain))
	return newAuthenticate(h, domain, user, workstation, c)
}

func newAuthenticate(h hash.Hash, domain, user, workstation string, c Challenge) NTLMv2Authentication {
	// Assumes domain, user, and workstation are not unicode
	var timestamp []byte
	for k, av := range *c.TargetInfo {
		if av.AvID == MsvAvTimestamp {
			timestamp = (*c.TargetInfo)[k].Value
		}
	}
	if timestamp == nil {
		// Credit to https://github.com/Azure/go-ntlmssp/blob/master/unicode.go for logic
		ft := uint64(time.Now().UnixNano()) / 100
		ft += 116444736000000000 // add time between unix & windows offset
		timestamp = make([]byte, 8)
		binary.LittleEndian.PutUint64(timestamp, ft)
	}

	clientChallenge := make([]byte, 8)
	rand.Reader.Read(clientChallenge)
	serverChallenge := make([]byte, 8)
	w := bytes.NewBuffer(make([]byte, 0))
	binary.Write(w, binary.LittleEndian, c.ServerChallenge)
	serverChallenge = w.Bytes()
	w = bytes.NewBuffer(make([]byte, 0))
	for _, av := range *c.TargetInfo {
		binary.Write(w, binary.LittleEndian, av.AvID)
		binary.Write(w, binary.LittleEndian, av.AvLen)
		binary.Write(w, binary.LittleEndian, av.Value)
	}
	ntChallengeResponse, lmChallengeResponse, sessionBaseKey := ComputeNTLMv2Response(h, clientChallenge, serverChallenge, timestamp, w.Bytes())

	return NTLMv2Authentication{
		Header: Header{
			Signature:   []byte(NTLMSecSignature),
			MessageType: NTLMAuthenticate,
		},
		DomainName:  encoder.ToUnicode(domain),
		UserName:    encoder.ToUnicode(user),
		Workstation: encoder.ToUnicode(workstation),
		NegotiateFlags: FlgNeg56 |
			FlgNeg128 |
			FlgNegTargetInfo |
			FlgNegExtendedSessionSecurity |
			FlgNegNTLM |
			FlgNegRequestTarget |
			FlgNegUnicode,
		NtChallengeResponse:       ntChallengeResponse,
		LmChallengeResponse:       lmChallengeResponse,
		EncryptedRandomSessionKey: sessionBaseKey,
	}
}
