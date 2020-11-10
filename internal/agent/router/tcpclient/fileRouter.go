package tcpclient

import (
	"encoding/json"
	"github.com/26597925/EastCloud/internal/agent/core"
	"github.com/26597925/EastCloud/internal/agent/msg"
	"github.com/26597925/EastCloud/pkg/logger"
	"github.com/26597925/EastCloud/pkg/network/common"
)

type FileRouter struct {
	manager *core.Manager
	common.BaseRouter
}

func NewFileRouter(manager *core.Manager) *FileRouter {
	return &FileRouter{
		manager:    manager,
	}
}

func (fr *FileRouter) Handle(request *common.Request) {
	logger.Info("recv from client : msgId=", request.GetMsgID(), ", data=", request.GetData())

	file := &msg.FileData{}
	file.UnPack(request.GetData())

	info, err := fr.manager.AddFile(file)
	if err != nil {
		logger.Error(err)

		m := map[string]interface{}{
			"Status": msg.FileTransferError,
			"Error": err.Error(),
		}
		by,err := json.Marshal(m)
		if err != nil {
			logger.Error(err)
			return
		}

		msgError := &msg.ErrorData{
			BaseMessage: file.BaseMessage,
			Origin:      request.GetMsgID(),
			Data:        string(by),
		}
		data, err := msgError.Pack()
		err = request.GetConnection().SendMsg(msg.ERROR, data)
		if err != nil {
			logger.Error(err)
		}
		return
	}

	ok := msg.OkData{
		BaseMessage: file.BaseMessage,
		Origin:      request.GetMsgID(),
		Data:        info,
	}
	data, err := ok.Pack()
	if err != nil {
		logger.Error(err)
	}

	err = request.GetConnection().SendMsg(msg.OK, data)
	if err != nil {
		logger.Error(err)
	}
}