package models

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type YamlConfigTestCase struct {
	path    string
	isValid bool
}

func TestNewConfigYML(t *testing.T) {
	assert := assert.New(t)
	testCases := []YamlConfigTestCase{
		{path: "testdata/valid-config.yaml", isValid: true},
		{path: "testdata/invalid-config.yaml", isValid: false},
		{path: "not/existed.yaml", isValid: false},
	}
	for _, testCase := range testCases {
		_, err := NewConfigYML(testCase.path)
		if testCase.isValid {
			assert.NoError(err, "This config was valid")
		} else {
			assert.Error(err, "This config was invalid")
		}
	}
}

func TestNewConfigENV(t *testing.T) {
	assert := assert.New(t)
	testCases := []MediaAPIConfig{
		{Host: "test-host", Port: 1000, Prefix: "test-prefix", AdminPass: "test-pass", TokenLiveTime: 0, StorageRootPath: "test-root", DataStorageRoute: "test-route"},
		{Port: 100, Prefix: "test-prefix1", AdminPass: "test-pass-1", TokenLiveTime: 1, StorageRootPath: "test-root-1", DataStorageRoute: "test-route1"},
		{Host: "test-host-2", Prefix: "test-prefix2", AdminPass: "test-pass-2", TokenLiveTime: 2, StorageRootPath: "test-root-2", DataStorageRoute: "test-route2"},
		{Host: "test-host-3", Port: 300, AdminPass: "test-pass-3", TokenLiveTime: 3, StorageRootPath: "test-root-3", DataStorageRoute: "test-route3"},
		{Host: "test-host-4", Port: 400, Prefix: "test-prefix4", TokenLiveTime: 4, StorageRootPath: "test-root-4", DataStorageRoute: "test-route4"},
		{Host: "test-host-5", Port: 500, Prefix: "test-prefix5", AdminPass: "test-pass-5", StorageRootPath: "test-root-5", DataStorageRoute: "test-route5"},
		{Host: "test-host-6", Port: 100, Prefix: "test-prefix6", AdminPass: "test-pass-6", TokenLiveTime: 6, StorageRootPath: "test-root-6"},
	}
	for _, testCase := range testCases {
		t.Setenv("MEDIA_SERVER_HOST", testCase.Host)
		t.Setenv("MEDIA_SERVER_PORT", fmt.Sprintf("%d", testCase.Port))
		t.Setenv("MEDIA_SERVER_ADMIN_PASS", testCase.AdminPass)
		t.Setenv("MEDIA_SERVER_TOKEN_LIVE_TIME", fmt.Sprintf("%d", testCase.TokenLiveTime))
		t.Setenv("MEDIA_SERVER_STORAGE_ROOT_PATH", testCase.StorageRootPath)
		t.Setenv("MEDIA_SERVER_DATA_ROUTE_STORAGE_ROUTE", testCase.DataStorageRoute)
		config, err := NewConfigENV()
		assert.Equal(testCase.Host, config.MediaAPIConfig.Host)
		assert.Equal(testCase.Port, config.MediaAPIConfig.Port)
		assert.Equal(testCase.AdminPass, config.MediaAPIConfig.AdminPass)
		assert.Equal(testCase.TokenLiveTime, config.MediaAPIConfig.TokenLiveTime)
		assert.Equal(testCase.StorageRootPath, config.MediaAPIConfig.StorageRootPath)
		assert.Equal(testCase.DataStorageRoute, config.MediaAPIConfig.DataStorageRoute)
		assert.NoError(err, "Error must be nil because all params equal")
	}
}

func TestNewConfigENVInvalidPort(t *testing.T) {
	assert := assert.New(t)
	t.Setenv("MEDIA_SERVER_PORT", "not-a-number")
	config, err := NewConfigENV()
	assert.Error(err, "Port env was not a number")
	assert.Nil(config, "Config parsing can be finished with invalid port value")
}

func TestNewConfigENVInvalidTokenLiveTime(t *testing.T) {
	assert := assert.New(t)
	t.Setenv("MEDIA_SERVER_TOKEN_LIVE_TIME", "not-a-number")
	config, err := NewConfigENV()
	assert.Error(err, "Token live time env was not a number")
	assert.Nil(config, "Config parsing can be finished with invalid token live time value")
}
