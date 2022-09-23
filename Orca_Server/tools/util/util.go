package util

import (
	"crypto/md5"
	"encoding/hex"
	uuid "github.com/satori/go.uuid"
	"strings"
)

//GenUUID 生成uuid
func GenUUID() string {
	uuidFunc := uuid.NewV4()
	uuidStr := uuidFunc.String()
	uuidStr = strings.Replace(uuidStr, "-", "", -1)
	uuidByt := []rune(uuidStr)
	return string(uuidByt[8:24])
}

//生成clientId
func GenClientId() string {
	//raw := []byte(setting.GlobalSetting.LocalHost + ":" + setting.CommonSetting.HttpPort)
	//str, err := crypto.Encrypt(raw, []byte(setting.CommonSetting.CryptoKey))
	//if err != nil {
	//	panic(err)
	//}
	uuidFunc := uuid.NewV4()
	uuidStr := uuidFunc.String()
	return uuidStr[:18]
}

func GenGroupKey(systemId, groupName string) string {
	return systemId + ":" + groupName
}

func Md5(s string) string {
	m := md5.New()
	m.Write([]byte(s))
	return hex.EncodeToString(m.Sum(nil))
}
