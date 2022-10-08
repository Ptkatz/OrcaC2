package smb

import (
	"Orca_Puppet/pkg/psexec/gss"
)

// 此文件定义SMB协议头

// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-smb2/fb188936-5050-48d3-b350-dc43059638a4

// SMB协议版本头
const (
	ProtocolSMB  = "\xFFSMB"
	ProtocolSMB2 = "\xFESMB"
)

// SMB签名 开启/关闭
const (
	_ uint16 = iota
	SecurityModeSigningEnabled
	SecurityModeSigningRequired
)

// SMB2 Command代码
const (
	SMB2_NEGOTIATE       = 0x0000
	SMB2_SESSION_SETUP   = 0x0001
	SMB2_LOGOFF          = 0x0002
	SMB2_TREE_CONNECT    = 0x0003
	SMB2_TREE_DISCONNECT = 0x0004
	SMB2_CREATE          = 0x0005
	SMB2_CLOSE           = 0x0006
	SMB2_FLUSH           = 0x0007
	SMB2_READ            = 0x0008
	SMB2_WRITE           = 0x0009
	SMB2_LOCK            = 0x000A
	SMB2_IOCTL           = 0x000B
	SMB2_CANCEL          = 0x000C
	SMB2_ECHO            = 0x000D
	SMB2_QUERY_DIRECTORY = 0x000E
	SMB2_CHANGE_NOTIFY   = 0x000F
	SMB2_QUERY_INFO      = 0x0010
	SMB2_SET_INFO        = 0x0011
	SMB2_OPLOCK_BREAK    = 0x0012
)

// SMB2标准头结构
type SMB2Header struct {
	ProtocolId    []byte `smb:"fixed:4"` //4字节，协议标识符，必须设置为 0x424D53FE
	StructureSize uint16 //2字节，协议结构大小，必修设置为64
	CreditCharge  uint16 //2字节，smb2.0.2发送方必须设置为0
	Status        uint32 //4字节，客户端设置0，服务端返回
	Command       uint16 //2字节，需要包含smb2有效命令
	Credits       uint16 //2字节，标识客户端请求信用数
	Flags         uint32 //4字节，标识如何处理请求
	NextCommand   uint32 //4字节，复合请求使用，不用归零
	MessageId     uint64 //8字节，消息唯一标识符
	Reserved      uint32 //2字节，保留字段，归零
	TreeId        uint32 //4字节，标识树连接，必须设置0
	SessionId     uint64 //8字节，会话唯一标识符，必须设置0
	Signature     []byte `smb:"fixed:16"` //16字节，消息未签名则设0
}

// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-smb2/e14db7ff-763a-4263-8b10-0c3944f52fc5
//SMB 修订号
const (
	SMB2_0_2_Dialect = 0x0202
	SMB2_1_Dialect   = 0x0210
	SMB3_0_Dialect   = 0x0300
	SMB3_0_2_Dialect = 0x0302
	SMB3_1_1_Dialect = 0x0311
)

// SMB2 Negotiate 请求头结构
type SMB2NegotiateRequestStruct struct {
	SMB2Header
	StructureSize   uint16   //2字节，客户端必须设置36
	DialectCount    uint16   `smb:"count:Dialects"` //2字节，必须大于0
	SecurityMode    uint16   //2字节，设置是否启用SMB签名
	Reserved        uint16   //2字节，必须设置0
	Capabilities    uint32   //4字节，如果客户端使用SMB3.x，必须使用SMB2_GLOBAL_CAP_*构造，否则设置为0
	ClientGuid      []byte   `smb:"fixed:16"` //16字节，客户端自身生成
	ClientStartTime uint64   //8字节，保留字段，归零
	Dialects        []uint16 //16位整数数组
}

// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-smb2/63abf97c-0d09-47e2-88d6-6bfa552949a5
// SMB2 Negotiate 响应头结构
type SMB2NegotiateResponseStruct struct {
	SMB2Header
	StructureSize        uint16            //2字节，客户端必须设置36
	SecurityMode         uint16            //2字节，设置是否启用SMB签名
	DialectRevision      uint16            //2字节，SMB协议号
	Reserved             uint16            //2字节，保留字段，归零
	ServerGuid           []byte            `smb:"fixed:16"` //16字节，服务器标识符
	Capabilities         uint32            //4字节，服务器协议作用
	MaxTransactSize      uint32            //4字节，客户端set_info请求缓冲区大小
	MaxReadSize          uint32            //4字节，服务器接受smb read请求最大长度
	MaxWriteSize         uint32            //4字节，服务器接受smb write请求最大长度
	SystemTime           uint64            //8字节，处理协商请求服务器系统时间
	ServerStartTime      uint64            //8字节，服务器启动时间
	SecurityBufferOffset uint16            `smb:"offset:SecurityBlob"` //2字节，smb2表头开始到安全缓存区的偏移量
	SecurityBufferLength uint16            `smb:"len:SecurityBlob"`    //2字节，安全缓冲区长度
	Reserved2            uint32            //4字节，协商上下文偏移量
	SecurityBlob         *gss.NegTokenInit //服务器返回二进制安全对象，遵循RFC2743标准
}

// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-smb2/5a3c2c28-d6b0-48ed-b917-a86b2ca4575f
// 质询请求结构体
type SMB2SessionSetupRequestStruct struct {
	SMB2Header
	StructureSize        uint16
	Flags                byte
	SecurityMode         byte
	Capabilities         uint32
	Channel              uint32 //4字节，保留字段，归零
	SecurityBufferOffset uint16 `smb:"offset:SecurityBlob"`
	SecurityBufferLength uint16 `smb:"len:SecurityBlob"`
	PreviousSessionID    uint64 //8字节，会话标识符。服务端用来标识客户端会话
	SecurityBlob         *gss.NegTokenInit
}

// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-smb2/0324190f-a31b-4666-9fa9-5c624273a694
// 质询响应结构体
type SMB2SessionSetupResponseStruct struct {
	SMB2Header
	StructureSize        uint16
	Flags                uint16
	SecurityBufferOffset uint16 `smb:"offset:SecurityBlob"`
	SecurityBufferLength uint16 `smb:"len:SecurityBlob"`
	SecurityBlob         *gss.NegTokenResp
}

// 质询请求认证结构体、需要带上响应
type SMB2SessionSetup2RequestStruct struct {
	SMB2Header
	StructureSize        uint16
	Flags                byte
	SecurityMode         byte
	Capabilities         uint32
	Channel              uint32 //4字节，保留字段，归零
	SecurityBufferOffset uint16 `smb:"offset:SecurityBlob"`
	SecurityBufferLength uint16 `smb:"len:SecurityBlob"`
	PreviousSessionID    uint64 //8字节，会话标识符。服务端用来标识客户端会话
	SecurityBlob         *gss.NegTokenResp
}
