package tcp

import (
	"bytes"
	"encoding/binary"
	"github.com/26597925/EastCloud/pkg/network/common"
)

type DataPack struct{
	MaxPacketSize    uint32
}

func newDataPack(maxPacketSize uint32) *DataPack {
	return &DataPack{
		MaxPacketSize: maxPacketSize,
	}
}

func (dp *DataPack) GetHeadLen() uint32 {
	return 5
}

func (dp *DataPack) Pack(msg *common.Message) ([]byte, error) {
	dataBuff := bytes.NewBuffer([]byte{})

	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}

	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

func (dp *DataPack) Unpack(binaryData []byte) (*common.Message, error) {

	dataBuff := bytes.NewReader(binaryData)

	msg := &common.Message{}

	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	//if dp.MaxPacketSize > 0 && msg.DataLen > dp.MaxPacketSize {
	//	return nil, errors.New("too large msg data received")
	//}

	return msg, nil
}
