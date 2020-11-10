package main

import (
	"encoding/json"
	"fmt"
	"github.com/26597925/EastCloud/pkg/process"
	"github.com/26597925/EastCloud/pkg/util/array"
	"os"
	"strconv"
)

const format  = "%d,%s"
var cmd = []string{"start", "stop", "reload","restart"}

//start '{"Pid":0,"Cmdline":"./ttorrent test.torrent", "Env":""}'
func main() {
	args := os.Args
	if len(args) <= 2 {
		errOut(-1, "core params not valid")
	}

	c := args[1]
	if !array.In(cmd, c) {
		errOut(-2, "command params not valid")
	}

	var procInfo process.ProcInfo
	err := json.Unmarshal([]byte(args[2]), &procInfo)
	if err != nil {
		errOut(-3, args[2])
	}

	switch c {
	case "start":
		err = process.SpawnProc(&procInfo)
		break
	case "stop":
		err = process.StopProc(&procInfo)
		break
	case "restart":
		process.StopProc(&procInfo)
		err = process.SpawnProc(&procInfo)
		break
	case "reload":
		err = process.TerminateProc(&procInfo, os.Interrupt)
		break
	}

	if err != nil {
		errOut(-4, err.Error())
	}

	info := fmt.Sprintf(format, 1, strconv.FormatInt(int64(procInfo.Pid), 10))
	_, _ = os.Stdout.WriteString(info)
}

func errOut(status int, info string) {
	data := fmt.Sprintf(format, status, info)
	_, _ = os.Stderr.WriteString(data)
	os.Exit(2)
}