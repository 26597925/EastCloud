package fileext

import (
	"io"
	"os"
	"strings"
)

func Format(p string) string {
	parts := strings.Split(p, ".")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return "unknown"
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func WriteOffsetFile(path string, offset int64, data []byte) (uint64, error) {
	file2, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0766)
	if err != nil {
		file2.Close()
		return 0, err
	}

	file2.Seek(offset, io.SeekStart)
	file2.Write(data)
	file2.Sync()
	fileInfo, err := file2.Stat()
	if err != nil {
		file2.Close()
		return 0, err
	}
	fileSize := uint64(fileInfo.Size())
	file2.Close()
	return fileSize, nil
}