package plugins

import (
	"github.com/soniah/gosnmp"
	"time"
)

func SnmpCrack(serv *Service) int {
	gosnmp.Default.Target = serv.Ip
	gosnmp.Default.Port = uint16(serv.Port)
	gosnmp.Default.Community = serv.Pass
	gosnmp.Default.Timeout = time.Duration(serv.Timeout)
	err := gosnmp.Default.Connect()
	if err == nil {
		oids := []string{"1.3.6.1.2.1.1.4.0", "1.3.6.1.2.1.1.7.0"}
		_, err = gosnmp.Default.Get(oids)
		if err == nil {
			return CrackSuccess
		}
		return CrackFail
	}
	return CrackFail
}
