package smb2

import (
	"Orca_Puppet/pkg/psexec/encoder"
	"Orca_Puppet/pkg/psexec/ms"
	"Orca_Puppet/pkg/psexec/smb"
	"encoding/hex"
	"errors"
)

// 此文件用于smb2创建文件请求

// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-fscc/ca28ec38-f155-4768-81d6-4bfeb8586fc9
// FileAttributes属性
const (
	FILE_ATTRIBUTE_READONLY              = 0x00000001
	FILE_ATTRIBUTE_HIDDEN                = 0x00000002
	FILE_ATTRIBUTE_SYSTEM                = 0x00000004
	FILE_ATTRIBUTE_DIRECTORY             = 0x00000010
	FILE_ATTRIBUTE_ARCHIVE               = 0x00000020
	FILE_ATTRIBUTE_NORMAL                = 0x00000080
	FILE_ATTRIBUTE_TEMPORARY             = 0x00000100
	FILE_ATTRIBUTE_SPARSE_FILE           = 0x00000200
	FILE_ATTRIBUTE_REPARSE_POINT         = 0x00000400
	FILE_ATTRIBUTE_COMPRESSED            = 0x00000800
	FILE_ATTRIBUTE_OFFLINE               = 0x00001000
	FILE_ATTRIBUTE_NOT_CONTENT_INDEXED   = 0x00002000
	FILE_ATTRIBUTE_ENCRYPTED             = 0x00004000
	FILE_ATTRIBUTE_INTEGRITY_STREAM      = 0x00008000
	FILE_ATTRIBUTE_NO_SCRUB_DATA         = 0x00020000
	FILE_ATTRIBUTE_RECALL_ON_OPEN        = 0x00040000
	FILE_ATTRIBUTE_PINNED                = 0x00080000
	FILE_ATTRIBUTE_UNPINNED              = 0x00100000
	FILE_ATTRIBUTE_RECALL_ON_DATA_ACCESS = 0x00400000
)

// RequestedOplockLevel属性
const (
	SMB2_OPLOCK_LEVEL_NONE      = 0x00
	SMB2_OPLOCK_LEVEL_II        = 0x01
	SMB2_OPLOCK_LEVEL_EXCLUSIVE = 0x08
	SMB2_OPLOCK_LEVEL_BATCH     = 0x09
	SMB2_OPLOCK_LEVEL_LEASE     = 0xFF
)

// ImpersonationLevel属性
const (
	Anonymous      = 0x00000000
	Identification = 0x00000001
	Impersonation  = 0x00000002
	Delegate       = 0x00000003
)

// AccessMask、CreateDisposition属性
const (
	FILE_SUPERSEDE           = 0x00000000
	FILE_OPEN                = 0x00000001
	FILE_CREATE              = 0x00000002
	FILE_OPEN_IF             = 0x00000003
	FILE_OVERWRITE           = 0x00000004
	FILE_OVERWRITE_IF        = 0x00000005
	FILE_ACTION_ADDED_STREAM = 0x00000006
)

// ShareAccess属性
const (
	FILE_SHARE_READ   = 0x00000001
	FILE_SHARE_WRITE  = 0x00000002
	FILE_SHARE_DELETE = 0x00000004
)

// CreateOptionss属性
const (
	FILE_DIRECTORY_FILE            = 0x00000001
	FILE_WRITE_THROUGH             = 0x00000002
	FILE_SEQUENTIAL_ONLY           = 0x00000004
	FILE_NO_INTERMEDIATE_BUFFERING = 0x00000008
	FILE_SYNCHRONOUS_IO_ALERT      = 0x00000010
	FILE_SYNCHRONOUS_IO_NONALERT   = 0x00000020
	FILE_NON_DIRECTORY_FILE        = 0x00000040
	FILE_COMPLETE_IF_OPLOCKED      = 0x00000100
	FILE_NO_EA_KNOWLEDGE           = 0x00000200
	FILE_RANDOM_ACCESS             = 0x00000800
	FILE_DELETE_ON_CLOSE           = 0x00001000
	FILE_OPEN_BY_FILE_ID           = 0x00002000
	FILE_OPEN_FOR_BACKUP_INTENT    = 0x00004000
	FILE_NO_COMPRESSION            = 0x00008000
	FILE_OPEN_REMOTE_INSTANCE      = 0x00000400
	FILE_OPEN_REQUIRING_OPLOCK     = 0x00010000
	FILE_DISALLOW_EXCLUSIVE        = 0x00020000
	FILE_RESERVE_OPFILTER          = 0x00100000
	FILE_OPEN_REPARSE_POINT        = 0x00200000
	FILE_OPEN_NO_RECALL            = 0x00400000
	FILE_OPEN_FOR_FREE_SPACE_QUERY = 0x00800000
)

// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-smb2/77b36d0f-6016-458a-a7a0-0f4a72ae1534
// DesiredAccess属性
const (
	FILE_READ_DATA         = 0x00000001
	FILE_WRITE_DATA        = 0x00000002
	FILE_APPEND_DATA       = 0x00000004
	FILE_READ_EA           = 0x00000008
	FILE_WRITE_EA          = 0x00000010
	FILE_DELETE_CHILD      = 0x00000040
	FILE_EXECUTE           = 0x00000020
	FILE_READ_ATTRIBUTES   = 0x00000080
	FILE_WRITE_ATTRIBUTES  = 0x00000100
	DELETE                 = 0x00010000
	READ_CONTROL           = 0x00020000
	WRITE_DAC              = 0x00040000
	WRITE_OWNER            = 0x00080000
	SYNCHRONIZE            = 0x00100000
	ACCESS_SYSTEM_SECURITY = 0x01000000
	MAXIMUM_ALLOWED        = 0x02000000
	GENERIC_ALL            = 0x10000000
	GENERIC_EXECUTE        = 0x20000000
	GENERIC_WRITE          = 0x40000000
	GENERIC_READ           = 0x80000000
)

// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-smb2/e8fb45c1-a03d-44ca-b7ae-47385cfd7997
// 创建请求结构
type CreateRequestStruct struct {
	smb.SMB2Header
	StructureSize        uint16
	SecurityFlags        uint8  //1字节，保留字段，不得使用
	OpLock               uint8  //1字节，对应文档RequestedOplockLevel字段
	ImpersonationLevel   uint32 //4字节，模拟等级
	CreateFlags          []byte `smb:"fixed:8"` //8字节，保留字段，不得使用
	Reserved             []byte `smb:"fixed:8"`
	AccessMask           uint32 //4字节，访问权限
	FileAttributes       uint32 //4字节，文件属性
	ShareAccess          uint32 //4字节，共享模式
	CreateDisposition    uint32
	CreateOptions        uint32
	FilenameBufferOffset uint16 `smb:"offset:Filename"`
	FilenameBufferLength uint16 `smb:"len:Filename"`
	CreateContextsOffset uint32
	CreateContextsLength uint32
	Filename             []byte `smb:"unicode"`
}

// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-smb2/d166aa9e-0b53-410e-b35e-3933d8131927
// 创建请求响应结构
type CreateResponseStruct struct {
	smb.SMB2Header
	StructureSize        uint16
	Oplock               uint8 //1字节，对应文档RequestedOplockLevel字段
	ResponseFlags        uint8
	CreateAction         uint32
	CreationTime         []byte `smb:"fixed:8"` //8字节，创建时间
	LastAccessTime       []byte `smb:"fixed:8"` //8字节
	LastWriteTime        []byte `smb:"fixed:8"` //8字节
	LastChangeTime       []byte `smb:"fixed:8"` //8字节
	AllocationSize       []byte `smb:"fixed:8"` //8字节，文件大小
	EndofFile            []byte `smb:"fixed:8"` //8字节
	FileAttributes       uint32
	Reserved2            uint32 `smb:"fixed:4"`
	FileId               []byte `smb:"fixed:16"` //16字节，文件句柄
	CreateContextsOffset uint32
	CreateContextsLength uint32
}

// 创建文件请求
func (c *Client) NewCreateRequest(treeId uint32, filename string, r CreateRequestStruct) CreateRequestStruct {
	smb2Header := NewSMB2Header()
	smb2Header.Command = smb.SMB2_CREATE
	smb2Header.CreditCharge = 1
	smb2Header.MessageId = c.GetMessageId()
	smb2Header.SessionId = c.GetSessionId()
	smb2Header.TreeId = treeId
	r.SMB2Header = smb2Header
	r.StructureSize = 57
	r.SecurityFlags = 0
	r.CreateFlags = make([]byte, 8)
	r.Reserved = make([]byte, 8)
	r.CreateContextsOffset = 0
	r.CreateContextsLength = 0
	r.Filename = encoder.ToUnicode(filename)
	return r
}

// 创建请求响应
func NewCreateResponse() CreateResponseStruct {
	smb2Header := NewSMB2Header()
	return CreateResponseStruct{
		SMB2Header:     smb2Header,
		CreationTime:   make([]byte, 8),
		LastAccessTime: make([]byte, 8),
		LastWriteTime:  make([]byte, 8),
		LastChangeTime: make([]byte, 8),
		AllocationSize: make([]byte, 8),
		EndofFile:      make([]byte, 8),
		FileId:         make([]byte, 16),
	}
}

func (c *Client) CreateRequest(treeId uint32, filename string, r CreateRequestStruct) (fileId []byte, err error) {
	c.Debug("Sending Create file request ["+filename+"]", nil)
	req := c.NewCreateRequest(treeId, filename, r)
	buf, err := c.Send(req)
	if err != nil {
		c.Debug("", err)
		return nil, err
	}
	res := NewCreateResponse()
	c.Debug("Unmarshalling Create file response ["+filename+"]", nil)
	if err := encoder.Unmarshal(buf, &res); err != nil {
		c.Debug("Raw:\n"+hex.Dump(buf), err)
	}
	if res.SMB2Header.Status != ms.STATUS_SUCCESS {
		return nil, errors.New("Failed to create file to [" + filename + "]: " + ms.StatusMap[res.SMB2Header.Status])
	}
	c.Debug("Completed CreateFile ["+filename+"]", nil)
	return res.FileId, nil
}

// 打开管道
func (c *Client) CreatePipeRequest(treeId uint32, pipename string) (fileId []byte, err error) {
	r := CreateRequestStruct{
		OpLock:             SMB2_OPLOCK_LEVEL_NONE,
		ImpersonationLevel: Impersonation,
		AccessMask:         FILE_READ_DATA | FILE_WRITE_DATA | FILE_APPEND_DATA | FILE_READ_EA | FILE_WRITE_EA | FILE_READ_ATTRIBUTES | FILE_WRITE_ATTRIBUTES | READ_CONTROL | SYNCHRONIZE,
		FileAttributes:     FILE_ATTRIBUTE_NORMAL,
		ShareAccess:        FILE_SHARE_READ,
		CreateDisposition:  FILE_OPEN,
		CreateOptions:      FILE_NON_DIRECTORY_FILE,
	}
	fileId, err = c.CreateRequest(treeId, pipename, r)
	if err != nil {
		c.Debug("", err)
		return nil, err
	}
	return fileId, nil
}
