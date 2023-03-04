package jwt

import (
	"github.com/stretchr/testify/assert"
	"media-server/internal/models"
	"testing"
	"time"
)

func TestCredentialsValidation(t *testing.T) {
	assert := assert.New(t)
	testCases := []models.Credentials{
		{Username: "test-user-1", Password: "test-pass-1"},
		{Username: "test-user-2", Password: "test-pass-2"},
		{Username: "", Password: "test-pass-3"},
		{Username: "test-user-4", Password: ""},
	}
	Users = make(map[string]string)
	for _, testCase := range testCases {
		Users[testCase.Username] = testCase.Password
	}
	for _, testCase := range testCases {
		assert.Equal(true, CredentialsValidation(testCase))
	}
}

func TestWrongCredentialsValidation(t *testing.T) {
	assert := assert.New(t)
	Users = make(map[string]string)
	testCreds := models.Credentials{Username: "user", Password: "pass"}
	assert.Equal(false, CredentialsValidation(testCreds))
	testCreds = models.Credentials{Username: "user"}
	assert.Equal(false, CredentialsValidation(testCreds))
}

func TestGenerateClaims(t *testing.T) {
	assert := assert.New(t)
	testUsername := "test-user"
	testPassword := "test-pass"
	testExpirationTime := 999

	claims := generateClaims(models.Credentials{
		Username: testUsername,
		Password: testPassword,
	}, testExpirationTime)
	assert.Equal(testUsername, claims.Username)
	assert.Equal(time.Now().Add(time.Minute*time.Duration(testExpirationTime)).Unix(), claims.ExpiresAt)
}

func TestGetKey(t *testing.T) {
	assert := assert.New(t)
	current := getKey()
	assert.Positive(len(Key), "Key not generated")
	assert.Equal(current, getKey(), "Key regenerated after first generation")
}

type GenerateNewTokenTestCase struct {
	creds   models.Credentials
	isValid bool
}

func TestGenerateNewToken(t *testing.T) {
	assert := assert.New(t)
	testCases := []GenerateNewTokenTestCase{
		{creds: models.Credentials{Username: "test-user-1", Password: "test-pass-1"}, isValid: true},
		{creds: models.Credentials{Username: "test-user-2", Password: "test-pass-2"}, isValid: true},
		{creds: models.Credentials{Password: "test-pass-1"}, isValid: true},
		{creds: models.Credentials{Username: "test-user-4"}, isValid: true},

		{creds: models.Credentials{Username: "no-reg-user-1", Password: "test-pass-not-reg"}, isValid: false},
		{creds: models.Credentials{Username: "no-reg-user-with-exist-pass", Password: "test-pass-1"}, isValid: false},
	}
	testExpirationTime := 999
	Users = make(map[string]string)
	for _, testCase := range testCases {
		if testCase.isValid {
			Users[testCase.creds.Username] = testCase.creds.Password
		}
	}
	current := ""
	for _, testCase := range testCases {
		if testCase.isValid {
			received, err := GenerateNewToken(testCase.creds, testExpirationTime)
			assert.NoError(err)
			assert.NotEqual(received, current, "Token repeat")
			current = received
		} else {
			_, err := GenerateNewToken(testCase.creds, testExpirationTime)
			assert.Error(err)
		}

	}
}
