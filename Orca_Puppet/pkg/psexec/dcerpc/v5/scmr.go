package v5

import (
	"Orca_Puppet/pkg/psexec/encoder"
	"Orca_Puppet/pkg/psexec/smb/smb2"
	"Orca_Puppet/pkg/psexec/util"
)

// 此文件提供访问windows服务管理封装

// DCE/RPC 扩展头
// 调用win ms service control api
type PDUExtHeaderStruct struct {
	Version            uint8
	VersionMinor       uint8
	PacketType         uint8
	PacketFlags        uint8
	DataRepresentation uint32 //4字节，小端排序，0x10
	FragLength         uint16 //2字节，整个结构的长度
	AuthLength         uint16
	CallId             uint32
	AllocHint          uint32 `smb:"len:Buffer"` //Buffer的长度
	ContextId          uint16
	OpNum              uint16
	Buffer             interface{}
}

// ms service control
// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-scmr/dc84adb3-d51d-48eb-820d-ba1c6ca5faf2
type OpenSCManagerWStruct struct {
	MachineName machineName
	Database    database
	AccessMask  uint32
}

type machineName struct {
	ReferentId  uint32 `smb:"offset:MachineName"`
	MaxCount    uint32
	Offset      uint32
	ActualCount uint32 //机器名的长度
	MachineName []byte //任意长度,unicode编码
	Reserved    uint16
}

type database struct {
	ReferentId  uint32 `smb:"offset:Database"`
	MaxCount    uint32
	Offset      uint32
	ActualCount uint32 //机器名的长度
	Database    []byte //任意长度,unicode编码
	Reserved    uint16
}

// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-scmr/0d7a7011-9f41-470d-ad52-8535b47ac282
// 安全描述符
const (
	SERVICE_ALL_ACCESS        = 0x000F01FF
	SC_MANAGER_CREATE_SERVICE = 0x00000002
	SC_MANAGER_CONNECT        = 0x00000001
)

// OpenSCManagerW响应结构
type OpenSCManagerWResponse struct {
	smb2.ReadResponseStruct
	Version            uint8
	VersionMinor       uint8
	PacketType         uint8
	PacketFlags        uint8
	DataRepresentation uint32
	FragLength         uint16
	AuthLength         uint16
	CallId             uint32
	AllocHint          uint32
	ContextId          uint16
	CancelCount        uint8
	Reserved           uint8
	ContextHandle      []byte `smb:"fixed:20"`
	ReturnCode         uint32
}

// opnum
// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-scmr/0d7a7011-9f41-470d-ad52-8535b47ac282
const (
	RCloseServiceHandle         = 0
	RControlService             = 1
	RDeleteService              = 2
	RLockServiceDatabase        = 3
	RQueryServiceObjectSecurity = 4
	RSetServiceObjectSecurity   = 5
	RQueryServiceStatus         = 6
	RSetServiceStatus           = 7
	RUnlockServiceDatabase      = 8
	RNotifyBootConfigStatus     = 9
	RChangeServiceConfigW       = 11
	RCreateServiceW             = 12
	REnumDependentServicesW     = 13
	REnumServicesStatusW        = 14
	ROpenSCManagerW             = 15
	ROpenServiceW               = 16
	RQueryServiceConfigW        = 17
	RQueryServiceLockStatusW    = 18
	RStartServiceW              = 19
	RGetServiceDisplayNameW     = 20
	RGetServiceKeyNameW         = 21
	RChangeServiceConfigA       = 23
	RCreateServiceA             = 24
	REnumDependentServicesA     = 25
	REnumServicesStatusA        = 26
	ROpenSCManagerA             = 27
	ROpenServiceA               = 28
	RQueryServiceConfigA        = 29
	RQueryServiceLockStatusA    = 30
	RStartServiceA              = 31
	RGetServiceDisplayNameA     = 32
	RGetServiceKeyNameA         = 33
	REnumServiceGroupW          = 35
	RChangeServiceConfig2A      = 36
	RChangeServiceConfig2W      = 37
	RQueryServiceConfig2A       = 38
	RQueryServiceConfig2W       = 39
	RQueryServiceStatusEx       = 40
	REnumServicesStatusExA      = 41
	REnumServicesStatusExW      = 42
	RCreateServiceWOW64A        = 44
	RCreateServiceWOW64W        = 45
	RNotifyServiceStatusChange  = 47
	RGetNotifyResults           = 48
	RCloseNotifyHandle          = 49
	RControlServiceExA          = 50
	RControlServiceExW          = 51
	RQueryServiceConfigEx       = 56
	RCreateWowService           = 60
	ROpenSCManager2             = 64
)

