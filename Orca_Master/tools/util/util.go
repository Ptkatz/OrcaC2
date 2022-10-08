package util

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	uuid "github.com/satori/go.uuid"
	"os"
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

func ReadLines(filename string) (lines []string, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		lines = append(lines, s.Text())
	}
	return
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) || err != nil || info == nil {
		return false
	}
	return !info.IsDir()
}

func RemoveDuplicate(list []string) []string {
	var set []string
	hashSet := make(map[string]struct{})
	for _, v := range list {
		hashSet[v] = struct{}{}
	}
	for k := range hashSet {
		// 去除空字符串
		if k == "" {
			continue
		}
		set = append(set, k)
	}
	return set
}

func Md5(s string) string {
	m := md5.New()
	m.Write([]byte(s))
	return hex.EncodeToString(m.Sum(nil))
}
