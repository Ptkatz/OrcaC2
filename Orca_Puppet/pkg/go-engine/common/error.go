package common

import (
	"errors"
	"fmt"
	"Orca_Puppet/pkg/go-engine/loggo"
	"runtime"
)

func CrashLog() {
	if r := recover(); r != nil {
		var err error
		switch x := r.(type) {
		case string:
			err = errors.New(x)
		case error:
			err = x
		default:
			err = errors.New("Unknown panic")
		}
		if err != nil {
			loggo.Error("crash %s \n%s", err, DumpStacks())
		}
		panic(err)
	}
}

func DumpStacks() string {
	buf := make([]byte, 16384)
	buf = buf[:runtime.Stack(buf, true)]
	return fmt.Sprintf("=== BEGIN goroutine stack dump ===\n%s\n=== END goroutine stack dump ===", buf)
}
