package utils

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type CheckMethodTestCase struct {
	first          string
	second         string
	expectedResult bool
}

func TestCheckMethod(t *testing.T) {
	assert := assert.New(t)
	testCases := []CheckMethodTestCase{
		{http.MethodGet, http.MethodGet, true},
		{http.MethodGet, http.MethodPost, false},
		{http.MethodPost, http.MethodPost, true},
		{http.MethodPost, http.MethodOptions, false},
		{http.MethodDelete, http.MethodPut, false},
	}
	for _, testCase := range testCases {
		assert.Equal(testCase.expectedResult, CheckMethod(testCase.first, testCase.second))
	}
}
