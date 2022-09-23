package util

import (
	"crypto/rand"
	"encoding/hex"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
)

// 提供一些常用方法

// 计算结构体大小
func SizeOfStruct(data interface{}) int {
	return sizeof(reflect.ValueOf(data))
}

func sizeof(v reflect.Value) int {
	var sum int
	switch v.Kind() {
	case reflect.Map:
		sum = 0
		keys := v.MapKeys()
		for i := 0; i < len(keys); i++ {
			mapkey := keys[i]
			s := sizeof(mapkey)
			if s < 0 {
				return -1
			}
			sum += s
			s = sizeof(v.MapIndex(mapkey))
			if s < 0 {
				return -1
			}
			sum += s
		}
	case reflect.Slice, reflect.Array:
		sum = 0
		for i, n := 0, v.Len(); i < n; i++ {
			s := sizeof(v.Index(i))
			if s < 0 {
				return -1
			}
			sum += s
		}
	case reflect.String:
		sum = 0
		for i, n := 0, v.Len(); i < n; i++ {
			s := sizeof(v.Index(i))
			if s < 0 {
				return -1
			}
			sum += s
		}
	case reflect.Struct:
		sum = 0
		for i, n := 0, v.NumField(); i < n; i++ {
			s := sizeof(v.Field(i))
			if s < 0 {
				return -1
			}
			sum += s
		}
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128,
		reflect.Int:
		sum = int(v.Type().Size())
	default:
		return 0
	}
	return sum
}

// 读文件
func ReadFile(filename string) ([]byte, error) {
	fp, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	buf, err := ioutil.ReadAll(fp)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

// 处理PDU uuid 转成字节数组
func PDUUuidFromBytes(uuid string) []byte {
	s := strings.ReplaceAll(uuid, "-", "")
	b, _ := hex.DecodeString(s)
	r := []byte{b[3], b[2], b[1], b[0], b[5], b[4], b[7], b[6], b[8], b[9], b[10], b[11], b[12], b[13], b[14], b[15]}
	return r
}

func Random(n int) []byte {
	const alpha = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alpha[b%byte(len(alpha))]
	}
	return bytes
}
