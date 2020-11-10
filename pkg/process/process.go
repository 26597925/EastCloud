package process

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type ProcInfo struct {
	Pid       int
	Cmdline    string
	Env		   string
}

func NewCmd(proc *ProcInfo) *exec.Cmd {
	cs := append(cmdStart, proc.Cmdline)
	cmd := exec.Command(cs[0], cs[1:]...)
	cmd.SysProcAttr = procAttrs
	env := append(os.Environ(), strings.Split(proc.Env, " ")...)
	cmd.Env = env
	return cmd
}

func SpawnProc(proc *ProcInfo) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		env := append(os.Environ(), strings.Split(proc.Env, "\n")...)
		cmds := strings.Split(proc.Cmdline, "\n")

		cmd = exec.Command(cmds[0], cmds[1:]...)
		cmd.SysProcAttr = procAttrs
		cmd.Env = env
	} else {
		cmd = NewCmd(proc)
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	proc.Pid = cmd.Process.Pid
	return nil
}

func StartProc(proc *ProcInfo) ([]byte, error) {
	return NewCmd(proc).Output()
}

func TerminateProc(proc *ProcInfo, signal os.Signal) error {
	if signal == nil {
		signal = os.Interrupt
	}

	target, err := os.FindProcess(proc.Pid)
	if err != nil {
		return err
	}
	return target.Signal(signal)
}

func StopProc(proc *ProcInfo) error {
	target, err := os.FindProcess(proc.Pid)
	if err != nil {
		return err
	}
	return target.Kill()
}

