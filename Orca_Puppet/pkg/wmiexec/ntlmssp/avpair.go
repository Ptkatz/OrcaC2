package ntlmssp

//https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-nlmp/83f5e789-660d-4781-8491-5f8c6641f75e

/*
The AV_PAIR structure defines an attribute/value pair.
Sequences of AV_PAIR structures are used in the CHALLENGE_MESSAGE (section 2.2.1.2) directly.
They are also in the AUTHENTICATE_MESSAGE (section 2.2.1.3) via the NTLMv2_CLIENT_CHALLENGE (section 2.2.2.7) structure.
*/

type AVID uint16

const (
	MsvAvEOL             = 0x0000 //Indicates that this is the last AV_PAIR in the list. AvLen MUST be 0. This type of information MUST be present in the AV pair list.
	MsvAvNbComputerName  = 0x0001 //The server's NetBIOS computer name. The name MUST be in Unicode, and is not null-terminated. This type of information MUST be present in the AV_pair list.
	MsvAvNbDomainName    = 0x0002 //The server's NetBIOS domain name. The name MUST be in Unicode, and is not null-terminated. This type of information MUST be present in the AV_pair list.
	MsvAvDnsComputerName = 0x0003 //The fully qualified domain name (FQDN) of the computer. The name MUST be in Unicode, and is not null-terminated.
	MsvAvDnsDomainName   = 0x0004 //The FQDN of the domain. The name MUST be in Unicode, and is not null-terminated.
	MsvAvDnsTreeName     = 0x0005 //The FQDN of the forest. The name MUST be in Unicode, and is not null-terminated.<13>
	MsvAvFlags           = 0x0006 //A 32-bit value indicating server or client configuration. 0x00000001: Indicates to the client that the account authentication is constrained. 0x00000002: Indicates that the client is providing message integrity in the MIC field (section 2.2.1.3) in the AUTHENTICATE_MESSAGE.<14> 0x00000004: Indicates that the client is providing a target SPN generated from an untrusted source.<15>
	MsvAvTimestamp       = 0x0007 //A FILETIME structure ([MS-DTYP] section 2.3.3) in little-endian byte order that contains the server local time. This structure is always sent in the CHALLENGE_MESSAGE.<16>
	MsvAvSingleHost      = 0x0008 //A Single_Host_Data (section 2.2.2.2) structure. The Value field contains a platform-specific blob, as well as a MachineID created at computer startup to identify the calling machine.<17>
	MsvAvTargetName      = 0x0009 //The SPN of the target server. The name MUST be in Unicode and is not null-terminated.<18>
	MsvAvChannelBindings = 0x000A //A channel bindings hash. The Value field contains an MD5 hash ([RFC4121] section 4.1.1.2) of a gss_channel_bindings_struct ([RFC2744] section 3.11). An all-zero value of the hash is used to indicate absence of channel bindings.<19>
)

type AV_Pair struct {
	AvID  AVID
	AvLen uint16
	Value []byte
}
