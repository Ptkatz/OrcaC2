// Copyright (c) 2016 Andreas Auernhammer. All rights reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package skein256

import (
	"bytes"
	"encoding/hex"
	"Orca_Puppet/pkg/go-engine/crypto/cryptonight/inter/skein"
	"testing"
)

func fromHex(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return b
}

var testVectors = []struct {
	hashsize  int
	conf      *skein.Config
	msg, hash string
}{
	{
		hashsize: 32,
		conf:     nil,
		msg:      "",
		hash:     "C8877087DA56E072870DAA843F176E9453115929094C3A40C463A196C29BF7BA",
	},
	{
		hashsize: 32,
		conf:     nil,
		msg:      "FF",
		hash:     "0B98DCD198EA0E50A7A244C444E25C23DA30C10FC9A1F270A6637F1F34E67ED2",
	},
	{
		hashsize: 32,
		conf:     nil,
		msg:      "FFFEFDFCFBFAF9F8F7F6F5F4F3F2F1F0EFEEEDECEBEAE9E8E7E6E5E4E3E2E1E0",
		hash:     "8D0FA4EF777FD759DFD4044E6F6A5AC3C774AEC943DCFC07927B723B5DBF408B",
	},
	{
		hashsize: 32,
		conf:     nil,
		msg: "FFFEFDFCFBFAF9F8F7F6F5F4F3F2F1F0EFEEEDECEBEAE9E8E7E6E5E4E3E2E1E0" +
			"DFDEDDDCDBDAD9D8D7D6D5D4D3D2D1D0CFCECDCCCBCAC9C8C7C6C5C4C3C2C1C0",
		hash: "DF28E916630D0B44C4A849DC9A02F07A07CB30F732318256B15D865AC4AE162F",
	},
	{
		hashsize: 20,
		conf:     nil,
		msg: "FBD17C26B61A82E12E125F0D459B96C91AB4837DFF22B39B78439430CDFC5DC8" +
			"78BB393A1A5F79BEF30995A85A12923339BA8AB7D8FC6DC5FEC6F4ED22C122BB" +
			"E7EB61981892966DE5CEF576F71FC7A80D14DAB2D0C03940B95B9FB3A727C66A" +
			"6E1FF0DC311B9AA21A3054484802154C1826C2A27A0914152AEB76F1168D4410",
		hash: "0CD491B7715704C3A15A45A1CA8D93F8F646D3A1",
	},
	{
		hashsize: 28,
		conf:     nil,
		msg: "FBD17C26B61A82E12E125F0D459B96C91AB4837DFF22B39B78439430CDFC5DC8" +
			"78BB393A1A5F79BEF30995A85A12923339BA8AB7D8FC6DC5FEC6F4ED22C122BB" +
			"E7EB61981892966DE5CEF576F71FC7A80D14DAB2D0C03940B95B9FB3A727C66A" +
			"6E1FF0DC311B9AA21A3054484802154C1826C2A27A0914152AEB76F1168D4410",
		hash: "AFD1E2D0F5B6CD4E1F8B3935FA2497D27EE97E72060ADAC099543487",
	},
	{
		hashsize: 32,
		conf:     nil,
		msg: "FBD17C26B61A82E12E125F0D459B96C91AB4837DFF22B39B78439430CDFC5DC8" +
			"78BB393A1A5F79BEF30995A85A12923339BA8AB7D8FC6DC5FEC6F4ED22C122BB" +
			"E7EB61981892966DE5CEF576F71FC7A80D14DAB2D0C03940B95B9FB3A727C66A" +
			"6E1FF0DC311B9AA21A3054484802154C1826C2A27A0914152AEB76F1168D4410",
		hash: "4DE6FE2BFDAA3717A4261030EF0E044CED9225D066354610842A24A3EAFD1DCF",
	},
	{
		hashsize: 48,
		conf:     nil,
		msg: "FBD17C26B61A82E12E125F0D459B96C91AB4837DFF22B39B78439430CDFC5DC8" +
			"78BB393A1A5F79BEF30995A85A12923339BA8AB7D8FC6DC5FEC6F4ED22C122BB" +
			"E7EB61981892966DE5CEF576F71FC7A80D14DAB2D0C03940B95B9FB3A727C66A" +
			"6E1FF0DC311B9AA21A3054484802154C1826C2A27A0914152AEB76F1168D4410",
		hash: "954620FB31E8B782A2794C6542827026FE069D715DF04261629FCBE81D7D529B" +
			"95BA021FA4239FB00AFAA75F5FD8E78B",
	},
	{
		hashsize: 64,
		conf:     nil,
		msg: "FBD17C26B61A82E12E125F0D459B96C91AB4837DFF22B39B78439430CDFC5DC8" +
			"78BB393A1A5F79BEF30995A85A12923339BA8AB7D8FC6DC5FEC6F4ED22C122BB" +
			"E7EB61981892966DE5CEF576F71FC7A80D14DAB2D0C03940B95B9FB3A727C66A" +
			"6E1FF0DC311B9AA21A3054484802154C1826C2A27A0914152AEB76F1168D4410",
		hash: "51347E27C7EABBA514959F899A6715EF6AD5CF01C23170590E6A8AF399470BF9" +
			"0EA7409960A708C1DBAA90E86389DF254ABC763639BB8CDF7FB663B29D9557C3",
	},
	{
		hashsize: 32,
		conf:     &skein.Config{Key: fromHex("CB41F1706CDE09651203C2D0EFBADDF8")},
		msg:      "",
		hash:     "886E4EFEFC15F06AA298963971D7A25398FFFE5681C84DB39BD00851F64AE29D",
	},
	{
		hashsize: 32,
		conf:     &skein.Config{Key: fromHex("")},
		msg:      "D3090C72167517F7C7AD82A70C2FD3F6443F608301591E59",
		hash:     "DCBD5C8BD09021A840B0EA4AAA2F06E67D7EEBE882B49DE6B74BDC56B60CC48F",
	},
	{
		hashsize: 48,
		conf:     &skein.Config{Key: fromHex("CB41F1706CDE09651203C2D0EFBADDF847A0D315CB2E53FF8BAC41DA0002672E92")},
		msg: "D3090C72167517F7C7AD82A70C2FD3F6443F608301591E598EADB195E8357135" +
			"BA26FEDE2EE187417F816048D00FC23512737A2113709A77E4170C49A94B7FDF" +
			"F45FF579A72287743102E7766C35CA5ABC5DFE2F63A1E726CE5FBD2926DB03A2" +
			"DD18B03FC1508A9AAC45EB362440203A323E09EDEE6324EE2E37B4432C1867ED",
		hash: "96E6CEBB23573D0A70CE36A67AA05D2403148093F25C695E1254887CC97F9771" +
			"D2518413AF4286BF2A06B61A53F7FCEC",
	},
}

