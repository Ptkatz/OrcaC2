package smb2

import (
	"Orca_Puppet/pkg/psexec/encoder"
	"Orca_Puppet/pkg/psexec/ms"
	"Orca_Puppet/pkg/psexec/smb"
	"encoding/hex"
	"errors"
)

// 此文件用于smb2读数据请求

// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-smb2/320f04f3-1b28-45cd-aaa1-9e5aed810dca
type ReadRequestStruct struct {
	smb.SMB2Header
	StructureSize  uint16 //2字节，必须设置49/0x0031
	Padding        uint8
	Flags          uint8
	ReadLength     uint32
	FileOffset     []byte `smb:"fixed:8"`  //8字节
	FileId         []byte `smb:"fixed:16"` //8字节
	MinCount       uint32
	Channel        uint32
	RemainingBytes uint32
	BlobOffset     uint16
	BlobLength     uint16 `smb:"len:Buffer"`
	Buffer         []byte
	Reserved       uint8
}

// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-smb2/3e3d2f2c-0e2f-41ea-ad07-fbca6ffdfd90
type ReadResponseStruct struct {
	smb.SMB2Header
	StructureSize uint16
	DataOffset    uint32 `smb:"fixed:4"`
	Reserved      uint32 `smb:"fixed:4"`
	BlobOffset    uint8
	Reserved2     uint8
	BlobLength    uint32
	//Info          []byte `smb:"count:BlobLength"` //写入的数据
}

type ReadResponseStruct2 struct {
	smb.SMB2Header
	StructureSize uint16
	DataOffset    uint32 `smb:"fixed:4"`
	Reserved      uint32 `smb:"fixed:4"`
	BlobOffset    uint8
	Reserved2     uint8
	BlobLength    uint32
	Info          []byte `smb:"count:BlobLength"` //写入的数据
}

func (c *Client) NewReadRequest(treeId uint32, fileId []byte) ReadRequestStruct {
	smb2Header := NewSMB2Header()
	smb2Header.Command = smb.SMB2_READ
	smb2Header.CreditCharge = 1
	smb2Header.MessageId = c.GetMessageId()
	smb2Header.SessionId = c.GetSessionId()
	smb2Header.TreeId = treeId
	smb2Header.Credits = 127
	return ReadRequestStruct{
		SMB2Header:    smb2Header,
		StructureSize: 49,
		Padding:       0x50,
		ReadLength:    65536,
		FileOffset:    make([]byte, 8),
		FileId:        fileId,
		Channel:       SMB2_CHANNEL_NONE,
		BlobOffset:    0,
		Buffer:        make([]byte, 0),
		Reserved:      48,
	}
}

func NewReadResponse() ReadResponseStruct2 {
	return ReadResponseStruct2{}
}

func (c *Client) ReadRequest(treeId uint32, fileId []byte) (info []byte, err error) {
	c.Debug("Sending Read request", nil)
	req := c.NewReadRequest(treeId, fileId)
	req.ReadLength = 1024
	buf, err := c.Send(req)
	if err != nil {
		c.Debug("", err)
		return nil, err
	}
	res := NewReadResponse()
	c.Debug("Unmarshalling Read response", nil)
	if err = encoder.Unmarshal(buf, &res); err != nil {
		c.Debug("Raw:\n"+hex.Dump(buf), err)
	}
	c.Debug("Raw:\n"+hex.Dump(buf), err)
	if res.SMB2Header.Status != ms.STATUS_SUCCESS {
		return nil, errors.New("Failed to Read response to :" + ms.StatusMap[res.SMB2Header.Status])
	}
	c.Debug("Completed Read response", nil)
	return res.Info, nil
}
