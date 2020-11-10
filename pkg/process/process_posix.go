// +build !windows

package process

import (
	"errors"
	"os"
	"os/signal"
	"fmt"
	"strconv"
	"strings"
	"golang.org/x/sys/unix"
)

const sigint = unix.SIGINT
const sigterm = unix.SIGTERM
const sighup = unix.SIGHUP

var cmdStart = []string{"/bin/sh", "-c"}
var procAttrs = &unix.SysProcAttr{Setsid: true}

func FindPid(name string) (int, error) {
	cmdline := fmt.Sprintf("ps -ef | grep '%s ' | grep -v grep | awk '{print $2}'", name)
	info := &ProcInfo{
		Cmdline: cmdline,
	}
	output, err := NewCmd(info).Output()
	if err != nil {
		return 0, err
	}

	fields := strings.Fields(string(output))
	if len(fields) == 0 {
		return 0, errors.New("no data")
	}

	pid, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, err
	}

	return pid, nil
}

func notifyCh() <-chan os.Signal {
	sc := make(chan os.Signal, 10)
	signal.Notify(sc, sigterm, sigint, sighup)
	return sc
}
