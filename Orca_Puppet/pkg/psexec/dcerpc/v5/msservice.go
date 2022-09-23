package v5

import (
	"Orca_Puppet/pkg/psexec/encoder"
	"Orca_Puppet/pkg/psexec/ms"
	"Orca_Puppet/pkg/psexec/smb/smb2"
	"encoding/hex"
	"errors"
	"fmt"
)

// 此文件提供访问windows服务管理安装/删除

// 上传文件
func (c *Client) FileUpload(file, Path string) error {
	treeId, err := c.TreeConnect("C$")
	if err != nil {
		c.Debug("", err)
		return err
	}
	createRequestStruct := smb2.CreateRequestStruct{
		OpLock:             smb2.SMB2_OPLOCK_LEVEL_NONE,
		ImpersonationLevel: smb2.Impersonation,
		AccessMask:         smb2.FILE_CREATE,
		FileAttributes:     smb2.FILE_ATTRIBUTE_NORMAL,
		ShareAccess:        smb2.FILE_SHARE_WRITE,
		CreateDisposition:  smb2.FILE_OVERWRITE_IF,
		CreateOptions:      smb2.FILE_NON_DIRECTORY_FILE,
	}
	fileId, err := c.CreateRequest(treeId, file, createRequestStruct)
	if err != nil {
		c.Debug("", err)
		return err
	}
	err = c.WriteRequest(treeId, Path, file, fileId)
	if err != nil {
		c.Debug("", err)
		return err
	}
	// 关闭目录连接
	c.TreeDisconnect("C$")
	return nil
}

// 打开scm，返回scm服务句柄
func (c *Client) OpenSvcManager(treeId uint32) (fileid, handler []byte, err error) {
	createRequestStruct := smb2.CreateRequestStruct{
		OpLock:             smb2.SMB2_OPLOCK_LEVEL_NONE,
		ImpersonationLevel: smb2.Impersonation,
		AccessMask:         smb2.FILE_OPEN_IF,
		FileAttributes:     smb2.FILE_ATTRIBUTE_NORMAL,
		ShareAccess:        smb2.FILE_SHARE_READ,
		CreateDisposition:  smb2.FILE_OPEN,
		CreateOptions:      smb2.FILE_NON_DIRECTORY_FILE,
	}
	fileId, err := c.CreateRequest(treeId, "svcctl", createRequestStruct)
	if err != nil {
		c.Debug("", err)
		return nil, nil, err
	}
	// 绑定svcctl函数
	err = c.PDUBind(treeId, fileId, ms.NTSVCS_UUID, ms.NTSVCS_VERSION)
	if err != nil {
		c.Debug("", err)
		return nil, nil, err
	}
	req := c.NewOpenSCManagerWRequest(treeId, fileId)
	_, err = c.Send(req)
	if err != nil {
		c.Debug("", err)
		return nil, nil, err
	}
	c.Debug("Read svcctl response", nil)
	reqRead := c.NewReadRequest(treeId, fileId)
	buf, err := c.Send(reqRead)
	if err != nil {
		c.Debug("", err)
		return nil, nil, err
	}
	res := NewOpenSCManagerWResponse()
	c.Debug("Unmarshalling OpenSCManagerW response", nil)
	if err = encoder.Unmarshal(buf, &res); err != nil {
		c.Debug("Raw:\n"+hex.Dump(buf), err)
	}
	if res.SMB2Header.Status != ms.STATUS_SUCCESS {
		return nil, nil, errors.New("Failed to OpenSCManagerW service active to " + ms.StatusMap[res.SMB2Header.Status])
	}
	c.Debug("Completed OpenSCManagerW ", nil)
	// 获取OpenSCManagerW句柄
	contextHandle := res.ContextHandle
	return fileId, contextHandle, nil
}

// 打开服务
func (c *Client) OpenService(treeId uint32, fileId, contextHandle []byte, servicename string) error {
	// 打开服务
	c.Debug("Sending svcctl OpenServiceW request", nil)
	req := c.NewROpenServiceWRequest(treeId, fileId, contextHandle, servicename)
	buf, err := c.Send(req)
	if err != nil {
		c.Debug("", err)
		return err
	}
	c.Debug("Read svcctl OpenServiceW response", nil)
	reqRead := c.NewReadRequest(treeId, fileId)
	buf, err = c.Send(reqRead)
	if err != nil {
		c.Debug("", err)
		return err
	}
	res := NewROpenServiceWResponse()
	c.Debug("Unmarshalling ROpenServiceW response", nil)
	if err = encoder.Unmarshal(buf, &res); err != nil {
		c.Debug("Raw:\n"+hex.Dump(buf), err)
	}
	//if res.SMB2Header.Status != ms.STATUS_SUCCESS {
	//	return errors.New("Failed to ROpenServiceW to " + ms.StatusMap[res.SMB2Header.Status])
	//}
	c.Debug("Completed ROpenServiceW ", nil)
	return nil
}

