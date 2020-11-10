package websocket

import (
	"encoding/json"
	"github.com/26597925/EastCloud/pkg/network/common"
)

type DataPack struct{}

type Msg struct {
	Id 	   byte
	Data   string
	Len    int64
}

func newDataPack() *DataPack {
	return &DataPack{}
}

func (dp *DataPack) GetHeadLen() uint32 {
	return 0
}

func (dp *DataPack) Pack(m *common.Message) ([]byte, error) {
	msg := Msg{
		Id: m.GetMsgId(),
		Data: string(m.GetData()),
		Len: int64(m.GetDataLen()),
	}

	return json.Marshal(msg)
}

func (dp *DataPack) Unpack(binaryData []byte) (*common.Message, error) {

	var m Msg
	json.Unmarshal(binaryData, &m)

	msg := &common.Message{}
	msg.SetMsgId(m.Id)
	msg.SetData([]byte(m.Data))
	msg.SetDataLen(uint32(m.Len))
	return msg, nil
}