package web

import (
	"encoding/json"
	"github.com/26597925/EastCloud/internal/agent/msg"
	"github.com/26597925/EastCloud/pkg/logger"
	"github.com/26597925/EastCloud/pkg/network/common"
	"github.com/26597925/EastCloud/pkg/network/tcp"
)

type File struct {
	ConnID uint32
	SessionId string
	Flags byte

	FileName string
	FileSize uint64 //总大小
	ChunkCount uint32 //总块数
	Sign string //md5签名验证文件是否正常

	StartSize uint64 //起始大小
	ChunkSize uint64
	ChunkData string
}

type FileRouter struct {
	Server *tcp.Server
	common.BaseRouter
}

func (fr *FileRouter) Handle(request *common.Request) {
	logger.Info("recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))

	var file File
	err := json.Unmarshal(request.GetData(), &file)
	if err != nil {
		logger.Error(err)
	}

	fl := msg.FileData{
		BaseMessage: msg.BaseMessage{
			ConnID: request.GetConnection().ConnID,
			SessionId: file.SessionId,
			Flags: file.Flags,
		},
		FileName:    file.FileName,
		FileSize:    file.FileSize,
		ChunkCount:  file.ChunkCount,
		Sign:        file.Sign,
		StartSize:   file.StartSize,
		ChunkSize:   file.ChunkSize,
		ChunkData:   file.ChunkData,
	}

	data, err := fl.Pack()
	if err != nil {
		logger.Error(err)
	}

	conn, err := fr.Server.GetConnMgr().Get(file.ConnID)
	if err != nil {
		logger.Error(err)
	}else{
		err = conn.SendMsg(request.GetMsgID(), data)
		if err != nil {
			logger.Error(err)
		}
	}

}