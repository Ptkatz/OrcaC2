package controller

import (
	"Orca_Master/cli/cmdopt/fileopt"
	"Orca_Master/cli/common"
	"Orca_Master/define/colorcode"
	"Orca_Master/define/config"
	"Orca_Master/define/retcode"
	"Orca_Master/tools/crypto"
	"fmt"
	"github.com/desertbit/grumble"
	"github.com/schollz/progressbar/v3"
	"os"
	"path/filepath"
	"time"
)

// 文件命令
var fileCmd = &grumble.Command{
	Name:  "file",
	Help:  "execute file upload or download",
	Usage: "file upload|download [-h | --help]",
}

// 文件上传
var fileUploadCmd = &grumble.Command{
	Name: "upload",
	Help: "execute file upload",
	Usage: "file upload [-h | --help] [-l | --local local_file] [-r | --remote remote_file]\n" +
		"  Notes: The local flag must be specified.\n  If the remote flag is not specified, the default is the current working directory of client.\n" +
		"  eg: \n   file upload -l \"C:\\file\\upload.png\"\n   file upload -l \"C:\\file\\upload.png\" -r \"C:\\Windows\\temp\\save.png\"",
	Flags: func(f *grumble.Flags) {
		f.String("l", "local", "", "the local file path to upload")
		f.String("r", "remote", "", "uploaded to the remote file path")
	},
	Run: func(c *grumble.Context) error {
		if SelectId == -1 {
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "please select the id first")
			return nil
		}
		localFile := c.Flags.String("local")
		remoteFile := c.Flags.String("remote")
		if !fileopt.IsLocalFileLegal(localFile) {
			return nil
		}
		if remoteFile == "" {
			_, fileName := filepath.Split(localFile)
			remoteFile = fileName
		}

		// 发送文件元信息
		data := fileopt.GetFileMetaInfo(localFile, remoteFile)
		retData := fileopt.SendFileMetaMsg(SelectClientId, data, "fileUpload")
		if retData.Code != retcode.SUCCESS {
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "file upload failed")
			return nil
		}
		// 分片发送文件
		fileopt.SendFileData(SelectClientId, localFile)
		// 接收消息，显示是否发送成功
		select {
		case msg := <-common.DefaultMsgChan:
			outputMsg, _ := crypto.Decrypt(common.GetHttpRetData(msg), []byte(config.AesKey))
			fmt.Println(outputMsg)
		case <-time.After(10 * time.Second):
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "request timed out")
			return nil
		}
		time.Sleep(100 * time.Millisecond)
		return nil
	},
}

// 文件下载
var fileDownloadCmd = &grumble.Command{
	Name: "download",
	Help: "execute file download",
	Usage: "file upload [-h | --help] [-r | remote remote_file] [-l | --local local_file]\n" +
		"  Notes: The remote flag must be specified.\n  If the local flag is not specified, the default is the current working directory of master.\n" +
		"  eg: \n   file download -r \"C:\\file\\download.png\"\n   file upload -r \"C:\\file\\download.png\" -l \"C:\\Windows\\temp\\save.png\"",
	Flags: func(f *grumble.Flags) {
		f.String("r", "remote", "", "download from remote file path")
		f.String("l", "local", "", "download to local file path")
	},
	Run: func(c *grumble.Context) error {
		if SelectId == -1 {
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "please select the id first")
			return nil
		}
		remoteFile := c.Flags.String("remote")
		localFile := c.Flags.String("local")
		if !fileopt.IsRemoteFileLegal(remoteFile) {
			return nil
		}
		if localFile == "" {
			_, fileName := filepath.Split(remoteFile)
			localFile = fileName
		}
		// 发送下载文件请求
		retData := fileopt.SendDownloadRequestMsg(SelectClientId, common.ClientId, remoteFile, localFile)
		if retData.Code != retcode.SUCCESS {
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "download request failed")
			return nil
		}
		// 接收消息，显示是否可以下载文件
		var requestFileMsg string
		select {
		case requestFileMsg = <-common.DefaultMsgChan:
			if common.GetHttpRetCode(requestFileMsg) != retcode.SUCCESS {
				return fmt.Errorf(common.GetHttpRetMsg(requestFileMsg))
			}
		case <-time.After(10 * time.Second):
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "request timed out")
			return nil
		}

		// 下载文件元信息
		fileMetaInfo := fileopt.GetMetaInfo(requestFileMsg)
		sliceNum := fileMetaInfo.SliceNum
		md5sum := fileMetaInfo.Md5sum
		pSaveFile, _ := os.OpenFile(localFile, os.O_CREATE|os.O_RDWR, 0600)
		defer pSaveFile.Close()
		// 循环从管道中获取文件元数据并写入
		bar := progressbar.Default(int64(sliceNum) + 1)
		for i := 0; i < sliceNum+1; i++ {
			var err error
			select {
			case metaData := <-common.FileSliceMsgChan:
				_, err = pSaveFile.Write(metaData)
			case <-time.After(8 * time.Second):
				colorcode.PrintMessage(colorcode.SIGN_FAIL, "file download failed")
				return nil
			}
			bar.Add(1)
			if err != nil {
				break
			}
		}
		// 文件md5校验
		saveFileMd5 := fileopt.GetFileMd5Sum(localFile)
		if md5sum == saveFileMd5 {
			colorcode.PrintMessage(colorcode.SIGN_SUCCESS, "file download success")
		} else {
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "file download failed")
		}
		time.Sleep(100 * time.Millisecond)
		return nil
	},
}
