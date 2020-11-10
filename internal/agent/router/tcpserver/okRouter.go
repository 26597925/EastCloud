package tcpserver

import (
	"github.com/26597925/EastCloud/internal/agent/msg"
	"github.com/26597925/EastCloud/pkg/logger"
	"github.com/26597925/EastCloud/pkg/network/common"
	"github.com/26597925/EastCloud/pkg/network/websocket"
)

type OkRouter struct {
	Web *websocket.Server
	common.BaseRouter
}

func (or *OkRouter) Handle(request *common.Request) {
	logger.Info("recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))

	ok := &msg.OkData{}
	err := ok.UnPack(request.GetData())
	if err != nil {
		logger.Error(err)
	}

	conn, err := or.Web.GetConnMgr().Get(ok.ConnID)
	if err != nil {
		logger.Error(err)
	}

	switch ok.Origin {
	case msg.Command:
		err = conn.SendMsg(msg.Command, []byte(ok.Data))
		break
	case msg.File:
		err = conn.SendMsg(msg.File, []byte(ok.Data))
		break
	}

	if err != nil {
		logger.Error(err)
	}
}