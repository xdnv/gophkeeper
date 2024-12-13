package domain

import (
	"fmt"
	"os"
	"path/filepath"
)

// Secret structure for datatype "binary"
type KeeperBinary struct {
	FileName  string `db:"file_name" json:"file_name"` // Binary file name
	Extension string `db:"extension" json:"extension"` // Binary file extension
	FileSize  int64  `db:"file_size" json:"file_size"` // Binary file size
	Data      []byte `db:"data" json:"data"`           // Binary file data
}

// String representation for "binary" datatype
func (k KeeperBinary) ToString() string {
	return fmt.Sprintf("File name:%s\nExtension:%s\nFile size: %v", k.FileName, k.Extension, k.FileSize)
}

// Create new KeepBinary and read file from disk
func NewBinarySecret(filePath string) (*KeeperBinary, error) {
	k := new(KeeperBinary)

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return k, err
	}

	k.FileName = fileInfo.Name()
	k.Extension = filepath.Ext(k.FileName)
	k.FileSize = fileInfo.Size()
	k.Data, err = os.ReadFile(filePath)
	if err != nil {
		return k, err
	}
	return k, nil
}

// Dump binary file from
func DumpBinary(k *KeeperBinary, filePath string) error {
	if k.FileSize == 0 || len(k.Data) == 0 {
		return fmt.Errorf("empty binary data storage, nothing to save")
	}
	return os.WriteFile(filePath, k.Data, 0666)
}
