package ntlmssp

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"math/rand"
	"strings"
	"time"

	"golang.org/x/crypto/md4"
	"golang.org/x/text/encoding/unicode"
)

type SSP_Negotiate struct {
	Signature         [8]byte
	MessageType       uint32
	NegotiateFlags    NegotiateFlag
	DomainNameFields  SSP_FeildInformation
	WorkstationFields SSP_FeildInformation
	Version           Version
	Payload           NegotiatePayload
}

type NegotiatePayload struct {
	DomainName      []byte
	WorkstationName []byte
}

func NewSSPNegotiate(flags NegotiateFlag) SSP_Negotiate {
	r := SSP_Negotiate{}
	//offset := uint32(64)
	r.Signature = [8]byte{0x4e, 0x54, 0x4c, 0x4d, 0x53, 0x53, 0x50, 0x00}
	r.MessageType = 1
	r.NegotiateFlags = flags
	r.DomainNameFields = NewSSPFeildInformation(0, 0)
	r.WorkstationFields = NewSSPFeildInformation(0, 0)
	r.Version = Version{
		ProductMajor:        0x06,
		ProductMinor:        0x01,
		Build:               0x1db1, // (windows 7, probably a good signature for naughty activity)
		NTLMRevisionCurrent: 0x0f,
	}

	return r
}

func (s SSP_Negotiate) Bytes() []byte {
	buff := bytes.Buffer{}
	binary.Write(&buff, binary.LittleEndian, s.Signature)
	binary.Write(&buff, binary.LittleEndian, s.MessageType)
	binary.Write(&buff, binary.LittleEndian, s.NegotiateFlags)
	binary.Write(&buff, binary.LittleEndian, s.DomainNameFields)
	binary.Write(&buff, binary.LittleEndian, s.WorkstationFields)
	binary.Write(&buff, binary.LittleEndian, s.Version)
	binary.Write(&buff, binary.LittleEndian, s.Payload.DomainName)
	binary.Write(&buff, binary.LittleEndian, s.Payload.WorkstationName)
	return buff.Bytes()
}

type SSP_FeildInformation struct {
	Len          uint16
	MaxLen       uint16
	BufferOffset uint32
}

func NewSSPFeildInformation(len uint16, offset uint32) SSP_FeildInformation {
	return SSP_FeildInformation{Len: len, MaxLen: len, BufferOffset: offset}
}

type SSP_Challenge struct {
	Signature        [8]byte
	MessageType      uint32
	TargetNameFields SSP_FeildInformation
	NegotiateFlags   uint32
	ServerChallenge  [8]byte
	Reserved         [8]byte
	TargetInfoFields SSP_FeildInformation
	Version          [8]byte
	Payload          ChallengePayload
}

type ChallengePayload struct {
	TargetName []byte
	TargetInfo []AV_Pair
}

func (c ChallengePayload) GetTargetInfoBytes() []byte {
	buff := bytes.Buffer{}
	for _, av := range c.TargetInfo {
		binary.Write(&buff, binary.LittleEndian, av.AvID)
		binary.Write(&buff, binary.LittleEndian, av.AvLen)
		buff.Write(av.Value)
	}

	return buff.Bytes()
}

func (c ChallengePayload) GetTimeBytes() []byte {
	for _, av := range c.TargetInfo {
		if av.AvID == MsvAvTimestamp {
			return av.Value
		}
	}
	return nil
}