// OpenSCManagerW请求
//DWORD ROpenSCManagerW(
//	[in, string, unique, range(0, SC_MAX_COMPUTER_NAME_LENGTH)] SVCCTL_HANDLEW lpMachineName,
//	[in, string, unique, range(0, SC_MAX_NAME_LENGTH)] wchar_t* lpDatabaseName,
//	[in] DWORD dwDesiredAccess,
//	[out] LPSC_RPC_HANDLE lpScHandle
//	);
//lpMachineName：一种 SVCCTL_HANDLEW（第 2.2.3 节）数据类型，它定义指向以空字符结尾的 UNICODE 字符串的指针，该字符串指定服务器的机器名称。
//lpDatabaseName：指向以空结尾的 UNICODE 字符串的指针，该字符串指定要打开的 SCM 数据库的名称。该参数必须设置为 NULL、“ServicesActive”或“ServicesFailed”。
//dwDesiredAccess：一个值，指定对数据库的访问。这必须是第 3.1.4 节中指定的值之一。
//客户端还必须具有 SC_MANAGER_CONNECT 访问权限。
//lpScHandle：一种 LPSC_RPC_HANDLE 数据类型，用于定义新打开的 SCM 数据库的句柄。
func (c *Client) NewOpenSCManagerWRequest(treeId uint32, fileId []byte) PDUHeader {
	pduHeader := NewPDUHeader()
	pduHeader.SMB2Header.MessageId = c.GetMessageId()
	pduHeader.SMB2Header.SessionId = c.GetSessionId()
	pduHeader.SMB2Header.TreeId = treeId
	pduHeader.FileId = fileId
	machinename := string(util.Random(6)) + "\x00"
	databaseName := "ServicesActive\x00"
	buffer := OpenSCManagerWStruct{
		MachineName: machineName{
			MaxCount:    uint32(len(machinename)),
			ActualCount: uint32(len(machinename)),
			MachineName: encoder.ToUnicode(machinename),
		},
		Database: database{
			MaxCount:    uint32(len(databaseName)),
			ActualCount: uint32(len(databaseName)),
			Database:    encoder.ToUnicode(databaseName),
		},
		AccessMask: SC_MANAGER_CREATE_SERVICE | SC_MANAGER_CONNECT,
	}
	fragLength := 24 + util.SizeOfStruct(buffer) // 头固定大小24
	pduHeader.Buffer = PDUExtHeaderStruct{
		Version:            5,
		VersionMinor:       0,
		PacketType:         PDURequest,
		PacketFlags:        PDUFault,
		DataRepresentation: 16,
		FragLength:         uint16(fragLength),
		AuthLength:         0,
		CallId:             1,
		ContextId:          0,
		OpNum:              ROpenSCManagerW,
		Buffer:             buffer,
	}
	return pduHeader
}

func NewOpenSCManagerWResponse() OpenSCManagerWResponse {
	return OpenSCManagerWResponse{
		ContextHandle: make([]byte, 20),
	}
}

// 打开服务
// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-scmr/6d0a4225-451b-4132-894d-7cef7aecfd2d
type ROpenServiceWRequestStruct struct {
	ContextHandle []byte `smb:"fixed:20"` //OpenSCManagerW 句柄
	ServiceName   serviceName
	//Reserved      uint8
	AccessMask uint32
}

type ROpenServiceWResponseStruct struct {
	smb2.ReadResponseStruct
	Version            uint8
	VersionMinor       uint8
	PacketType         uint8
	PacketFlags        uint8
	DataRepresentation uint32
	FragLength         uint16
	AuthLength         uint16
	CallId             uint32
	AllocHint          uint32
	ContextId          uint16
	CancelCount        uint8
	Reserved           uint8
	ContextHandle      []byte `smb:"fixed:20"`
	ReturnCode         uint32
}

func (c *Client) NewROpenServiceWRequest(treeId uint32, fileId, contextHandle []byte, servicename string) PDUHeader {
	pduHeader := NewPDUHeader()
	pduHeader.SMB2Header.MessageId = c.GetMessageId()
	pduHeader.SMB2Header.SessionId = c.GetSessionId()
	pduHeader.SMB2Header.TreeId = treeId
	pduHeader.FileId = fileId
	serName := servicename + "\x00"
	buffer := ROpenServiceWRequestStruct{
		ContextHandle: contextHandle,
		ServiceName: serviceName{
			MaxCount:    uint32(len(serName)),
			ActualCount: uint32(len(serName)),
			ServiceName: encoder.ToUnicode(serName),
		},
		AccessMask: SERVICE_ALL_ACCESS,
	}
	fragLength := 24 + util.SizeOfStruct(buffer) // 头固定大小24
	pduHeader.Buffer = PDUExtHeaderStruct{
		Version:            5,
		VersionMinor:       0,
		PacketType:         PDURequest,
		PacketFlags:        PDUFault,
		DataRepresentation: 16,
		FragLength:         uint16(fragLength),
		AuthLength:         0,
		CallId:             2,
		ContextId:          0,
		OpNum:              ROpenServiceW,
		Buffer:             buffer,
	}
	return pduHeader
}

