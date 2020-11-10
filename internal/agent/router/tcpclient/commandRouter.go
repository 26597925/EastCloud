package tcpclient

import (
	"encoding/json"
	"fmt"
	"github.com/26597925/EastCloud/internal/agent/core"
	"github.com/26597925/EastCloud/internal/agent/msg"
	"github.com/26597925/EastCloud/pkg/logger"
	"github.com/26597925/EastCloud/pkg/network/common"
)

type CommandRouter struct {
	Process *core.Process

	cmd *msg.CommandData
	common.BaseRouter

	err error
	pid int
	cmdType string
}

func (cr *CommandRouter) PreHandle(request *common.Request)()  {
	logger.Info("recv from client : msgId=", request.GetMsgID(), ", data=", request.GetData())

	cmd := &msg.CommandData{}
	err := cmd.UnPack(request.GetData())
	if err != nil {
		logger.Error(err)

		m := map[string]interface{}{
			"Status": msg.ConmmandError,
			"Info": "command parse fail",
		}
		by,err := json.Marshal(m)
		if err != nil {
			logger.Error(err)
			return
		}

		errInfo := &msg.ErrorData{
			BaseMessage: msg.BaseMessage{
				SessionId: cmd.SessionId,
			},
			Origin: request.GetMsgID(),
			Data: string(by),
		}

		data, err := errInfo.Pack()
		if err != nil {
			logger.Error(err)
			return
		}

		err = request.GetConnection().SendMsg(msg.ERROR, data)
		if err != nil {
			logger.Error(err)
			return
		}
		return
	}

	cr.cmd = cmd
}

func (cr *CommandRouter) Handle(*common.Request) {
	if cr.cmd == nil {
		return
	}

	var info string
	switch cr.cmd.Type {
	case msg.Start:
		cr.cmdType = "start"
		cr.pid, cr.err = cr.Process.StartProxy("start", cr.cmd.Name)
		info = fmt.Sprint("start ", cr.cmd.Name, ", pid=", cr.pid, ", err=", cr.err)
		break
	case msg.Stop:
		cr.cmdType = "stop"
		cr.pid, cr.err = cr.Process.StartProxy("stop", cr.cmd.Name)
		info = fmt.Sprint("stop ", cr.cmd.Name, ", pid=", cr.pid, ", err=", cr.err)
		break
	case msg.Restart:
		cr.cmdType = "restart"
		cr.pid, cr.err = cr.Process.StartProxy("restart", cr.cmd.Name)
		info = fmt.Sprint("restart ", cr.cmd.Name, ", pid=", cr.pid, ", err=", cr.err)
		break
	}
	logger.Info(info)
}

func (cr *CommandRouter) PostHandle(request *common.Request)()  {
	if cr.err != nil {
		m := map[string]interface{}{
			"Status": msg.ConmmandError,
			"Cmd": cr.cmdType,
			"Pid": cr.pid,
			"Error": cr.err.Error(),
		}

		by,err := json.Marshal(m)
		if err != nil {
			logger.Error(err)
			return
		}

		errInfo := &msg.ErrorData{
			BaseMessage: msg.BaseMessage{
				SessionId: cr.cmd.SessionId,
			},
			Origin: request.GetMsgID(),
			Data: string(by),
		}

		data, err := errInfo.Pack()
		if err != nil {
			logger.Error(err)
			return
		}

		err = request.GetConnection().SendMsg(msg.ERROR, data)
		if err != nil {
			logger.Error(err)
			return
		}
	} else {
		m := map[string]interface{}{
			"Cmd": cr.cmdType,
			"Pid": cr.pid,
		}

		by,err := json.Marshal(m)
		if err != nil {
			logger.Error(err)
			return
		}

		okInfo := &msg.OkData{
			BaseMessage: msg.BaseMessage{
				SessionId: cr.cmd.SessionId,
			},
			Origin:	request.GetMsgID(),
			Data:   string(by),
		}

		data, err := okInfo.Pack()
		if err != nil {
			logger.Error(err)
			return
		}

		err = request.GetConnection().SendMsg(msg.OK, data)
		if err != nil {
			logger.Error(err)
			return
		}
	}
}