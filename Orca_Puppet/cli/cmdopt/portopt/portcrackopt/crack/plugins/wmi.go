package plugins

import (
	"github.com/niudaii/crack/pkg/crack/plugins/wmiexec"
	"strconv"
	"strings"
)

func WmiCrack(serv *Service) int {
	err := wmiexec.WMIExec(serv.Ip+":"+strconv.Itoa(serv.Port), serv.User, serv.Pass, "", "", "", "", serv.Timeout, nil)
	if err != nil {
		if strings.Contains(err.Error(), "timeout") {
			return CrackError
		}
		return CrackFail
	}
	return CrackSuccess
}
