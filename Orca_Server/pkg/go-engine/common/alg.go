package common

import (
	"bytes"
	"compress/zlib"
	"crypto/rc4"
	"github.com/google/uuid"
	"io"
)

func CompressData(src []byte) []byte {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	w.Write(src)
	w.Close()
	return b.Bytes()
}

func DeCompressData(src []byte) ([]byte, error) {
	b := bytes.NewReader(src)
	r, err := zlib.NewReader(b)
	if err != nil {
		return nil, err
	}
	var out bytes.Buffer
	io.Copy(&out, r)
	r.Close()
	return out.Bytes(), nil
}

func Rc4(key string, src []byte) ([]byte, error) {
	c, err := rc4.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	dst := make([]byte, len(src))
	c.XORKeyStream(dst, src)
	return dst, nil
}

func Guid() string {
	return uuid.New().String()
}