func ParseSSPChallenge(b []byte) SSP_Challenge {
	cursor := 0

	r := SSP_Challenge{}
	copy(r.Signature[:], b[:8])
	cursor += 8
	binary.Read(bytes.NewReader(b[cursor:]), binary.LittleEndian, &r.MessageType)
	cursor += 4
	binary.Read(bytes.NewReader(b[cursor:]), binary.LittleEndian, &r.TargetNameFields)
	cursor += 8
	binary.Read(bytes.NewReader(b[cursor:]), binary.LittleEndian, &r.NegotiateFlags)
	cursor += 4
	binary.Read(bytes.NewReader(b[cursor:]), binary.LittleEndian, &r.ServerChallenge)
	cursor += 8
	//reserved??
	copy(r.Reserved[:], b[cursor:cursor+8])
	cursor += 8
	binary.Read(bytes.NewReader(b[cursor:]), binary.LittleEndian, &r.TargetInfoFields)
	cursor += 8
	copy(r.Version[:], b[cursor:cursor+8])
	cursor += 8

	//complicated lol
	r.Payload = ChallengePayload{}
	r.Payload.TargetName = make([]byte, r.TargetNameFields.Len)
	copy(r.Payload.TargetName,
		b[r.TargetNameFields.BufferOffset:r.TargetNameFields.BufferOffset+uint32(r.TargetNameFields.Len)])

	tmpPairs := []AV_Pair{}
	avOffset := r.TargetInfoFields.BufferOffset

	//REALLY WHAT THE FUCK MS?
	for {
		if b[avOffset] == MsvAvEOL {
			break
		}
		tmpPair := AV_Pair{}
		binary.Read(bytes.NewReader(b[avOffset:]), binary.LittleEndian, &tmpPair.AvID)
		avOffset += 2
		binary.Read(bytes.NewReader(b[avOffset:]), binary.LittleEndian, &tmpPair.AvLen)
		avOffset += 2
		tmpPair.Value = make([]byte, tmpPair.AvLen)
		copy(tmpPair.Value, b[avOffset:avOffset+uint32(tmpPair.AvLen)])
		tmpPairs = append(tmpPairs, tmpPair)
		avOffset += uint32(tmpPair.AvLen)
	}
	tmpPairs = append(tmpPairs, AV_Pair{})
	r.Payload.TargetInfo = tmpPairs
	return r
}

type SSP_Authenticate struct {
	Signature                       [8]byte              //8
	MessageType                     uint32               //12
	LmChallengeResponseFields       SSP_FeildInformation //20
	NtChallengeResponseFields       SSP_FeildInformation //28
	DomainNameFields                SSP_FeildInformation //36
	UsernameFields                  SSP_FeildInformation //44
	WorkstationFields               SSP_FeildInformation //52
	EncryptedRandomSessionKeyFields SSP_FeildInformation //60
	NegotiateFlags                  uint32               //64
	//Version                         [8]byte              //72
	//MIC     [16]byte //88 //https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-nlmp/a211d894-21bc-4b8b-86ba-b83d0c167b00#Appendix_A_12 HMMMM
	Payload authenticatePayload
}

func NewSSPAuthenticate(response, domainName, username, workstation, sessionkey []byte) SSP_Authenticate {
	r := SSP_Authenticate{
		Signature:   [8]byte{0x4e, 0x54, 0x4c, 0x4d, 0x53, 0x53, 0x50, 0x00},
		MessageType: 3,
	}
	payloadOffset := 64                                                            //only because no MIC and no negotiate flag to do version
	r.LmChallengeResponseFields = NewSSPFeildInformation(0, uint32(payloadOffset)) //not supporting lm I guess
	r.NtChallengeResponseFields = NewSSPFeildInformation(uint16(len(response)), uint32(payloadOffset))
	r.Payload.NtChallengeResponse = response
	payloadOffset += len(response)
	r.DomainNameFields = NewSSPFeildInformation(uint16(len(domainName)), uint32(payloadOffset))
	r.Payload.DomainName = domainName
	payloadOffset += len(domainName)
	r.UsernameFields = NewSSPFeildInformation(uint16(len(username)), uint32(payloadOffset))
	r.Payload.UserName = username
	payloadOffset += len(username)
	r.WorkstationFields = NewSSPFeildInformation(uint16(len(workstation)), uint32(payloadOffset))
	r.Payload.Workstation = workstation
	payloadOffset += len(workstation)
	//r.WorkstationFields = NewSSPFeildInformation(uint16(len(workstation)), uint32(72+payloadOffset))
	r.Payload.EncryptedRandomSessionKey = sessionkey
	r.NegotiateFlags = 0xa2888215 // hard coded for now - flags should be selected sanely in the future

	//0x18, 0x00, 0x18, 0x00}

	return r
}

func (s SSP_Authenticate) Bytes() []byte {
	buff := bytes.Buffer{}
	binary.Write(&buff, binary.LittleEndian, s.Signature)
	binary.Write(&buff, binary.LittleEndian, s.MessageType)
	binary.Write(&buff, binary.LittleEndian, s.LmChallengeResponseFields)
	binary.Write(&buff, binary.LittleEndian, s.NtChallengeResponseFields)
	binary.Write(&buff, binary.LittleEndian, s.DomainNameFields)
	binary.Write(&buff, binary.LittleEndian, s.UsernameFields)
	binary.Write(&buff, binary.LittleEndian, s.WorkstationFields)
	binary.Write(&buff, binary.LittleEndian, s.EncryptedRandomSessionKeyFields)
	binary.Write(&buff, binary.LittleEndian, s.NegotiateFlags)

	buff.Write(s.Payload.LmChallengeResponse)
	buff.Write(s.Payload.NtChallengeResponse)
	buff.Write(s.Payload.DomainName)
	buff.Write(s.Payload.UserName)
	buff.Write(s.Payload.Workstation)
	buff.Write(s.Payload.EncryptedRandomSessionKey)

	return buff.Bytes()
}

