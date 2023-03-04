package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"path"
	"strings"
)

func GetFileExtension(fileName string) string {
	fileExtension := ""
	fileNameSplitted := strings.Split(fileName, ".")
	if len(fileNameSplitted) > 1 {
		fileExtension = fileNameSplitted[len(fileNameSplitted)-1]
	}
	return fileExtension
}

func EncodeFileName(fileName string, fileExtension string) string {
	encodedFileName := md5.Sum([]byte(fileName))
	return fmt.Sprintf("%s.%s", hex.EncodeToString(encodedFileName[:]), fileExtension)
}

func GetFullFilePath(rootPath string, dataStoragePath string, fileExtension string) string {
	//path := fmt.Sprintf("%s/%s/%s", rootPath, dataStoragePath, fileExtension)
	//reg := regexp.MustCompile("(/)*")
	//return reg.ReplaceAllString(path, "$1")
	return path.Join(rootPath, dataStoragePath, fileExtension)
}
