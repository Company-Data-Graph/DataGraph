package api

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"media-server/internal/jwt"
	"media-server/internal/models"
	"media-server/internal/utils"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strings"
	"testing"
)

func TestNewMediaAPI(t *testing.T) {
	assert := assert.New(t)
	testCases := []string{
		"testdata/config-1.yaml",
		"testdata/config-2.yaml",
		"testdata/config-3.yaml",
	}
	for _, testCase := range testCases {
		testConfig, _ := models.NewConfigYML(testCase)
		testApi := NewMediaAPI(&testConfig.MediaAPIConfig)
		assert.Equal(testConfig.MediaAPIConfig.Host, testApi.Host)
		assert.Equal(testConfig.MediaAPIConfig.Port, testApi.Port)
		assert.Equal(testConfig.MediaAPIConfig.Prefix, testApi.Prefix)
		assert.Equal(testConfig.MediaAPIConfig.TokenLiveTime, testApi.TokenLiveTime)
		assert.Equal(testConfig.MediaAPIConfig.StorageRootPath, testApi.RootPath)
		assert.Equal(testConfig.MediaAPIConfig.DataStorageRoute, testApi.DataStorageRoute)
	}
}

func SetupHeadersTest(assert *assert.Assertions, rec *httptest.ResponseRecorder) {
	testCases := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Headers": "Content-Type",
	}
	for key, value := range testCases {
		assert.Equal(value, rec.Header().Get(key))
	}
}

func MethodValidationTest(assert *assert.Assertions, pattern string, function func(w http.ResponseWriter, r *http.Request), correct string, testCase string) {
	handler := http.HandlerFunc(function)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(testCase, pattern, nil)
	handler.ServeHTTP(rec, req)
	if correct != testCase {
		assert.Equal(http.StatusMethodNotAllowed, rec.Code)
	} else {
		assert.NotEqual(http.StatusMethodNotAllowed, rec.Code)
	}
}

func TestPing(t *testing.T) {
	assert := assert.New(t)
	testConfig, _ := models.NewConfigYML("testdata/config-1.yaml")
	testApi := NewMediaAPI(&testConfig.MediaAPIConfig)
	expectedResponse := "pong"
	handler := http.HandlerFunc(testApi.ping)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	handler.ServeHTTP(rec, req)
	assert.Equal(expectedResponse, string(rec.Body.Bytes()))
	SetupHeadersTest(assert, rec)
	MethodValidationTest(assert, "/ping", testApi.ping, http.MethodGet, http.MethodGet)
	MethodValidationTest(assert, "/ping", testApi.ping, http.MethodGet, http.MethodPost)
}

type SignInTestCase struct {
	payload      string
	expectedCode int
}

func TestSignIn(t *testing.T) {
	assert := assert.New(t)
	testConfig, _ := models.NewConfigYML("testdata/config-1.yaml")
	testApi := NewMediaAPI(&testConfig.MediaAPIConfig)
	jwt.Users = map[string]string{
		"existed": "test",
	}
	testCases := []SignInTestCase{
		{payload: "{\"Username\": \"existed\", \"Password\": \"test\"}", expectedCode: http.StatusOK},
		{payload: "{\"Username\": \"not-existed\", \"Password\": \"test\"}", expectedCode: http.StatusUnauthorized},
		{payload: "{broken payload", expectedCode: http.StatusBadRequest},
	}
	for _, testCase := range testCases {
		handler := http.HandlerFunc(testApi.signIn)
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/signIn", strings.NewReader(testCase.payload))
		handler.ServeHTTP(rec, req)
		assert.Equal(testCase.expectedCode, rec.Code)
		SetupHeadersTest(assert, rec)
		MethodValidationTest(assert, "/signIn", testApi.signIn, http.MethodPost, http.MethodPost)
		MethodValidationTest(assert, "/signIn", testApi.signIn, http.MethodPost, http.MethodGet)
	}
}

func mockUploadWithoutFil(assert *assert.Assertions, testApi *MediaAPI, token string) {
	handler := http.HandlerFunc(testApi.signIn)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/upload", nil)
	handler.ServeHTTP(rec, req)
	assert.Equal(http.StatusBadRequest, rec.Code)
	SetupHeadersTest(assert, rec)
}

func mockUploadingTest(assert *assert.Assertions, testApi *MediaAPI, token string, filePath string) (*httptest.ResponseRecorder, *http.Request) {
	handler := http.HandlerFunc(testApi.upload)
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	fileWriter, _ := bodyWriter.CreateFormFile("file", filePath)
	fh, _ := os.Open(filePath)
	defer fh.Close()
	io.Copy(fileWriter, fh)
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/upload", bodyBuf)
	req.Header.Add("Token", token)
	req.Header.Add("Content-Type", contentType)
	handler.ServeHTTP(rec, req)
	return rec, req
}

type AuthTestCase struct {
	token   string
	isValid bool
}

