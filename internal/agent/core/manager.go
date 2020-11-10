package core

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/26597925/EastCloud/internal/agent/msg"
	"github.com/26597925/EastCloud/pkg/util/crypto"
	"github.com/26597925/EastCloud/pkg/util/fileext"
	"sync"
)

const (
	NotDownload = iota
	Downloading
	Downloaded
)

type Manager struct {
	sync.RWMutex
	SavePath string
	records map[string]*Record
}

type Record struct {
	Chunk uint32
	FileName string
	CurrentLength uint64
	FileLength uint64
	Sign string //md5签名验证文件是否正常
	Status int
}

func NewManager(savePath string) *Manager {
	return &Manager {
		SavePath: savePath,
		records:make(map[string]*Record),
	}
}

func (m *Manager) ListFile() map[string]*Record {
	return m.records
}

func (m *Manager) AddFile(file *msg.FileData) (string, error) {
	data, err := hex.DecodeString(file.ChunkData)
	if err != nil {
		return "", err
	}

	m.Lock()
	fileSize, err := fileext.WriteOffsetFile(m.SavePath + file.FileName, int64(file.StartSize), data)
	m.Unlock()
	if err != nil {
		return "", err
	}

	record, ok := m.records[file.SessionId]
	if ok {
		record.Chunk++
		record.CurrentLength = fileSize
	} else {
		record = &Record{
			Chunk: 1,
			FileName:      file.FileName,
			CurrentLength: fileSize,
			FileLength:    file.FileSize,
			Sign:          file.Sign,
			Status:        NotDownload,
		}
	}
	record.Status = Downloading
	m.records[file.SessionId] = record

	if record.Chunk == file.ChunkCount &&
		record.FileLength == record.CurrentLength {
		md5sign, err := crypto.Md5File(m.SavePath + file.FileName)
		if err != nil {
			return "", err
		}

		if file.Sign != md5sign {
			return "", errors.New("core verify fail")
		}

		record.Status = Downloaded
	}

	inf, err := json.Marshal(record)
	if err != nil {
		return "", err
	}
	return string(inf), nil
}
