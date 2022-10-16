package util

import (
	"crypto/md5"
	"encoding/hex"
	uuid "github.com/satori/go.uuid"
	"github.com/shirou/gopsutil/process"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//GenUUID 生成uuid
func GenUUID() string {
	uuidFunc := uuid.NewV4()
	uuidStr := uuidFunc.String()
	uuidStr = strings.Replace(uuidStr, "-", "", -1)
	uuidByt := []rune(uuidStr)
	return string(uuidByt[8:24])
}

// 获取文件md5值
func GetFileMd5Sum(fileName string) string {
	pFile, _ := os.Open(fileName)
	defer pFile.Close()
	md5a := md5.New()
	io.Copy(md5a, pFile)
	return hex.EncodeToString(md5a.Sum(nil))
}

func Md5(s string) string {
	m := md5.New()
	m.Write([]byte(s))
	return hex.EncodeToString(m.Sum(nil))
}

type Charset string

const (
	UTF8    = Charset("UTF-8")
	GB18030 = Charset("GB18030")
	GBK     = Charset("GBK")
)

func ConvertByte2String(byte []byte, charset Charset) string {
	var str string
	switch charset {
	case GB18030:
		var decodeBytes, _ = simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
		str = string(decodeBytes)
	case GBK:
		var decodeBytes, _ = simplifiedchinese.GBK.NewDecoder().Bytes(byte)
		str = string(decodeBytes)
	case UTF8:
		fallthrough
	default:
		str = string(byte)
	}
	return str
}

func GetExecPath() (string, error) {

	pid := int32(os.Getpid())
	processes, err := process.Processes()
	if err != nil {
		return "", err
	}
	for _, p := range processes {
		if p.Pid == pid {
			return p.Cmdline()
		}
	}
	return "", err

}

func GetRandomProcessName() string {
	path, err := GetExecPath()
	if err != nil {
		return ""
	}
	_, file := filepath.Split(path)
	processes, err := process.Processes()
	if err != nil {
		return ""
	}
	var name string
	for {
		rand.Seed(time.Now().Unix())
		proc := processes[rand.Intn(len(processes))]
		name, err = proc.Name()
		if err != nil {
			return ""
		}
		if file != name {
			break
		}
	}
	return name
}