func AuthTest(assert *assert.Assertions, method string, pattern string, function func(w http.ResponseWriter, r *http.Request), creds models.Credentials, tokenExpirationTime int) {
	validToken, _ := jwt.GenerateNewToken(creds, tokenExpirationTime)
	testCases := []AuthTestCase{
		{token: validToken, isValid: true},
		{token: "some-unvalid-token", isValid: false},
		{token: "some.unvalid.token", isValid: false},
	}
	for _, testCase := range testCases {
		handler := http.HandlerFunc(function)
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(method, pattern, nil)
		req.Header.Add("Token", testCase.token)
		handler.ServeHTTP(rec, req)
		if testCase.isValid {
			assert.NotEqual(http.StatusUnauthorized, rec.Code)
		} else {
			assert.Equal(http.StatusUnauthorized, rec.Code)
		}
		SetupHeadersTest(assert, rec)
	}

}

func FilesUploadingTest(assert *assert.Assertions, testApi *MediaAPI, token string) {
	mockUploadWithoutFil(assert, testApi, token)
	testCases := []string{
		"testdata/config-1.yaml",
		"testdata/config-2.yaml",
		"testdata/config-3.yaml",
	}
	for _, testCase := range testCases {
		rec, _ := mockUploadingTest(assert, testApi, token, testCase)
		assert.Equal(http.StatusOK, rec.Code)
		assert.Equal(utils.EncodeFileName(path.Base(testCase), utils.GetFileExtension(testCase)), string(rec.Body.Bytes()))
		rec, _ = mockUploadingTest(assert, testApi, token, testCase)
		assert.Equal(http.StatusBadRequest, rec.Code)
		SetupHeadersTest(assert, rec)
	}
}

func TestUpload(t *testing.T) {
	assert := assert.New(t)
	testConfig, _ := models.NewConfigYML("testdata/config-1.yaml")
	testConfig.MediaAPIConfig.StorageRootPath = t.TempDir()
	testApi := NewMediaAPI(&testConfig.MediaAPIConfig)
	testUsername := "test-user"
	testPassword := "test-pass"
	testExpirationTime := 999
	jwt.Users = map[string]string{testUsername: testPassword}
	token, _ := jwt.GenerateNewToken(models.Credentials{Username: testUsername, Password: testPassword}, testExpirationTime)
	FilesUploadingTest(assert, testApi, token)
}

func TestUploadAuth(t *testing.T) {
	assert := assert.New(t)
	testConfig, _ := models.NewConfigYML("testdata/config-1.yaml")
	testConfig.MediaAPIConfig.StorageRootPath = t.TempDir()
	testApi := NewMediaAPI(&testConfig.MediaAPIConfig)
	testUsername := "test-user"
	testPassword := "test-pass"
	testExpirationTime := 999
	jwt.Users = map[string]string{testUsername: testPassword}
	AuthTest(assert, http.MethodPost, "/upload", testApi.upload, models.Credentials{Username: testUsername, Password: testPassword}, testExpirationTime)
	MethodValidationTest(assert, "/upload", testApi.upload, http.MethodPost, http.MethodPost)
	MethodValidationTest(assert, "/upload", testApi.upload, http.MethodPost, http.MethodGet)
}

func TestDelete(t *testing.T) {
	assert := assert.New(t)
	testConfig, _ := models.NewConfigYML("testdata/config-1.yaml")
	testConfig.MediaAPIConfig.StorageRootPath = t.TempDir()
	testApi := NewMediaAPI(&testConfig.MediaAPIConfig)
	testUsername := "test-user"
	testPassword := "test-pass"
	testExpirationTime := 999
	jwt.Users = map[string]string{testUsername: testPassword}

	AuthTest(assert, http.MethodGet, "/delete", testApi.delete, models.Credentials{Username: testUsername, Password: testPassword}, testExpirationTime)
	MethodValidationTest(assert, "/delete", testApi.delete, http.MethodGet, http.MethodGet)
	MethodValidationTest(assert, "/delete", testApi.delete, http.MethodGet, http.MethodPost)
}

type DeleteTestCase struct {
	file    string
	isValid bool
}

func TestDeleteFilesCases(t *testing.T) {
	assert := assert.New(t)
	testConfig, _ := models.NewConfigYML("testdata/config-1.yaml")
	testConfig.MediaAPIConfig.StorageRootPath = t.TempDir()
	testApi := NewMediaAPI(&testConfig.MediaAPIConfig)
	testUsername := "test-user"
	testPassword := "test-pass"
	testExpirationTime := 999
	jwt.Users = map[string]string{testUsername: testPassword}
	token, _ := jwt.GenerateNewToken(models.Credentials{Username: testUsername, Password: testPassword}, testExpirationTime)

	testCases := []DeleteTestCase{
		{file: "not-existed.txt", isValid: false},
		{file: "testdata/test.txt", isValid: true},
	}

	for _, testCase := range testCases {
		if testCase.isValid {
			mockUploadingTest(assert, testApi, token, testCase.file)
		}
		handler := http.HandlerFunc(testApi.delete)
		rec := httptest.NewRecorder()

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/delete?name=%s", utils.EncodeFileName(path.Base(testCase.file), utils.GetFileExtension(testCase.file))), nil)
		req.Header.Add("Token", token)
		handler.ServeHTTP(rec, req)
		if testCase.isValid {
			assert.Equal(http.StatusOK, rec.Code)
		} else {
			assert.Equal(http.StatusBadRequest, rec.Code)
		}
	}
}