type authenticatePayload struct {
	LmChallengeResponse       []byte
	NtChallengeResponse       []byte
	DomainName                []byte
	UserName                  []byte
	Workstation               []byte
	EncryptedRandomSessionKey []byte
}

type Version struct {
	ProductMajor        byte
	ProductMinor        byte
	Build               uint16
	Reserved            [3]byte
	NTLMRevisionCurrent byte
}

//http://davenport.sourceforge.net/ntlm.html#ntlm2Signing

type MessageSignatureExtended struct {
	Version  uint32
	Checksum [8]byte
	SeqNum   uint32
}

func (m MessageSignatureExtended) Bytes() []byte {
	buff := bytes.Buffer{}
	binary.Write(&buff, binary.LittleEndian, m.Version)
	binary.Write(&buff, binary.LittleEndian, m.Checksum)
	binary.Write(&buff, binary.LittleEndian, m.SeqNum)
	return buff.Bytes()
}

func (m *MessageSignatureExtended) SignValue(seq, value, key []byte) {
	hmacer := hmac.New(md5.New, key)
	hmacer.Write(append(seq, value...))
	copy(m.Checksum[:], hmacer.Sum(nil))
}

func NewMessageSignature(value, key []byte, seq uint32) MessageSignatureExtended {
	r := MessageSignatureExtended{}
	r.Version = 1
	r.SeqNum = seq
	sq := []byte{0, 0, 0, 0}
	binary.LittleEndian.PutUint32(sq, seq)
	r.SignValue(sq, value, key)
	return r
}

//NTLMV2Hash returns the NTLMV2 hash provided a password or hash (if both are provided, the hash takes precidence), username and target info. Assumes all strings are UTF8, and have not yet been converted to UTF16
func NTLMV2Hash(password, hash, username, target string) ([]byte, error) {
	if hash == "" {
		h := md4.New()
		unipw, err := toUnicodeS(password)
		if err != nil {
			return nil, err
		}
		h.Write([]byte(unipw))
		hash = hex.EncodeToString(h.Sum(nil))
	}
	hashBytes, err := hex.DecodeString(hash)
	if err != nil {
		return nil, err
	}
	mac := hmac.New(md5.New, hashBytes)
	idkman, err := toUnicodeS(strings.ToUpper(username) + target)
	if err != nil {
		return nil, err
	}
	mac.Write([]byte(idkman))
	return mac.Sum(nil), nil
}

func toUnicodeS(s string) (string, error) {
	s, e := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder().String(s)
	if e != nil {
		return "", e
	}
	return s, nil
}

func NTLMV2Response(hash, servChal, timestamp, targetInfo []byte) []byte {

	v := []byte{1, 1, 0, 0, 0, 0, 0, 0}
	v = append(v, timestamp...)
	chal := make([]byte, 8)
	rand.Seed(time.Now().UnixNano())
	rand.Read(chal)
	v = append(v, chal...)
	v = append(v, 0, 0, 0, 0)
	v = append(v, targetInfo...)
	v = append(v, 0, 0, 0, 0, 0, 0, 0, 0)

	mac := hmac.New(md5.New, hash)
	mac.Write(servChal)
	mac.Write(v)
	hmacVal := mac.Sum(nil)
	return append(hmacVal, v...)
}

func GenerateClientSigningKey(clientNTLMV2Hash, generatedNTLMV2Response []byte) []byte {
	mac := hmac.New(md5.New, clientNTLMV2Hash)
	mac.Write(generatedNTLMV2Response[:mac.Size()])
	base := mac.Sum(nil)

	//signingConst := "session key to client-to-server signing key magic constant.\x00" //what on earth was MS smoking
	signingConst := []byte{0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x20, 0x6b, 0x65, 0x79, 0x20, 0x74, 0x6f, 0x20, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2d, 0x74, 0x6f, 0x2d, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x20, 0x73, 0x69, 0x67, 0x6e, 0x69, 0x6e, 0x67, 0x20, 0x6b, 0x65, 0x79, 0x20, 0x6d, 0x61, 0x67, 0x69, 0x63, 0x20, 0x63, 0x6f, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x00}
	md5er := md5.New()
	md5er.Write(append(base, signingConst...))
	return md5er.Sum(nil)
}