func NewROpenServiceWResponse() ROpenServiceWResponseStruct {
	return ROpenServiceWResponseStruct{
		ContextHandle: make([]byte, 20),
	}
}

// 创建服务
// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-scmr/6a8ca926-9477-4dd4-b766-692fab07227e
type RCreateServiceWRequestStruct struct {
	ContextHandle       []byte `smb:"fixed:20"` //OpenSCManagerW 句柄
	ServiceName         serviceName
	DisplayName         displayName
	AccessMask          uint32
	ServiceType         uint32
	ServiceStartType    uint32
	ServiceErrorControl uint32
	BinaryPathName      binaryPathName
	NULLPointer         uint32
	TagId               uint32
	NULLPointer2        uint32
	DependSize          uint32
	NULLPointer3        uint32
	NULLPointer4        uint32
	PasswordSize        uint32
}

type serviceName struct {
	MaxCount    uint32
	Offset      uint32
	ActualCount uint32
	ServiceName []byte
	Reserved    uint16
}

type displayName struct {
	ReferentId  uint32 `smb:"offset:DisplayName"`
	MaxCount    uint32
	Offset      uint32
	ActualCount uint32
	DisplayName []byte
	Reserved    uint16
}

type binaryPathName struct {
	MaxCount       uint32
	Offset         uint32
	ActualCount    uint32
	BinaryPathName []byte
}

// RCreateServiceW响应结构
type RCreateServiceWResponseStruct struct {
	smb2.ReadResponseStruct
	Version            uint8
	VersionMinor       uint8
	PacketType         uint8
	PacketFlags        uint8
	DataRepresentation uint32
	FragLength         uint16
	AuthLength         uint16
	CallId             uint32
	AllocHint          uint32
	ContextId          uint16
	CancelCount        uint8
	Reserved           uint8
	TagId              uint32
	ContextHandle      []byte `smb:"fixed:20"`
	ReturnCode         uint32
}

// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-scmr/6a8ca926-9477-4dd4-b766-692fab07227e
// dwServiceType 类型
const (
	SERVICE_KERNEL_DRIVER       = 0x00000001
	SERVICE_FILE_SYSTEM_DRIVER  = 0x00000002
	SERVICE_WIN32_OWN_PROCESS   = 0x00000010
	SERVICE_WIN32_SHARE_PROCESS = 0x00000020
	SERVICE_INTERACTIVE_PROCESS = 0x00000100
)

// dwStartType类型
const (
	SERVICE_BOOT_START   = 0x00000000
	SERVICE_SYSTEM_START = 0x00000001
	SERVICE_AUTO_START   = 0x00000002
	SERVICE_DEMAND_START = 0x00000003
	SERVICE_DISABLED     = 0x00000004
)

// dwErrorControl类型
const (
	SERVICE_ERROR_IGNORE   = 0x00000000
	SERVICE_ERROR_NORMAL   = 0x00000001
	SERVICE_ERROR_SEVERE   = 0x00000002
	SERVICE_ERROR_CRITICAL = 0x00000003
)

func (c *Client) NewRCreateServiceWRequest(treeId uint32, fileId, contextHandle []byte, servicename, uploadPathFile string) PDUHeader {
	pduHeader := NewPDUHeader()
	pduHeader.SMB2Header.MessageId = c.GetMessageId()
	pduHeader.SMB2Header.SessionId = c.GetSessionId()
	pduHeader.SMB2Header.TreeId = treeId
	pduHeader.FileId = fileId
	serName := servicename + "\x00"
	uploadpathFile := uploadPathFile + "\x00"
	buffer := RCreateServiceWRequestStruct{
		ContextHandle: contextHandle,
		ServiceName: serviceName{
			MaxCount:    uint32(len(serName)),
			ActualCount: uint32(len(serName)),
			ServiceName: encoder.ToUnicode(serName),
		},
		DisplayName: displayName{
			MaxCount:    uint32(len(serName)),
			ActualCount: uint32(len(serName)),
			DisplayName: encoder.ToUnicode(serName),
		},
		AccessMask:          SERVICE_ALL_ACCESS,
		ServiceType:         SERVICE_WIN32_OWN_PROCESS,
		ServiceStartType:    SERVICE_DEMAND_START,
		ServiceErrorControl: SERVICE_ERROR_IGNORE,
		BinaryPathName: binaryPathName{
			MaxCount:       uint32(len(uploadpathFile)),
			ActualCount:    uint32(len(uploadpathFile)),
			BinaryPathName: encoder.ToUnicode(uploadpathFile),
		},
	}
	fragLength := 24 + util.SizeOfStruct(buffer) // 头固定大小24
	pduHeader.Buffer = PDUExtHeaderStruct{
		Version:            5,
		VersionMinor:       0,
		PacketType:         PDURequest,
		PacketFlags:        PDUFault,
		DataRepresentation: 16,
		FragLength:         uint16(fragLength),
		AuthLength:         0,
		CallId:             3,
		ContextId:          0,
		OpNum:              RCreateServiceW,
		Buffer:             buffer,
	}
	return pduHeader
}

