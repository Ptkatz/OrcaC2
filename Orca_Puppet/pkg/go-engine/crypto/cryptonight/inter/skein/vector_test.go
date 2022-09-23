// Copyright (c) 2016 Andreas Auernhammer. All rights reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package skein

import (
	"bytes"
	"encoding/hex"
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
	conf      *Config
	msg, hash string
}{
	{
		hashsize: 64,
		conf:     nil,
		msg:      "",
		hash: "BC5B4C50925519C290CC634277AE3D6257212395CBA733BBAD37A4AF0FA06AF4" +
			"1FCA7903D06564FEA7A2D3730DBDB80C1F85562DFCC070334EA4D1D9E72CBA7A",
	},
	{
		hashsize: 64,
		conf:     nil,
		msg: "FBD17C26B61A82E12E125F0D459B96C91AB4837DFF22B39B78439430CDFC5DC8" +
			"78BB393A1A5F79BEF30995A85A12923339BA8AB7D8FC6DC5FEC6F4ED22C122BB",
		hash: "02D01535C2DF280FDE92146DF054B0609273C73056C93B94B82F5E7DCC5BE697" +
			"9978C4BE24331CAA85D892D2E710C6C9B4904CD056A53547B866BEE097C0FB17",
	},
	{
		hashsize: 20,
		conf:     nil,
		msg: "FBD17C26B61A82E12E125F0D459B96C91AB4837DFF22B39B78439430CDFC5DC8" +
			"78BB393A1A5F79BEF30995A85A12923339BA8AB7D8FC6DC5FEC6F4ED22C122BB" +
			"E7EB61981892966DE5CEF576F71FC7A80D14DAB2D0C03940B95B9FB3A727C66A" +
			"6E1FF0DC311B9AA21A3054484802154C1826C2A27A0914152AEB76F1168D4410",
		hash: "EF03079D61B57C6047E15FA2B35B46FA24279539",
	},
	{
		hashsize: 32,
		conf:     nil,
		msg: "FBD17C26B61A82E12E125F0D459B96C91AB4837DFF22B39B78439430CDFC5DC8" +
			"78BB393A1A5F79BEF30995A85A12923339BA8AB7D8FC6DC5FEC6F4ED22C122BB" +
			"E7EB61981892966DE5CEF576F71FC7A80D14DAB2D0C03940B95B9FB3A727C66A" +
			"6E1FF0DC311B9AA21A3054484802154C1826C2A27A0914152AEB76F1168D4410",
		hash: "809DD3F763A11AF90912BBB92BC0D94361CBADAB10142992000C88B4CEB88648",
	},
	{
		hashsize: 48,
		conf:     nil,
		msg: "FBD17C26B61A82E12E125F0D459B96C91AB4837DFF22B39B78439430CDFC5DC8" +
			"78BB393A1A5F79BEF30995A85A12923339BA8AB7D8FC6DC5FEC6F4ED22C122BB" +
			"E7EB61981892966DE5CEF576F71FC7A80D14DAB2D0C03940B95B9FB3A727C66A" +
			"6E1FF0DC311B9AA21A3054484802154C1826C2A27A0914152AEB76F1168D4410",
		hash: "825F5CBD5DA8807A7B4D3E7BD9CD089CA3A256BCC064CD73A9355BF3AE67F2BF" +
			"93AC7074B3B19907A0665BA3A878B262",
	},
	{
		hashsize: 64,
		conf:     nil,
		msg: "FBD17C26B61A82E12E125F0D459B96C91AB4837DFF22B39B78439430CDFC5DC8" +
			"78BB393A1A5F79BEF30995A85A12923339BA8AB7D8FC6DC5FEC6F4ED22C122BB" +
			"E7EB61981892966DE5CEF576F71FC7A80D14DAB2D0C03940B95B9FB3A727C66A" +
			"6E1FF0DC311B9AA21A3054484802154C1826C2A27A0914152AEB76F1168D4410",
		hash: "1A0D5ABF4432E7C612D658F8DCFA35B0D1AB68B8D6BD4DD115C23CC57B5C5BCD" +
			"DE9BFF0ECE4208596E499F211BC07594D0CB6F3C12B0E110174B2A9B4B2CB6A9",
	},
	{
		hashsize: 128,
		conf:     nil,
		msg: "FBD17C26B61A82E12E125F0D459B96C91AB4837DFF22B39B78439430CDFC5DC8" +
			"78BB393A1A5F79BEF30995A85A12923339BA8AB7D8FC6DC5FEC6F4ED22C122BB" +
			"E7EB61981892966DE5CEF576F71FC7A80D14DAB2D0C03940B95B9FB3A727C66A" +
			"6E1FF0DC311B9AA21A3054484802154C1826C2A27A0914152AEB76F1168D4410",
		hash: "8C25D314110D1C0D58054C96A19D571E26A45D5362AA8F47547E53E0BE4A830A" +
			"5F2C29CCD88E2185FEBAD024A4696F2DBE8307DC150E7A58B3793B1A93FAE252" +
			"3E2D239C59A23A1CC127B3C481A9809162E60B4CB01C011B9630322C8FE9745D" +
			"56D0F3AED54B3490578DB4692901EAFC1960C15359176A9C0990B32B8CA8F94B",
	},
	{
		hashsize: 64,
		conf:     &Config{Key: fromHex("")},
		msg:      "D3090C72",
		hash: "1259AFC2CB025EEF2F681E128F889BBCE57F9A502D57D1A17239A12E71603559" +
			"16B72223790FD9A8B367EC96212A3ED239331ED72EF3DEB17685A8D5FD75158D",
	},
	{
		hashsize: 64,
		conf: &Config{Key: fromHex("CB41F1706CDE09651203C2D0EFBADDF847A0D315CB2E53FF8BAC41DA0002672E" +
			"920244C66E02D5F0DAD3E94C42BB65F0D14157DECF4105EF5609D5B0984457C1")},
		msg: "D3090C72167517F7C7AD82A70C2FD3F6",
		hash: "478D7B6C0CC6E35D9EBBDEDF39128E5A36585DB6222891692D1747D401DE34CE" +
			"3DB6FCBAB6C968B7F2620F4A844A2903B547775579993736D2493A75FF6752A1",
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
