package meterpreter

import (
	"encoding/base64"
	"encoding/binary"
	"math/rand"
	"runtime"
	"strings"
	"time"
)

func generateUUID() string {

	//Much of this is based on the documentation on metasploit
	//This will generate the UUID which is used as the path in the URL

	tstamp := make([]byte, 4)
	tstampInt := uint32(time.Now().Unix())
	binary.LittleEndian.PutUint32(tstamp, tstampInt)
	var bitOS int
	var archID int
	if runtime.GOARCH == "amd64" {
		archID = 1
	}
	switch runtime.GOOS {
	case "windows":
		bitOS = 2
	case "darwin":
		bitOS = 9
	case "linux":
		bitOS = 6
	}
	platformXOR := rand.Intn(256)
	archXOR := rand.Intn(256)
	timeXOR := make([]byte, 4)
	binary.LittleEndian.PutUint32(timeXOR, uint32((platformXOR<<24)+(archXOR<<16)+(platformXOR<<8)+archXOR))
	puid := make([]byte, 8)
	rand.Read(puid)

	uuid := []byte{}
	uuid = append(uuid, puid...)
	uuid = append(uuid, byte(archXOR))
	uuid = append(uuid, byte(platformXOR))
	uuid = append(uuid, byte(archXOR)^byte(archID))
	uuid = append(uuid, byte(platformXOR)^byte(bitOS))

	xoredTstamp := make([]byte, 4)
	for i := 0; i < 4; i++ {
		xoredTstamp[i] = timeXOR[i] ^ tstamp[i]
	}
	xoredTstampInt := binary.LittleEndian.Uint32(xoredTstamp)

	uuid = append(uuid, byte(xoredTstampInt>>24)&byte(255))
	uuid = append(uuid, byte(xoredTstampInt>>16)&byte(255))
	uuid = append(uuid, byte(xoredTstampInt>>8)&byte(255))
	uuid = append(uuid, byte(xoredTstampInt)&byte(255))

	encodedUUID := base64.StdEncoding.EncodeToString(uuid)[:16]
	uri := strings.Replace(encodedUUID, "=", "", -1)
	if !strings.Contains(uri, "+") && !strings.Contains(uri, "/") {
		originalUri := uri
		length := rand.Intn(30) + 40 - len(uri)
		var sum int
		for strings.Contains(uri, "+") || strings.Contains(uri, "/") || sum%256 != 92 {
			uri = originalUri
			randBytes := make([]byte, length+2)
			rand.Read(randBytes)
			junk := base64.StdEncoding.EncodeToString(randBytes)[:length]
			uri = uri + junk
			sum = 0
			for i := 0; i < len(uri); i++ {
				sum += int(uri[i])
			}
		}
		return uri
	}
	uri = generateUUID()
	return uri
}
