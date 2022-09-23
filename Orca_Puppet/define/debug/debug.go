package debug

import "log"

var IsDebug = false

func DebugPrint(out string) {
	if IsDebug {
		log.Println(out)
	}
}
