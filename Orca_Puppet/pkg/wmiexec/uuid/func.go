package uuid

import (
	"encoding/hex"
	"strings"
)

func FromString(s string) (UUID, error) {

	s = strings.ReplaceAll(s, "-", "")
	b, err := hex.DecodeString(s)
	if err != nil {
		return UUID{}, err
	}
	//8a 88 5d 04 1c eb-11 c9-9f e8 08 00 2b 10 48 60
	//0  1  2  3  4  5  6  7  8  9  10 11 12 13 14 15
	//8a885d04-1ceb-11c9-9fe8-08002b104860

	r := UUID{b[3], b[2], b[1], b[0], b[5], b[4], b[7], b[6], b[8], b[9], b[10], b[11], b[12], b[13], b[14], b[15]}
	//0x04, 0x5d, 0x88, 0x8a, 0xeb, 0x1c, 0xc9, 0x11, 0x9f, 0xe8, 0x08, 0x00, 0x2b, 0x10, 0x48, 0x60
	//3     2     1     0     5     4     7     6     8     9     10    11    12   13     14    15
	return r, nil
}

func fromStringInternalOnly(s string) UUID {
	b, e := FromString(s)
	if e != nil {
		panic(e)
	}
	return b
}

func FromBytes(b []byte) string {
	tmp := [16]byte{b[3], b[2], b[1], b[0], b[5], b[4], b[7], b[6], b[8], b[9], b[10], b[11], b[12], b[13], b[14], b[15]}
	sb := strings.Builder{}
	s := hex.EncodeToString(tmp[:])
	sb.WriteString(s[0:8])
	sb.WriteString("-")
	sb.WriteString(s[8:12])
	sb.WriteString("-")
	sb.WriteString(s[12:16])
	sb.WriteString("-")
	sb.WriteString(s[16:20])
	sb.WriteString("-")
	sb.WriteString(s[20:])
	return sb.String()
}
