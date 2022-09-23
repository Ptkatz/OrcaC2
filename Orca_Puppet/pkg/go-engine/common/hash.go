package common

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/OneOfOne/xxhash"
	"hash/crc32"
	"strconv"
)

func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func GetXXHashString(s string) string {
	h := xxhash.New64()
	h.WriteString(s)
	return strconv.FormatUint(h.Sum64(), 10)
}

func GetCrc32String(s string) string {
	hash := crc32.New(crc32.IEEETable)
	hash.Write([]byte(s))
	hashInBytes := hash.Sum(nil)[:]
	return hex.EncodeToString(hashInBytes)
}

func GetCrc32(data []byte) string {
	hash := crc32.New(crc32.IEEETable)
	hash.Write(data)
	hashInBytes := hash.Sum(nil)[:]
	return hex.EncodeToString(hashInBytes)
}
