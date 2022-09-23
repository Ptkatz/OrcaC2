package rpce

import (
	"github.com/C-Sto/goWMIExec/pkg/uuid"
)

//2.2.1.1.1 (same as rpc_if_id_t)
type RPC_IF_ID struct {
	UUID                 uuid.UUID
	VersMajor, VersMinor uint16
}

type Version_t struct {
	Major, Minor uint16
}

type P_rt_versions_supported_t struct {
	Protocols   uint8
	P_Protocols []Version_t
}

type SecurityProviders byte

const (
	RPC_C_AUTHN_NONE          = SecurityProviders(0)
	RPC_C_AUTHN_GSS_NEGOTIATE = SecurityProviders(0x09) //SPNEGO
	RPC_C_AUTHN_WINNT         = SecurityProviders(0x0a) //NTLM
	RPC_C_AUTHN_GSS_SCHANNEL  = SecurityProviders(0x0e) //TLS
	RPC_C_AUTHN_GSS_KERBEROS  = SecurityProviders(0x10) //Kerberos
	RPC_C_AUTHN_NETLOGON      = SecurityProviders(0x44)
	RPC_C_AUTHN_DEFAULT       = SecurityProviders(0xff)
)

type AuthLevel uint8

const (
	RPC_C_AUTHN_LEVEL_DEFAULT       AuthLevel = iota //0x00 Same as RPC_C_AUTHN_LEVEL_CONNECT
	RPC_C_AUTHN_LEVEL_NONE                           // 0x01 No authentication.
	RPC_C_AUTHN_LEVEL_CONNECT                        // 0x02 Authenticates the credentials of the client and server.
	RPC_C_AUTHN_LEVEL_CALL                           //0x03 Same as RPC_C_AUTHN_LEVEL_PKT.
	RPC_C_AUTHN_LEVEL_PKT                            //0x04 Same as RPC_C_AUTHN_LEVEL_CONNECT but also prevents replay attacks.
	RPC_C_AUTHN_LEVEL_PKT_INTEGRITY                  //0x05 Same as RPC_C_AUTHN_LEVEL_PKT but also verifies that none of the data transferred between the client and server has been modified.
	RPC_C_AUTHN_LEVEL_PKT_PRIVACY                    // 0x06 Same as RPC_C_AUTHN_LEVEL_PKT_INTEGRITY but also ensures that the data transferred can only be
)

type Twr_t struct {
	TowerLength uint32 //max of 2000
	Tower       []byte //or string?
}

type ErrorStatus uint32
