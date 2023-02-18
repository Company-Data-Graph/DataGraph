package osProvider

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestCreateDir(t *testing.T) {
	testDirName := "test"
	testDir := path.Join(t.TempDir(), testDirName)
	CreateDir(testDir)
	_, err := os.Stat(testDir)
	if os.IsNotExist(err) {
		t.Errorf("Test directory not created!")
	}
}

func TestCheckFileExist(t *testing.T) {
	assert := assert.New(t)
	testFileName := "test.txt"
	testExtensionName := "test"
	tempDir := t.TempDir()
	CreateDir(path.Join(tempDir, testExtensionName))
	file, err := os.Create(path.Join(tempDir, testExtensionName, testFileName))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	assert.Equal(true, CheckFileExist(path.Join(tempDir, testExtensionName), testFileName), "Existed file not found")
	assert.Equal(false, CheckFileExist(path.Join(tempDir, testExtensionName), "not-exist.txt"), "Find file that not exist")
}

func TestCreateFileInNotExistedDirectory(t *testing.T) {
	assert := assert.New(t)
	testFileName := "test.txt"
	testExtensionName := "test"
	tempDir := t.TempDir()
	assert.Error(CreateFile(testFileName, path.Join(tempDir, testExtensionName), nil))
}

func TestCreateFileWith(t *testing.T) {
	assert := assert.New(t)
	testFileName := "test.txt"
	testExtensionName := "test"
	testData := "test data raw"
	tempDir := t.TempDir()
	testPath := path.Join(tempDir, testExtensionName)
	CreateDir(testPath)
	assert.NoError(CreateFile(testFileName, path.Join(tempDir, testExtensionName), []byte(testData)))
	assert.Equal(true, CheckFileExist(testPath, testFileName), "Created file not found")
	fileData, _ := ioutil.ReadFile(path.Join(testPath, testFileName))
	assert.Equal(testData, string(fileData), "Wrong data in the file")
}

func TestDeleteNotExistedFile(t *testing.T) {
	assert := assert.New(t)
	testFileName := "test.txt"
	testExtensionName := "test"
	tempDir := t.TempDir()
	testPath := path.Join(tempDir, testExtensionName)
	assert.Error(DeleteFile(testFileName, testPath), "Cant delete not existed file")
}

func TestDeleteFile(t *testing.T) {
	assert := assert.New(t)
	testFileName := "test.txt"
	testExtensionName := "test"
	tempDir := t.TempDir()
	testPath := path.Join(tempDir, testExtensionName)
	CreateDir(testPath)
	CreateFile(testFileName, testPath, nil)
	assert.NoError(DeleteFile(testFileName, testPath), "File exist but not deleted")
}