func NewRCreateServiceWResponse() RCreateServiceWResponseStruct {
	return RCreateServiceWResponseStruct{
		ContextHandle: make([]byte, 20),
	}
}

// 启动服务
// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-scmr/d9be95a2-cf01-4bdc-b30f-6fe4b37ada16
type RStartServiceWRequestStruct struct {
	ContextHandle []byte `smb:"fixed:20"` //20字节，创建服务返回的句柄
	Argc          uint32 //argv字符串数量
	Argv          []byte `smb:"fixed:4"` //4字节，unicode对象指针
}

type RStartServiceWResponseStruct struct {
	smb2.ReadResponseStruct
	Version            uint8
	VersionMinor       uint8
	PacketType         uint8
	PacketFlags        uint8
	DataRepresentation uint32
	FragLength         uint16
	AuthLength         uint16
	CallId             uint32
	AllocHint          uint32
	ContextId          uint16
	CancelCount        uint8
	Reserved           uint8
	StubData           uint32
}

// 启动服务封装
func (c *Client) NewRStartServiceWRequest(treeId uint32, fileId, contextHandle []byte) PDUHeader {
	pduHeader := NewPDUHeader()
	pduHeader.SMB2Header.MessageId = c.GetMessageId()
	pduHeader.SMB2Header.SessionId = c.GetSessionId()
	pduHeader.SMB2Header.Credits = 127
	pduHeader.SMB2Header.TreeId = treeId
	pduHeader.FileId = fileId
	argv := encoder.ToUnicode(string(make([]byte, 4)))
	buffer := RStartServiceWRequestStruct{
		ContextHandle: contextHandle,
		Argc:          0,
		Argv:          argv[:4],
	}
	fragLength := 24 + util.SizeOfStruct(buffer)
	pduHeader.Buffer = PDUExtHeaderStruct{
		Version:            5,
		VersionMinor:       0,
		PacketType:         PDURequest,
		PacketFlags:        PDUFault,
		DataRepresentation: 16,
		FragLength:         uint16(fragLength),
		AuthLength:         0,
		CallId:             4,
		ContextId:          0,
		OpNum:              RStartServiceW,
		Buffer:             buffer,
	}
	return pduHeader
}

// 启动服务响应封装
func NewRStartServiceWResponse() RStartServiceWResponseStruct {
	return RStartServiceWResponseStruct{}
}

// 关闭服务句柄
type RCloseServiceHandleRequestStruct struct {
	ContextHandle []byte `smb:"fixed:20"`
}

type RCloseServiceHandleResponseStruct struct {
	smb2.ReadResponseStruct
	Version            uint8
	VersionMinor       uint8
	PacketType         uint8
	PacketFlags        uint8
	DataRepresentation uint32
	FragLength         uint16
	AuthLength         uint16
	CallId             uint32
	AllocHint          uint32
	ContextId          uint16
	CancelCount        uint8
	Reserved           uint8
	ContextHandle      []byte `smb:"fixed:20"`
	ReturnCode         uint32
}

func (c *Client) NewRCloseServiceHandleRequest(treeId uint32, fileId, contextHandle []byte) PDUHeader {
	pduHeader := NewPDUHeader()
	pduHeader.SMB2Header.MessageId = c.GetMessageId()
	pduHeader.SMB2Header.SessionId = c.GetSessionId()
	pduHeader.SMB2Header.TreeId = treeId
	pduHeader.FileId = fileId
	buffer := RCloseServiceHandleRequestStruct{ContextHandle: contextHandle}
	fragLength := 24 + util.SizeOfStruct(buffer)
	pduHeader.Buffer = PDUExtHeaderStruct{
		Version:            5,
		VersionMinor:       0,
		PacketType:         PDURequest,
		PacketFlags:        PDUFault,
		DataRepresentation: 16,
		FragLength:         uint16(fragLength),
		AuthLength:         0,
		CallId:             6,
		ContextId:          0,
		OpNum:              RCloseServiceHandle,
		Buffer:             buffer,
	}
	return pduHeader
}

func NewRCloseServiceHandleResponse() RCloseServiceHandleResponseStruct {
	return RCloseServiceHandleResponseStruct{
		ContextHandle: make([]byte, 20),
	}
}
