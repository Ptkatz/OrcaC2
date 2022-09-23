package node

import (
	"context"
	"fmt"
	"Orca_Puppet/pkg/go-engine/loggo"
	"os/exec"
	"time"
)

func Run(script string, silent bool, timeout int, param ...string) string {

	d := time.Now().Add(time.Duration(timeout) * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), d)

	defer cancel() // releases resources if slowOperation completes before timeout elapses

	if !silent {
		loggo.Info("node Run start %v %v %v ", script, timeout, fmt.Sprint(param))
	}

	var tmpparam []string
	tmpparam = append(tmpparam, script)
	tmpparam = append(tmpparam, param...)

	begin := time.Now()
	cmd := exec.CommandContext(ctx, "node", tmpparam...)
	out, err := cmd.CombinedOutput()
	outstr := string(out)
	if err != nil {
		loggo.Warn("node Run fail %v %v %v %v", cmd.Args, outstr, ctx.Err(), err)
		return ""
	}

	if !silent {
		loggo.Info("node Run ok %v %v", cmd.Args, time.Now().Sub(begin))
		loggo.Info("%v", outstr)
	}

	return outstr
}
