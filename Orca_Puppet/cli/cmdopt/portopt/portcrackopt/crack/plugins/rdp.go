package plugins

import (
	"fmt"
	"github.com/niudaii/crack/pkg/crack/plugins/grdp"
	"strings"
)

func RdpCrack(serv *Service) int {
	addr := fmt.Sprintf("%v:%v", serv.Ip, serv.Port)
	err := grdp.Login(addr, "", serv.User, serv.Pass, serv.Timeout)
	if err != nil {
		if strings.Contains(err.Error(), "timeout") {
			return CrackError
		}
		return CrackFail
	}
	return CrackSuccess
}
