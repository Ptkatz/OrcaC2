package execopt

import (
	"Orca_Master/define/colorcode"
	"fmt"
	"strings"
)

var NeedPidMap = make(map[string]bool)

func InitNeedPidMap() {
	NeedPidMap["CreateThread"] = false
	NeedPidMap["CreateRemoteThread"] = true
	NeedPidMap["RtlCreateUserThread"] = true
	NeedPidMap["EtwpCreateEtwThread"] = false
}

func JudgeLoadType(loadFunc string, pid int) int {
	InitNeedPidMap()
	for i, v := range NeedPidMap {
		i = strings.ToLower(i)
		if i == loadFunc {
			if v == false {
				return 0
			}
			if pid != 0 {
				return 0
			}
			//stdErr := fmt.Sprintf("'%s' load method needs to specify pid", loadFunc)
			//colorcode.PrintMessage(colorcode.SIGN_ERROR, stdErr)
			return 1
		}
	}
	stdErr := fmt.Sprintf("'%s' load method does not exist", loadFunc)
	colorcode.PrintMessage(colorcode.SIGN_ERROR, stdErr)
	return -1
}
