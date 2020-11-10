package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/26597925/EastCloud/pkg/process"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

const (
	Stop = iota
	StartIng
	Start
	Running
)

type Info struct {
	process.ProcInfo

	Name	   string
	Status	   int

	SavePath  string
}

type Process struct {
	binPath string

	mu sync.Mutex
	procs []*Info
}

func NewProcess(binPath string) *Process {
	return &Process{
		binPath: binPath,
		procs: make([]*Info, 0),
	}
}

func (p *Process) AddProcess(info *Info)  {
	p.procs = append(p.procs, info)
}

func (p *Process) FindProc(name string) *Info {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, proc := range p.procs {
		if proc.Name == name {
			return proc
		}
	}
	return nil
}

func (p *Process) Monitor() error {
	for _, proc := range p.procs {
		if proc.Status != Running {
			continue
		}

		var err error
		if runtime.GOOS == "windows" {
			err = os.ErrNotExist
		} else {
			_, err = os.Stat(fmt.Sprintf("/proc/%d/status", proc.Pid))
			if err == nil {
				continue
			}
		}

		if os.IsNotExist(err) && proc.Status == Running {
			pid, _ := process.FindPid(proc.Name)
			if pid != 0 {
				proc.Pid = pid
				continue
			}

			_, err = p.StartProxy("restart", proc.Name)
			return err
		}
	}

	return nil
}

func (p *Process) StartProxy(types string, name string) (int, error) {
	info := p.FindProc(name)
	if info == nil {
		return -1, errors.New("unknown name: " + name)
	}

	if types == "start" {
		if info.Status >= StartIng {
			return -1, errors.New(name + " is starting")
		}
		info.Status = StartIng
	}

	if types == "stop" {
		info.Status = Stop
	}

	procInfo, err := p.buildProxy(types, info)
	if err != nil {
		return -1, err
	}

	data, err := process.StartProc(procInfo)
	if err != nil {
		info.Status = Stop
		return -1, err
	}

	var res []string
	if runtime.GOOS == "windows" {
		res = strings.Split(string(data), "\n")
		res = strings.Split(res[0], ",")
	} else {
		res = strings.Split(string(data), ",")
	}

	if len(res) < 2 {
		info.Status = Stop
		return -1, errors.New("back result fail")
	}

	status, err := strconv.Atoi(res[0])
	if err != nil {
		info.Status = Stop
		return -1, err
	}

	var pid int
	if status == 1 {
		pid, err = strconv.Atoi(res[1])
		if err != nil {
			info.Status = Stop
			return -1, err
		}
	}

	info.Pid = pid
	if types != "stop" {
		info.Status = Running
	}
	return pid, nil
}

func (p *Process) buildProxy(types string, info *Info) (*process.ProcInfo, error) {
	data, err := json.Marshal(info.ProcInfo)
	if err != nil {
		return nil, err
	}

	var workExec,cmdline string
	if runtime.GOOS == "windows" {
		workExec = strings.TrimRight(p.binPath, "\\") + "\\work.exe"
		cmdline = fmt.Sprintf("%s %s %s", workExec, types, string(data))
	} else {
		workExec = strings.TrimRight(p.binPath, "/") + "/work"
		cmdline = fmt.Sprintf("%s %s '%s'", workExec, types, string(data))
	}
	_, err = os.Stat(workExec)
	if err != nil {
		return nil, err
	}

	procInfo := &process.ProcInfo{
		Cmdline: cmdline,
	}

	return  procInfo, nil
}