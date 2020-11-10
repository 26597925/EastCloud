package tcpserver

import (
	"github.com/26597925/EastCloud/pkg/logger"
	"github.com/26597925/EastCloud/pkg/network/common"
)

type HeartbeatRouter struct {
	common.BaseRouter
}

func (hr *HeartbeatRouter) Handle(request *common.Request) {
	logger.Info("recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))
	err := request.GetConnection().SendMsg(request.GetMsgID(), request.GetData())
	if err != nil {
		logger.Error(err)
	}
}