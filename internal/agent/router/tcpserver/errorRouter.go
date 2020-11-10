package tcpserver

import (
	"github.com/26597925/EastCloud/internal/agent/msg"
	"github.com/26597925/EastCloud/pkg/logger"
	"github.com/26597925/EastCloud/pkg/network/common"
	"github.com/26597925/EastCloud/pkg/network/websocket"
)

type ErrorRouter struct {
	Web *websocket.Server
	common.BaseRouter
}

func (er *ErrorRouter) Handle(request *common.Request) {
	logger.Info("recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))

	error := &msg.ErrorData{}
	err := error.UnPack(request.GetData())
	if err != nil {
		logger.Error(err)
	}

	conn, err := er.Web.GetConnMgr().Get(error.ConnID)
	if err != nil {
		logger.Error(err)
	}

	switch error.Origin {
	case msg.Command:
		err = conn.SendMsg(msg.Command, []byte(error.Data))
		break
	case msg.File:
		err = conn.SendMsg(msg.File, []byte(error.Data))
		break
	}

	if err != nil {
		logger.Error(err)
		return
	}
}