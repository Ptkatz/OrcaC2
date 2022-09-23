package smb2

import (
	"Orca_Puppet/pkg/psexec/encoder"
	"Orca_Puppet/pkg/psexec/ms"
	"Orca_Puppet/pkg/psexec/smb"
	"encoding/hex"
	"errors"
	"fmt"
)

// 此文件用于目录树连接/断开

// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-smb2/832d2130-22e8-4afb-aafd-b30bb0901798
// 树连接请求结构
type TreeConnectRequestStruct struct {
	smb.SMB2Header
	StructureSize uint16
	Reserved      uint16 //2字节，smb3.x使用，其他忽略
	PathOffset    uint16 `smb:"offset:Path"`
	PathLength    uint16 `smb:"len:Path"`
	Path          []byte
}

// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-smb2/dd34e26c-a75e-47fa-aab2-6efc27502e96
// 树连接响应结构
type TreeConnectResponseStruct struct {
	smb.SMB2Header
	StructureSize uint16
	ShareType     uint8 //1字节，访问共享类型
	Reserved      uint8 //1字节
	ShareFlags    uint32
	Capabilities  uint32
	MaximalAccess uint32
}

// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-smb2/8a622ecb-ffee-41b9-b4c4-83ff2d3aba1b
// 断开树连接请求结构
type TreeDisconnectRequestStruct struct {
	smb.SMB2Header
	StructureSize uint16 //2字节，客户端必须设为4,表示请求大小
	Reserved      uint16
}

// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-smb2/aeac92de-8db3-48f8-a8b7-bfee28b9fd9e
// 断开树连接响应结构
type TreeDisconnectResponseStruct struct {
	smb.SMB2Header
	StructureSize uint16
	Reserved      uint16
}

func (c *Client) NewTreeConnectRequest(name string) (TreeConnectRequestStruct, error) {
	smb2Header := NewSMB2Header()
	smb2Header.Command = smb.SMB2_TREE_CONNECT
	smb2Header.CreditCharge = 1
	smb2Header.MessageId = c.GetMessageId()
	smb2Header.SessionId = c.GetSessionId()
	smb2Header.Credits = 127
	//格式 \\172.20.10.5:445\IPC$
	path := fmt.Sprintf("\\\\%s:%d\\%s", c.GetOptions().Host, c.GetOptions().Port, name)
	return TreeConnectRequestStruct{
		SMB2Header:    smb2Header,
		StructureSize: 9,
		Reserved:      0,
		PathOffset:    0,
		PathLength:    0,
		Path:          encoder.ToUnicode(path),
	}, nil
}

func NewTreeConnectResponse() TreeConnectResponseStruct {
	smb2Header := NewSMB2Header()
	return TreeConnectResponseStruct{
		SMB2Header: smb2Header,
	}
}

func (c *Client) NewTreeDisconnectRequest(treeId uint32) (TreeDisconnectRequestStruct, error) {
	smb2Header := NewSMB2Header()
	smb2Header.Command = smb.SMB2_TREE_DISCONNECT
	smb2Header.CreditCharge = 1
	smb2Header.MessageId = c.GetMessageId()
	smb2Header.SessionId = c.GetSessionId()
	smb2Header.TreeId = treeId
	smb2Header.Credits = 127
	return TreeDisconnectRequestStruct{
		SMB2Header:    smb2Header,
		StructureSize: 4,
		Reserved:      0,
	}, nil
}

func NewTreeDisconnectResponse() TreeDisconnectResponseStruct {
	smb2Header := NewSMB2Header()
	return TreeDisconnectResponseStruct{
		SMB2Header: smb2Header,
	}
}

// 树连接
func (c *Client) TreeConnect(name string) (treeId uint32, err error) {
	c.Debug("Sending TreeConnect request ["+name+"]", nil)
	req, err := c.NewTreeConnectRequest(name)
	if err != nil {
		c.Debug("", err)
		return 0, err
	}
	buf, err := c.Send(req)
	if err != nil {
		c.Debug("", err)
		return 0, err
	}
	res := NewTreeConnectResponse()
	c.Debug("Unmarshalling TreeConnect response ["+name+"]", nil)
	if err = encoder.Unmarshal(buf, &res); err != nil {
		c.Debug("Raw:\n"+hex.Dump(buf), err)
		//return err
	}
	if res.SMB2Header.Status != ms.STATUS_SUCCESS {
		return 0, errors.New("Failed to connect to [" + name + "]: " + ms.StatusMap[res.SMB2Header.Status])
	}
	treeID := res.SMB2Header.TreeId
	trees := make(map[string]uint32)
	trees[name] = treeID
	c.WithTrees(trees)
	c.Debug("Completed TreeConnect ["+name+"]", nil)
	return treeID, nil
}

// 断开树连接
func (c *Client) TreeDisconnect(name string) error {
	var (
		treeid    uint32
		pathFound bool
	)
	trees := c.GetTrees()
	for k, v := range trees {
		if k == name {
			treeid = v
			pathFound = true
			break
		}
	}
	if !pathFound {
		err := errors.New("Unable to find tree path for disconnect")
		c.Debug("", err)
		return err
	}
	c.Debug("Sending TreeDisconnect request ["+name+"]", nil)
	req, err := c.NewTreeDisconnectRequest(treeid)
	if err != nil {
		c.Debug("", err)
		return err
	}
	buf, err := c.Send(req)
	if err != nil {
		c.Debug("", err)
		return err
	}
	c.Debug("Unmarshalling TreeDisconnect response for ["+name+"]", nil)
	res := NewTreeDisconnectResponse()
	if err = encoder.Unmarshal(buf, &res); err != nil {
		c.Debug("Raw:\n"+hex.Dump(buf), err)
		return err
	}
	if res.SMB2Header.Status != ms.STATUS_MORE_PROCESSING_REQUIRED {
		return errors.New("Failed to connect to tree: " + ms.StatusMap[res.SMB2Header.Status])
	}
	delete(trees, name)
	c.WithTrees(trees)
	c.Debug("TreeDisconnect completed ["+name+"]", nil)
	return nil
}
