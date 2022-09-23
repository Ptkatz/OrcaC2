package shell

import (
	"context"
	"fmt"
	"Orca_Puppet/pkg/go-engine/loggo"
	"os/exec"
	"path/filepath"
	"time"
)

func Run(script string, silent bool, param ...string) (string, error) {

	script = filepath.Clean(script)
	script = filepath.ToSlash(script)

	if !silent {
		loggo.Info("shell Run start %v %v ", script, fmt.Sprint(param))
	}

	var tmpparam []string
	tmpparam = append(tmpparam, script)
	tmpparam = append(tmpparam, param...)

	begin := time.Now()
	cmd := exec.Command("sh", tmpparam...)
	out, err := cmd.CombinedOutput()
	outstr := string(out)
	if err != nil {
		loggo.Warn("shell Run fail %v %v", cmd.Args, outstr)
		return "", err
	}

	if !silent {
		loggo.Info("shell Run ok %v %v", cmd.Args, time.Now().Sub(begin))
		loggo.Info("%v", outstr)
	}

	return outstr, nil
}

func RunTimeout(script string, silent bool, timeout int, param ...string) (string, error) {

	script = filepath.Clean(script)
	script = filepath.ToSlash(script)

	d := time.Now().Add(time.Duration(timeout) * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), d)

	defer cancel() // releases resources if slowOperation completes before timeout elapses
	if !silent {
		loggo.Info("shell Run start %v %v %v ", script, timeout, fmt.Sprint(param))
	}

	var tmpparam []string
	tmpparam = append(tmpparam, script)
	tmpparam = append(tmpparam, param...)

	begin := time.Now()
	cmd := exec.CommandContext(ctx, "sh", tmpparam...)
	out, err := cmd.CombinedOutput()
	outstr := string(out)
	if err != nil {
		loggo.Warn("shell Run fail %v %v %v", cmd.Args, outstr, ctx.Err())
		return "", err
	}

	if !silent {
		loggo.Info("shell Run ok %v %v", cmd.Args, time.Now().Sub(begin))
		loggo.Info("%v", outstr)
	}
	return outstr, nil
}

func RunCommand(command string, silent bool) (string, error) {

	if !silent {
		loggo.Info("shell RunCommand start %v ", command)
	}

	begin := time.Now()
	cmd := exec.Command("bash", "-c", command)
	out, err := cmd.CombinedOutput()
	outstr := string(out)
	if err != nil {
		loggo.Warn("shell RunCommand fail %v %v %v", cmd.Args, outstr, err)
		return "", err
	}

	if !silent {
		loggo.Info("shell RunCommand ok %v %v", cmd.Args, time.Now().Sub(begin))
		loggo.Info("%v", outstr)
	}

	return outstr, nil
}

func RunExe(exe string, silent bool, param ...string) (string, error) {

	exe = filepath.Clean(exe)
	exe = filepath.ToSlash(exe)

	if !silent {
		loggo.Info("shell Run start %v %v ", exe, fmt.Sprint(param))
	}

	begin := time.Now()
	cmd := exec.Command(exe, param...)
	out, err := cmd.CombinedOutput()
	outstr := string(out)
	if err != nil {
		loggo.Warn("shell Run fail %v %v %v", cmd.Args, outstr, err)
		return "", err
	}

	if !silent {
		loggo.Info("shell Run ok %v %v", cmd.Args, time.Now().Sub(begin))
		loggo.Info("%v", outstr)
	}

	return outstr, nil
}

func RunExeTimeout(exe string, silent bool, timeout int, param ...string) (string, error) {

	d := time.Now().Add(time.Duration(timeout) * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), d)

	defer cancel() // releases resources if slowOperation completes before timeout elapses

	exe = filepath.Clean(exe)
	exe = filepath.ToSlash(exe)

	if !silent {
		loggo.Info("shell Run start %v %v ", exe, fmt.Sprint(param))
	}

	begin := time.Now()
	cmd := exec.CommandContext(ctx, exe, param...)
	out, err := cmd.CombinedOutput()
	outstr := string(out)
	if err != nil {
		loggo.Warn("shell Run fail %v %v %v", cmd.Args, outstr, err)
		return "", err
	}

	if !silent {
		loggo.Info("shell Run ok %v %v", cmd.Args, time.Now().Sub(begin))
		loggo.Info("%v", outstr)
	}

	return outstr, nil
}

func RunExeRaw(exe string, silent bool, param ...string) string {

	exe = filepath.Clean(exe)
	exe = filepath.ToSlash(exe)

	if !silent {
		loggo.Info("shell Run start %v %v ", exe, fmt.Sprint(param))
	}

	begin := time.Now()
	cmd := exec.Command(exe, param...)
	out, _ := cmd.CombinedOutput()
	outstr := string(out)

	if !silent {
		loggo.Info("shell Run ok %v %v", cmd.Args, time.Now().Sub(begin))
		loggo.Info("%v", outstr)
	}

	return outstr
}
