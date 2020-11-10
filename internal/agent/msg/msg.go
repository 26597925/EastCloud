package msg

import (
	"bytes"
	"encoding/binary"
)

const (
	Heartbeat = iota
	Login //登录成功后，配置的客户端获取
	Command //start,stop,uninstall
	File
	OK
	Online //定时更新数据
	ERROR
	ACK
	UNKNOWN
)

const (
	FlagCompress = iota
	FlagAutoAck
)

const (
	FileTransferError = iota
	ConmmandError
)

const (
	Start = iota + 1
	Stop
	Restart
	Reload
	Uninstall
)

type BaseMessage struct {
	ConnID uint32
	SessionId string
	Flags byte  ///特性，如是否加密，是否压缩, 是否ack等
}

func (bm *BaseMessage) addFlag(flag byte) {
	bm.Flags |= flag
}

func (bm *BaseMessage) hasFlag(flag byte) bool {
	return (bm.Flags & flag) != 0
}

func (bm *BaseMessage) WriteBytes(buff *bytes.Buffer, data []byte) error {
	l := uint64(len(data))
	if err := binary.Write(buff, binary.LittleEndian, l); err != nil {
		return err
	}
	if err := binary.Write(buff, binary.LittleEndian, data); err != nil {
		return err
	}

	return nil
}

func (bm *BaseMessage) ReadBytes(dataBuff *bytes.Reader) ([]byte, error) {
	var l uint64
	if err := binary.Read(dataBuff, binary.LittleEndian, &l); err != nil {
		return nil, err
	}

	data := make([]byte, l)
	if err := binary.Read(dataBuff, binary.LittleEndian, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func (bm *BaseMessage) WriteString(buff *bytes.Buffer, str string) error {
	return bm.WriteBytes(buff, []byte(str))
}

func (bm *BaseMessage) ReadString(dataBuff *bytes.Reader) (string, error) {
	data,err := bm.ReadBytes(dataBuff)
	if err != nil {
		return "", err
	}
	return string(data), err
}

func (bm *BaseMessage) BasePack(dataBuff *bytes.Buffer) error {
	if err := bm.WriteString(dataBuff, bm.SessionId); err != nil {
		return err
	}
	if err := binary.Write(dataBuff, binary.LittleEndian, bm.Flags); err != nil {
		return err
	}

	return nil
}

func (bm *BaseMessage) BaseUnpack(dataBuff *bytes.Reader) error {
	var err error
	if bm.SessionId, err = bm.ReadString(dataBuff); err != nil {
		return err
	}

	if err := binary.Read(dataBuff, binary.LittleEndian, &bm.Flags); err != nil {
		return err
	}
	return nil
}

type CommandData struct {
	BaseMessage
	Name string
	Type byte
}

func (cmd  *CommandData) Pack() ([]byte, error) {
	dataBuff := bytes.NewBuffer([]byte{})

	if err := cmd.BasePack(dataBuff); err != nil {
		return nil, err
	}
	if err := cmd.WriteString(dataBuff, cmd.Name); err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuff, binary.LittleEndian, cmd.Type); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

func (cmd  *CommandData) UnPack(binaryData []byte) error {
	dataBuff := bytes.NewReader(binaryData)
	if err := cmd.BaseUnpack(dataBuff); err != nil {
		return err
	}
	var err error
	if cmd.Name, err = cmd.ReadString(dataBuff); err != nil {
		return err
	}
	if err := binary.Read(dataBuff, binary.LittleEndian, &cmd.Type); err != nil {
		return err
	}
	return  nil
}

type FileData struct {
	BaseMessage
	FileName string
	FileSize uint64 //总大小
	ChunkCount uint32 //总块数
	Sign string //md5签名验证文件是否正常

	StartSize uint64 //起始大小
	ChunkSize uint64
	ChunkData string
}

func (file *FileData) Pack() ([]byte, error) {
	dataBuff := bytes.NewBuffer([]byte{})
	if err := file.BasePack(dataBuff); err != nil {
		return nil, err
	}
	if err := file.WriteString(dataBuff, file.FileName); err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuff, binary.LittleEndian, file.FileSize); err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuff, binary.LittleEndian, file.ChunkCount); err != nil {
		return nil, err
	}
	if err := file.WriteString(dataBuff, file.Sign); err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuff, binary.LittleEndian, file.StartSize); err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuff, binary.LittleEndian, file.ChunkSize); err != nil {
		return nil, err
	}
	if err := file.WriteString(dataBuff, file.ChunkData); err != nil {
		return nil, err
	}
	return dataBuff.Bytes(), nil
}

