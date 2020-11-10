package tcpclient

import (
	"github.com/26597925/EastCloud/pkg/network/common"
)

type HeartbeatRouter struct {
	common.BaseRouter
}

func (hr *HeartbeatRouter) Handle(request *common.Request) {
	//fmt.Println("recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))
}