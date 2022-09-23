package smb2

import (
	"Orca_Puppet/pkg/psexec/encoder"
	"Orca_Puppet/pkg/psexec/smb"
	"encoding/hex"
)

// loctl/fsctl封装

// Function属性
const (
	FSCTL_DFS_GET_REFERRALS            = 0x00060194
	FSCTL_PIPE_PEEK                    = 0x0011400C
	FSCTL_PIPE_WAIT                    = 0x00110018
	FSCTL_PIPE_TRANSCEIVE              = 0x0011C017
	FSCTL_SRV_COPYCHUNK                = 0x001440F2
	FSCTL_SRV_ENUMERATE_SNAPSHOTS      = 0x00144064
	FSCTL_SRV_REQUEST_RESUME_KEY       = 0x00140078
	FSCTL_SRV_READ_HASH                = 0x001441bb
	FSCTL_SRV_COPYCHUNK_WRITE          = 0x001480F2
	FSCTL_LMR_REQUEST_RESILIENCY       = 0x001401D4
	FSCTL_QUERY_NETWORK_INTERFACE_INFO = 0x001401FC
	FSCTL_SET_REPARSE_POINT            = 0x000900A4
	FSCTL_DFS_GET_REFERRALS_EX         = 0x000601B0
	FSCTL_FILE_LEVEL_TRIM              = 0x00098208
	FSCTL_VALIDATE_NEGOTIATE_INFO      = 0x00140204
)

// Flags属性
const (
	SMB2_0_IOCTL_IS_IOCTL = 0x00000000
	SMB2_0_IOCTL_IS_FSCTL = 0x00000001
)

// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-smb2/5c03c9d6-15de-48a2-9835-8fb37f8a79d8
type IOCTLRequestStruct struct {
	smb.SMB2Header
	StructureSize     uint16
	Reserved          uint16
	Function          uint32
	GUIDHandle        []byte `smb:"fixed:16"`
	InputOffset       uint32 `smb:"offset:Buffer"`
	InputCount        uint32 `smb:"len:Buffer"`
	MaxInputResponse  uint32
	OutputOffset      uint32
	OutputCount       uint32
	MaxOutputResponse uint32
	Flags             uint32
	Reserved2         uint32
	Buffer            interface{}
}

func (c *Client) NewIOCTLRequest(treeId uint32) IOCTLRequestStruct {
	smb2Header := NewSMB2Header()
	smb2Header.Command = smb.SMB2_IOCTL
	smb2Header.CreditCharge = 1
	smb2Header.MessageId = c.GetMessageId()
	smb2Header.SessionId = c.GetSessionId()
	smb2Header.TreeId = treeId
	smb2Header.Credits = 127
	return IOCTLRequestStruct{
		SMB2Header:    smb2Header,
		StructureSize: 57,
		GUIDHandle:    make([]byte, 16),
	}
}

// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-smb2/f70eccb6-e1be-4db8-9c47-9ac86ef18dbb
type IOCTLResponseStruct struct {
	smb.SMB2Header
	StructureSize uint16
	Reserved      uint16
	Function      uint32
	GUIDHandle    []byte `smb:"fixed:16"`
	BlobOffset    uint32
	BlobLength    uint32
	BlobOffset2   uint32
	BlobLength2   uint32
	Flags         uint32
	Reserved2     uint32
}

func NewIOCTLResponse() IOCTLResponseStruct {
	return IOCTLResponseStruct{
		GUIDHandle: make([]byte, 20),
	}
}

// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-fscc/f030a3b9-539c-4c7b-a893-86b795b9b711
// 请求服务器等待连接
type FSCTLPIPEWAITRequestStruct struct {
	Timeout          uint64 // 毫秒单位
	NameLength       uint32
	TimeoutSpecified uint8 // 一个布尔值，指定是否忽略Timeout参数
	Padding          uint8
	Name             []byte // 命名管道名称的Unicode字符串，名称不得包含“\pipe\”
}

func (c *Client) NewFSCTLPIPEWAITRequest(pipename string) FSCTLPIPEWAITRequestStruct {
	pipeName := encoder.ToUnicode(pipename)
	return FSCTLPIPEWAITRequestStruct{
		NameLength:       uint32(len(pipeName)),
		TimeoutSpecified: 1,
		Padding:          0,
		Name:             pipeName,
	}
}