func (file *FileData) UnPack(binaryData []byte) error {
	dataBuff := bytes.NewReader(binaryData)
	if err := file.BaseUnpack(dataBuff); err != nil {
		return err
	}
	var  err error
	if file.FileName, err = file.ReadString(dataBuff); err != nil {
		return err
	}
	if err = binary.Read(dataBuff, binary.LittleEndian, &file.FileSize); err != nil {
		return err
	}
	if err = binary.Read(dataBuff, binary.LittleEndian, &file.ChunkCount); err != nil {
		return err
	}
	if file.Sign, err = file.ReadString(dataBuff); err != nil {
		return err
	}

	if err = binary.Read(dataBuff, binary.LittleEndian, &file.StartSize); err != nil {
		return err
	}
	if err = binary.Read(dataBuff, binary.LittleEndian, &file.ChunkSize); err != nil {
		return err
	}
	if file.ChunkData, err = file.ReadString(dataBuff); err != nil {
		return err
	}
	return  nil
}

type ErrorData struct {
	BaseMessage
	Origin uint8
	Data string
}

func (ed  *ErrorData) Pack() ([]byte, error) {
	dataBuff := bytes.NewBuffer([]byte{})

	if err := ed.BasePack(dataBuff); err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuff, binary.LittleEndian, ed.Origin); err != nil {
		return nil, err
	}
	if err := ed.WriteString(dataBuff, ed.Data); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

func (ed  *ErrorData) UnPack(binaryData []byte) error {
	dataBuff := bytes.NewReader(binaryData)
	if err := ed.BaseUnpack(dataBuff); err != nil {
		return err
	}
	if err := binary.Read(dataBuff, binary.LittleEndian, &ed.Origin); err != nil {
		return err
	}
	var  err error
	if ed.Data, err = ed.ReadString(dataBuff); err != nil {
		return err
	}
	return  nil
}

type OkData struct {
	BaseMessage
	Origin uint8
	Data string
}

func (ok  *OkData) Pack() ([]byte, error) {
	dataBuff := bytes.NewBuffer([]byte{})

	if err := ok.BasePack(dataBuff); err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuff, binary.LittleEndian, ok.Origin); err != nil {
		return nil, err
	}
	if err := ok.WriteString(dataBuff, ok.Data); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

func (ok  *OkData) UnPack(binaryData []byte) error {
	dataBuff := bytes.NewReader(binaryData)
	if err := ok.BaseUnpack(dataBuff); err != nil {
		return err
	}
	if err := binary.Read(dataBuff, binary.LittleEndian, &ok.Origin); err != nil {
		return err
	}
	var  err error
	if ok.Data, err = ok.ReadString(dataBuff); err != nil {
		return err
	}
	return  nil
}

type AckData struct {
	BaseMessage
	Origin uint8
}

func (ack  *AckData) Pack() ([]byte, error) {
	dataBuff := bytes.NewBuffer([]byte{})

	if err := ack.BasePack(dataBuff); err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuff, binary.LittleEndian, ack.Origin); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

func (ack  *AckData) UnPack(binaryData []byte) error {
	dataBuff := bytes.NewReader(binaryData)
	if err := ack.BaseUnpack(dataBuff); err != nil {
		return err
	}
	if err := binary.Read(dataBuff, binary.LittleEndian, &ack.Origin); err != nil {
		return err
	}
	return  nil
}