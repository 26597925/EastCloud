package web

import (
	"encoding/json"
	"github.com/26597925/EastCloud/internal/agent/msg"
	"github.com/26597925/EastCloud/pkg/logger"
	"github.com/26597925/EastCloud/pkg/network/common"
	"github.com/26597925/EastCloud/pkg/network/tcp"
)

type Command struct {
	ConnID uint32
	SessionId string
	Flags byte
	Name string
	Type byte
}

type CommandRouter struct {
	Server *tcp.Server
	common.BaseRouter
}

func (wcr *CommandRouter) Handle(request *common.Request) {
	logger.Info("recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))

	var command Command
	err := json.Unmarshal(request.GetData(), &command)
	if err != nil {
		logger.Error(err)
	}

	cmd := &msg.CommandData{
		BaseMessage:msg.BaseMessage{
			ConnID: request.GetConnection().ConnID,
			SessionId: command.SessionId,
			Flags: command.Flags,
		},
		Name: command.Name,
		Type: command.Type,
	}
	data, err := cmd.Pack()
	if err != nil {
		logger.Error(err)
	}

	conn, err := wcr.Server.GetConnMgr().Get(command.ConnID)
	if err != nil {
		logger.Error(string(request.GetData()))
	}else{
		err = conn.SendMsg(request.GetMsgID(), data)
		if err != nil {
			logger.Error(err)
		}
	}
}