// 连接并绑定命名管道，并拿到管道句柄
func (c *Client) ConnectAndWriteStdInPipes(pipename string) (treeid uint32, pipehandle []byte, err error) {
	timeout := uint64(500000)
	treeId, err := c.TreeConnect("IPC$")
	if err != nil {
		c.Debug("", err)
		return 0, nil, err
	}
	IOCTLRequest := c.NewIOCTLRequest(treeId)
	// 使用FSCTL_PIPE_WAIT，FileId必须为0xFFFFFFFFFFFFFFFF
	IOCTLRequest.Function = FSCTL_PIPE_WAIT
	IOCTLRequest.Flags = SMB2_0_IOCTL_IS_FSCTL
	IOCTLRequest.GUIDHandle = []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
	// 连接管道
	FSCTLPIPEWAITRequestIn := c.NewFSCTLPIPEWAITRequest(pipename)
	FSCTLPIPEWAITRequestIn.Timeout = timeout
	IOCTLRequest.Buffer = FSCTLPIPEWAITRequestIn
	c.Debug("Sending Ioctl stdin pipe request ["+pipename+"]", nil)
	buf, err := c.Send(IOCTLRequest)
	if err != nil {
		c.Debug("", err)
		return 0, nil, err
	}
	res := NewIOCTLResponse()
	c.Debug("Unmarshalling Ioctl stdin pipe response ["+pipename+"]", nil)
	if err = encoder.Unmarshal(buf, &res); err != nil {
		c.Debug("Raw:\n"+hex.Dump(buf), err)
		return 0, nil, err
	}
	// 创建管道请求
	pipeHander, err := c.CreatePipeRequest(treeId, pipename)
	if err != nil {
		return 0, nil, err
	}
	// 将数据写入管道
	err = c.WritePipeRequest(treeId, []byte("cmd"), pipeHander)
	if err != nil {
		return 0, nil, err
	}
	return treeId, pipeHander, nil
}

// 拿到stdin、out、err句柄
func (c *Client) ConnectAndBindNamedPipes(pipename string) (stdinpipe, stdoutpipe, stderrpipe []byte, err error) {
	var stdIn, stdOut, stdErr []byte
	timeout := uint64(500000)
	treeId, err := c.TreeConnect("IPC$")
	if err != nil {
		c.Debug("", err)
		//return nil, err
	}
	IOCTLRequest := c.NewIOCTLRequest(treeId)
	// 使用FSCTL_PIPE_WAIT，FileId必须为0xFFFFFFFFFFFFFFFF
	IOCTLRequest.Function = FSCTL_PIPE_WAIT
	IOCTLRequest.Flags = SMB2_0_IOCTL_IS_FSCTL
	IOCTLRequest.GUIDHandle = []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
	// 连接输入管道
	pipeIn := pipename + "_in"
	FSCTLPIPEWAITRequestIn := c.NewFSCTLPIPEWAITRequest(pipeIn)
	FSCTLPIPEWAITRequestIn.Timeout = timeout
	IOCTLRequest.Buffer = FSCTLPIPEWAITRequestIn
	c.Debug("Sending Ioctl stdin pipe request ["+pipeIn+"]", nil)
	buf, err := c.Send(IOCTLRequest)
	if err != nil {
		c.Debug("", err)
		//return nil, err
	}
	res := NewIOCTLResponse()
	c.Debug("Unmarshalling Ioctl stdin pipe response ["+pipeIn+"]", nil)
	if err := encoder.Unmarshal(buf, &res); err != nil {
		c.Debug("Raw:\n"+hex.Dump(buf), err)
	}
	stdIn = res.GUIDHandle
	c.Debug("Completed Ioctl stdin pipe ["+pipeIn+"]", nil)
	// 连接输出管道
	pipeOut := pipename + "_out"
	FSCTLPIPEWAITRequestOut := c.NewFSCTLPIPEWAITRequest(pipeOut)
	FSCTLPIPEWAITRequestOut.Timeout = timeout
	IOCTLRequest.Buffer = FSCTLPIPEWAITRequestOut
	c.Debug("Sending Ioctl stdout pipe request ["+pipeOut+"]", nil)
	buf, err = c.Send(IOCTLRequest)
	if err != nil {
		c.Debug("", err)
		//return nil, err
	}
	res = NewIOCTLResponse()
	c.Debug("Unmarshalling Ioctl stdout pipe response ["+pipeOut+"]", nil)
	if err = encoder.Unmarshal(buf, &res); err != nil {
		c.Debug("Raw:\n"+hex.Dump(buf), err)
	}
	c.Debug("Completed Ioctl stdout pipe ["+pipeOut+"]", nil)
	// 创建
	// 连接错误管道
	pipeErr := pipename + "_err"
	FSCTLPIPEWAITRequestErr := c.NewFSCTLPIPEWAITRequest(pipeErr)
	FSCTLPIPEWAITRequestErr.Timeout = timeout
	IOCTLRequest.Buffer = FSCTLPIPEWAITRequestErr
	c.Debug("Sending Ioctl stderr pipe request ["+pipeErr+"]", nil)
	buf, err = c.Send(IOCTLRequest)
	if err != nil {
		c.Debug("", err)
		//return nil, err
	}
	res = NewIOCTLResponse()
	c.Debug("Unmarshalling Ioctl stderr pipe response ["+pipeErr+"]", nil)
	if err = encoder.Unmarshal(buf, &res); err != nil {
		c.Debug("Raw:\n"+hex.Dump(buf), err)
	}
	c.Debug("Completed Ioctl stderr pipe ["+pipeErr+"]", nil)
	return stdIn, stdOut, stdErr, nil
}
