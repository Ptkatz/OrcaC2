package config

import (
	"Orca_Puppet/tools/util"
	"encoding/base64"
)

var (
	AesKey = "Adba723b4fe06819"
)

const (
	Version = "0.10.8"
	sysver  = "OrcaC2_" + Version
)

var SystemId = util.Md5(base64.StdEncoding.EncodeToString([]byte(sysver)))
