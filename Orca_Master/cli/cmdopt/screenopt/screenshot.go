package screenopt

import (
	"Orca_Master/cli/common"
	"encoding/json"
)

// 屏幕截图元信息
type ScreenShotMetaInfo struct {
	Fid        string // 操作文件ID，随机生成的UUID
	SliceNum   int    // 基础分片数量
	SliceSize  int64  // 基础分片大小
	RemainSize int64  // 剩余分片大小
}

// 发送屏幕截图请求
func SendScreenshotRequestMsg(clientId, sendUserId string) common.HttpRetData {
	msg := "screenshot"
	data := ""
	return common.SendSuccessMsg(clientId, sendUserId, msg, data, "")
}

func GetMetaInfo(metaInfoMsg string) ScreenShotMetaInfo {
	_, _, decData := common.SettleRetData(metaInfoMsg)

	var screenShotMetaInfo ScreenShotMetaInfo
	json.Unmarshal([]byte(decData), &screenShotMetaInfo)
	return screenShotMetaInfo
}
