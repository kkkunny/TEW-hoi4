package util

import (
	"bytes"
	"os"
)

// WriteFileWithBOM 以utf-8编码格式写入文件
func WriteFileWithBOM(fp string, data []byte) error {
	var dataBuf bytes.Buffer
	_, err := dataBuf.Write([]byte{0xEF, 0xBB, 0xBF})
	if err != nil {
		return err
	}
	_, err = dataBuf.Write(data)
	if err != nil {
		return err
	}
	return os.WriteFile(fp, dataBuf.Bytes(), 0644)
}