// 创建服务，返回创建服务后的实例句柄
func (c *Client) CreateService(treeId uint32, fileId, contextHandle []byte, servicename, uploadPathFile string) (handler []byte, err error) {
	// 创建服务
	c.Debug("Sending svcctl RCreateServiceW request", nil)
	req := c.NewRCreateServiceWRequest(treeId, fileId, contextHandle, servicename, uploadPathFile)
	buf, err := c.Send(req)
	if err != nil {
		c.Debug("", err)
		return nil, err
	}
	c.Debug("Read svcctl RCreateServiceW response", nil)
	reqRead := c.NewReadRequest(treeId, fileId)
	buf, err = c.Send(reqRead)
	if err != nil {
		c.Debug("", err)
		return nil, err
	}
	res := NewRCreateServiceWResponse()
	c.Debug("Unmarshalling RCreateServiceW response", nil)
	if err = encoder.Unmarshal(buf, &res); err != nil {
		c.Debug("Raw:\n"+hex.Dump(buf), err)
	}
	c.Debug("Completed RCreateServiceW to ["+servicename+"] ", nil)
	// 得到创建服务后的服务句柄
	serviceHandle := res.ContextHandle
	return serviceHandle, nil
}

// 启动服务
func (c *Client) StartService(treeId uint32, fileId, serviceHandle []byte) error {
	// 启动服务
	c.Debug("Sending svcctl RStartServiceW request", nil)
	req := c.NewRStartServiceWRequest(treeId, fileId, serviceHandle)
	buf, err := c.Send(req)
	if err != nil {
		c.Debug("", err)
		return err
	}
	c.Debug("Read svcctl RStartServiceW response", nil)
	reqRead := c.NewReadRequest(treeId, fileId)
	buf, err = c.Send(reqRead)
	if err != nil {
		c.Debug("", err)
		return err
	}
	res := NewRStartServiceWResponse()
	c.Debug("Unmarshalling RStartServiceW response", nil)
	if err = encoder.Unmarshal(buf, &res); err != nil {
		c.Debug("Raw:\n"+hex.Dump(buf), err)
	}
	if res.StubData != ms.STATUS_SUCCESS {
		return errors.New("Failed to RStartServiceW: " + ms.StatusMap[res.StubData])
	}
	c.Debug("Completed RStartServiceW ", nil)
	return nil
}

// 删除服务
func DeleteService() {

}

// 关闭scm句柄
func (c *Client) CloseService(treeId uint32, fileId, serviceHandle []byte) error {
	// 关闭服务管理句柄
	c.Debug("Sending svcctl RCloseServiceHandle request", nil)
	req := c.NewRCloseServiceHandleRequest(treeId, fileId, serviceHandle)
	buf, err := c.Send(req)
	if err != nil {
		c.Debug("", err)
		return err
	}
	c.Debug("Read svcctl RCloseServiceHandle response", nil)
	reqRead := c.NewReadRequest(treeId, fileId)
	buf, err = c.Send(reqRead)
	if err != nil {
		c.Debug("", err)
		return err
	}
	res := NewRCloseServiceHandleResponse()
	c.Debug("Unmarshalling RCloseServiceHandle response", nil)
	if err = encoder.Unmarshal(buf, &res); err != nil {
		c.Debug("Raw:\n"+hex.Dump(buf), err)
	}
	if res.ReturnCode != ms.STATUS_SUCCESS {
		return errors.New("Failed to RCloseServiceHandle to " + ms.StatusMap[res.ReturnCode])
	}
	c.Debug("Completed RCloseServiceHandle ", nil)
	return nil
}

// 服务安装
func (c *Client) ServiceInstall(servicename string, file, path string) (service string, err error) {
	// 上传文件
	err = c.FileUpload(file, path)
	if err != nil {
		fmt.Println("[-]", err)
		return "", err
	}
	//建立ipc$管道
	treeId, err := c.TreeConnect("IPC$")
	if err != nil {
		fmt.Println("[-]", err)
		return "", err
	}
	// 打开服务管理
	svcctlFileId, svcctlHandler, err := c.OpenSvcManager(treeId)
	if err != nil {
		fmt.Println("[-]", err)
		return "", err
	}
	// 打开服务
	err = c.OpenService(treeId, svcctlFileId, svcctlHandler, servicename)
	if err != nil {
		fmt.Println("[-]", err)
		return "", err
	}
	// 创建服务
	uploadFilePath := "%systemdrive%\\" + file
	serviceHandle, err := c.CreateService(treeId, svcctlFileId, svcctlHandler, servicename, uploadFilePath)
	if err != nil {
		fmt.Println("[-]", err)
		return "", err
	}
	// 启动服务
	err = c.StartService(treeId, svcctlFileId, serviceHandle)
	if err != nil {
		fmt.Println("[-]", err)
		return servicename, err
	}
	// 关闭服务管理
	err = c.CloseService(treeId, svcctlFileId, svcctlHandler)
	if err != nil {
		fmt.Println("[-]", err)
		return servicename, err
	}
	return servicename, nil
}