func TestGetAvailableExtension(t *testing.T) {
	assert := assert.New(t)
	testConfig, _ := models.NewConfigYML("testdata/config-1.yaml")
	testConfig.MediaAPIConfig.StorageRootPath = t.TempDir()
	testApi := NewMediaAPI(&testConfig.MediaAPIConfig)
	testUsername := "test-user"
	testPassword := "test-pass"
	testExpirationTime := 999
	jwt.Users = map[string]string{testUsername: testPassword}
	token, _ := jwt.GenerateNewToken(models.Credentials{Username: testUsername, Password: testPassword}, testExpirationTime)

	MethodValidationTest(assert, "/extensions", testApi.getAvailableExtension, http.MethodGet, http.MethodGet)
	MethodValidationTest(assert, "/extensions", testApi.getAvailableExtension, http.MethodGet, http.MethodPost)

	mockUploadingTest(assert, testApi, token, "testdata/test.txt")

	expected := "[{\"name\":\"txt\",\"isDir\":true}]"

	handler := http.HandlerFunc(testApi.getAvailableExtension)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/extensions", nil)
	handler.ServeHTTP(rec, req)
	assert.Equal(expected, string(rec.Body.Bytes()))
}

func TestGetAvailableExtensionNotFoundDir(t *testing.T) {
	assert := assert.New(t)
	testConfig, _ := models.NewConfigYML("testdata/config-1.yaml")
	testConfig.MediaAPIConfig.StorageRootPath = t.TempDir()
	testApi := NewMediaAPI(&testConfig.MediaAPIConfig)

	handler := http.HandlerFunc(testApi.getAvailableExtension)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/extensions", nil)
	handler.ServeHTTP(rec, req)
	assert.Equal(http.StatusNotFound, rec.Code)
}

func TestGetFileNamesWithDates(t *testing.T) {
	assert := assert.New(t)
	testConfig, _ := models.NewConfigYML("testdata/config-1.yaml")
	testConfig.MediaAPIConfig.StorageRootPath = t.TempDir()
	testApi := NewMediaAPI(&testConfig.MediaAPIConfig)
	testUsername := "test-user"
	testPassword := "test-pass"
	testFileName := "testdata/test.txt"
	testExpirationTime := 999
	jwt.Users = map[string]string{testUsername: testPassword}
	token, _ := jwt.GenerateNewToken(models.Credentials{Username: testUsername, Password: testPassword}, testExpirationTime)

	MethodValidationTest(assert, "/dir", testApi.getFileNamesWithDates, http.MethodGet, http.MethodGet)
	MethodValidationTest(assert, "/dir", testApi.getFileNamesWithDates, http.MethodGet, http.MethodPost)

	mockUploadingTest(assert, testApi, token, testFileName)

	handler := http.HandlerFunc(testApi.getFileNamesWithDates)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/dir?extension=%s", utils.GetFileExtension(testFileName)), nil)
	handler.ServeHTTP(rec, req)
	assert.Equal(true, strings.Contains(
		string(rec.Body.Bytes()),
		utils.EncodeFileName(path.Base(testFileName), utils.GetFileExtension(testFileName))))
}

func TestGetDataByUrl(t *testing.T) {
	assert := assert.New(t)
	testConfig, _ := models.NewConfigYML("testdata/config-1.yaml")
	testConfig.MediaAPIConfig.StorageRootPath = t.TempDir()
	testApi := NewMediaAPI(&testConfig.MediaAPIConfig)
	testUsername := "test-user"
	testPassword := "test-pass"
	testFileName := "testdata/test.txt"
	testExpirationTime := 999
	jwt.Users = map[string]string{testUsername: testPassword}
	token, _ := jwt.GenerateNewToken(models.Credentials{Username: testUsername, Password: testPassword}, testExpirationTime)

	fmt.Println(path.Base(path.Join(testApi.Prefix, "/data/", testFileName)))
	MethodValidationTest(assert, path.Join("/data/", testFileName), testApi.getDataByUrl, http.MethodGet, http.MethodGet)
	MethodValidationTest(assert, path.Join("/data/", testFileName), testApi.getDataByUrl, http.MethodGet, http.MethodPost)

	loadedFileName := testFileName
	handler := http.HandlerFunc(testApi.getDataByUrl)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, path.Join("/data/", loadedFileName), nil)
	handler.ServeHTTP(rec, req)
	assert.Equal(http.StatusNotFound, rec.Code)

	rec, _ = mockUploadingTest(assert, testApi, token, testFileName)
	loadedFileName = string(rec.Body.Bytes())
	handler = http.HandlerFunc(testApi.getDataByUrl)
	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, path.Join("/data/", loadedFileName), nil)
	handler.ServeHTTP(rec, req)
	assert.Positive(rec.Body.Len())
	assert.Equal(http.StatusOK, rec.Code)
}
