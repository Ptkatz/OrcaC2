package controller

import (
	"Orca_Master/cli/cmdopt/fileopt"
	"Orca_Master/cli/cmdopt/screenopt"
	"Orca_Master/cli/common"
	"Orca_Master/define/colorcode"
	"Orca_Master/define/retcode"
	"fmt"
	"github.com/desertbit/grumble"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var screenCmd = &grumble.Command{
	Name:  "screen",
	Help:  "screenshot and screensteam",
	Usage: "screen [-h | --help]",
}

var screenShotCmd = &grumble.Command{
	Name:  "shot",
	Help:  "get a screenshot of the remote host",
	Usage: "screen shot [-h | --help]",
	Run: func(c *grumble.Context) error {
		if SelectId == -1 {
			colorcode.PrintMessage(colorcode.SIGN_ERROR, "please select the id first")
			return nil
		}

		// 发送屏幕截图请求
		retData := screenopt.SendScreenshotRequestMsg(SelectClientId, common.ClientId)
		if retData.Code != retcode.SUCCESS {
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "screenshot request failed")
			return nil
		}
		// 接收消息，显示是否可以接收截图数据
		var requestScreenMsg string
		select {
		case requestScreenMsg = <-common.DefaultMsgChan:
			if common.GetHttpRetCode(requestScreenMsg) != retcode.SUCCESS {
				colorcode.PrintMessage(colorcode.SIGN_FAIL, common.GetHttpRetMsg(requestScreenMsg))
				return nil
			}
		case <-time.After(10 * time.Second):
			colorcode.PrintMessage(colorcode.SIGN_FAIL, "request timed out")
			return nil
		}

		// 保存截图路径
		savePath := "tmp/screenshot"
		if !fileopt.IsDir(savePath) {
			err := os.Mkdir("tmp", 0666)
			err = os.Mkdir(savePath, 0666)
			if err != nil {
				return fmt.Errorf("%s", err)
			}
		}
		timeStr := strconv.FormatInt(time.Now().Unix(), 10)
		temp := savePath + "/" + timeStr + ".jpg"
		saveFile, err := filepath.Abs(temp)
		if err != nil {
			return fmt.Errorf(err.Error())
		}
		screenShotMetaInfo := screenopt.GetMetaInfo(requestScreenMsg)
		sliceNum := screenShotMetaInfo.SliceNum
		pSaveFile, _ := os.OpenFile(saveFile, os.O_CREATE|os.O_RDWR, 0600)
		defer pSaveFile.Close()
		// 循环从管道中获取截图元数据并写入
		for i := 0; i < sliceNum+1; i++ {
			select {
			case metaData := <-common.ScreenSliceMsgChan:
				_, err := pSaveFile.Write(metaData)
				if err != nil {
					break
				}
			case <-time.After(10 * time.Second):
				colorcode.PrintMessage(colorcode.SIGN_FAIL, "request timed out")
				return nil
			}
		}
		colorcode.PrintMessage(colorcode.SIGN_SUCCESS, "The picture has been saved to\n\t"+saveFile)
		time.Sleep(100 * time.Millisecond)
		return nil
	},
}
