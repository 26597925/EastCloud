package process

import (
	"errors"
	"fmt"
	"golang.org/x/sys/windows"
	"os"
	"os/signal"
	"strconv"
	"strings"
)

var cmdStart = []string{"cmd", "/c"}
var procAttrs = &windows.SysProcAttr{
	CreationFlags: windows.CREATE_UNICODE_ENVIRONMENT | windows.CREATE_NEW_PROCESS_GROUP,
}

func FindPid(name string) (int, error) {
	cmdline := fmt.Sprintf("tasklist | findstr %s", name)
	info := &ProcInfo{
		Cmdline: cmdline,
	}

	output, err := NewCmd(info).Output()
	if err != nil {
		return 0, err
	}

	fields := strings.Fields(string(output))
	if len(fields) == 6 {
		pid, err := strconv.Atoi(fields[1])
		if err != nil {
			return 0, err
		}
		return  pid, nil
	}

	return 0, errors.New("no data")
}

func notifyCh() <-chan os.Signal {
	sc := make(chan os.Signal, 10)
	signal.Notify(sc, os.Interrupt)
	return sc
}
