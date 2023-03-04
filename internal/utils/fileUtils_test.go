package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetFileExtension(t *testing.T) {
	assert := assert.New(t)
	testCases := map[string]string{
		"":                      "",
		".":                     "",
		"test":                  "",
		"test.go":               "go",
		"test.txt.go":           "go",
		"/first.second/test.go": "go",
		"/first/second/test.go": "go",
	}
	for key, value := range testCases {
		assert.Equal(value, GetFileExtension(key))
	}
}

type GetFullFilePathTestCase struct {
	rootPath        string
	dataStoragePath string
	fileExtension   string
	expectedResult  string
}

func TestGetFullFilePath(t *testing.T) {
	assert := assert.New(t)
	testCases := []GetFullFilePathTestCase{
		{rootPath: "", dataStoragePath: "", fileExtension: "", expectedResult: ""},
		{rootPath: "root", dataStoragePath: "", fileExtension: "", expectedResult: "root"},
		{rootPath: "root", dataStoragePath: "storage", fileExtension: "", expectedResult: "root/storage"},
		{rootPath: "/root", dataStoragePath: "storage", fileExtension: "", expectedResult: "/root/storage"},
		{rootPath: "/root", dataStoragePath: "/storage", fileExtension: "", expectedResult: "/root/storage"},
		{rootPath: "/root/", dataStoragePath: "/storage", fileExtension: "", expectedResult: "/root/storage"},
		{rootPath: "/root/", dataStoragePath: "/storage", fileExtension: "/test/", expectedResult: "/root/storage/test"},
		{rootPath: "/root/", dataStoragePath: "/storage", fileExtension: "/test", expectedResult: "/root/storage/test"},
		{rootPath: "/root//", dataStoragePath: "//storage//", fileExtension: "//test/", expectedResult: "/root/storage/test"},
	}
	for _, el := range testCases {
		assert.Equalf(el.expectedResult, GetFullFilePath(el.rootPath, el.dataStoragePath, el.fileExtension), "params: '%s' '%s' '%s'", el.rootPath, el.dataStoragePath, el.fileExtension)
	}
}
