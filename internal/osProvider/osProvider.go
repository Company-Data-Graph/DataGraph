package osProvider

import (
	"os"
	"path"
)

// Create directory
func CreateDir(root string) {
	os.MkdirAll(root, os.ModePerm)
}

// Check that file current file in current dir exist
func CheckFileExist(filePath string, fileName string) bool {
	if _, err := os.Stat(path.Join(filePath, fileName)); err == nil {
		return true
	}
	return false
}

func CreateFile(fileName string, filePath string, data []byte) error {
	file, err := os.Create(path.Join(filePath, fileName))
	if err != nil {
		return err
	}
	_, err = file.Write(data)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

func DeleteFile(fileName string, filePath string) error {
	if err := os.Remove(path.Join(filePath, fileName)); err != nil {
		return err
	}
	os.Remove(filePath)
	return nil
}
