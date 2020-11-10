package common

type Message struct {
	DataLen uint32
	Id      byte
	Data    []byte
}

func NewMsg(id byte, data []byte) *Message {
	return &Message{
		DataLen: uint32(len(data)),
		Id:      id,
		Data:    data,
	}
}

func (msg *Message) GetDataLen() uint32 {
	return msg.DataLen
}

func (msg *Message) GetMsgId() byte {
	return msg.Id
}

func (msg *Message) GetData() []byte {
	return msg.Data
}

func (msg *Message) SetDataLen(len uint32) {
	msg.DataLen = len
}

func (msg *Message) SetMsgId(msgId byte) {
	msg.Id = msgId
}

func (msg *Message) SetData(data []byte) {
	msg.Data = data
}
