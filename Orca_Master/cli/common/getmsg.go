package common

import (
	"Orca_Master/define/retcode"
	"encoding/json"
)

func GetHttpRet(message string) HttpRetData {
	var httpRetData HttpRetData
	err := json.Unmarshal([]byte(message), &httpRetData)
	if err != nil {
		return HttpRetData{retcode.FAIL, "", ""}
	}
	return httpRetData
}

func GetHttpRetData(message string) string {
	httpRetData := GetHttpRet(message)
	return httpRetData.Data.(string)
}

func GetHttpRetMsg(message string) string {
	httpRetData := GetHttpRet(message)
	return httpRetData.Msg
}

func GetHttpRetCode(message string) int {
	httpRetData := GetHttpRet(message)
	return httpRetData.Code
}