func TestVectors(t *testing.T) {
	for i, v := range testVectors {
		conf, msg, ref := v.conf, fromHex(v.msg), fromHex(v.hash)

		h := New(v.hashsize, conf)

		h.Write(msg)
		sum := h.Sum(nil)
		if !bytes.Equal(sum, ref) {
			t.Fatalf("Test vector %d : Hash does not match:\nFound:      %s\nExpected: %s", i, hex.EncodeToString(sum), hex.EncodeToString(ref))
		}

		sum = Sum(msg, v.hashsize, conf)
		if !bytes.Equal(sum, ref) {
			t.Fatalf("Test vector %d : Hash does not match:\nFound:      %s\nExpected: %s", i, hex.EncodeToString(sum), hex.EncodeToString(ref))
		}

		var key []byte
		if conf != nil {
			key = conf.Key
		}

		switch v.hashsize {
		case 64:
			{
				var out [64]byte
				Sum512(&out, msg, key)
				if !bytes.Equal(out[:], ref) {
					t.Fatalf("Test vector %d : Hash does not match:\nFound:      %s\nExpected: %s", i, hex.EncodeToString(out[:]), hex.EncodeToString(ref))
				}
			}
		case 48:
			{
				var out [48]byte
				Sum384(&out, msg, key)
				if !bytes.Equal(out[:], ref) {
					t.Fatalf("Test vector %d : Hash does not match:\nFound:      %s\nExpected: %s", i, hex.EncodeToString(out[:]), hex.EncodeToString(ref))
				}
			}
		case 32:
			{
				var out [32]byte
				Sum256(&out, msg, key)
				if !bytes.Equal(out[:], ref) {
					t.Fatalf("Test vector %d : Hash does not match:\nFound:      %s\nExpected: %s", i, hex.EncodeToString(out[:]), hex.EncodeToString(ref))
				}
			}
		case 20:
			{
				var out [20]byte
				Sum160(&out, msg, key)
				if !bytes.Equal(out[:], ref) {
					t.Fatalf("Test vector %d : Hash does not match:\nFound:      %s\nExpected: %s", i, hex.EncodeToString(out[:]), hex.EncodeToString(ref))
				}
			}
		}
	}
}
