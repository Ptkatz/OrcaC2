package ntlmssp

//goodness me
//https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-nlmp/99d90ff4-957f-4c8a-80e4-5bfe5a9a9832

type NegotiateFlag uint32

const (
	NTLMSSP_NEGOTIATE_UNICODE NegotiateFlag = 1 << iota
	NTLM_NEGOTIATE_OEM
	NTLMSSP_REQUEST_TARGET
	_
	NTLMSSP_NEGOTIATE_SIGN
	NTLMSSP_NEGOTIATE_SEAL
	NTLMSSP_NEGOTIATE_DATAGRAM
	NTLMSSP_NEGOTIATE_LM_KEY
	_
	NTLMSSP_NEGOTIATE_NTLM
	_
	NTLMSSP_ANONYMOUS_CONNECTIONS
	NTLMSSP_NEGOTIATE_OEM_DOMAIN_SUPPLIED
	NTLMSSP_NEGOTIATE_OEM_WORKSTATION_SUPPLIED
	_
	NTLMSSP_NEGOTIATE_ALWAYS_SIGN
	NTLMSSP_TARGET_TYPE_DOMAIN
	NTLMSSP_TARGET_TYPE_SERVER
	_
	NTLMSSP_NEGOTIATE_EXTENDED_SESSIONSECURITY
	NTLMSSP_NEGOTIATE_IDENTIFY
	_
	NTLMSSP_REQUEST_NON_NT_SESSION_KEY
	NTLMSSP_NEGOTIATE_TARGET_INFO
	_
	NTLMSSP_NEGOTIATE_VERSION
	_
	_
	_
	NTLMSSP_NEGOTIATE_128
	NTLMSSP_NEGOTIATE_KEY_EXCH
	NTLMSSP_NEGOTIATE_56
)

//go:generate stringer -type=NegotiateFlag

func Flags(n NegotiateFlag) []string {
	var result []string
	for i := 0; i <= 32; i++ {
		if n&(1<<i) != 0 {
			result = append(result, (NegotiateFlag)(1<<i).String())
		}
	}
	return result
}
