package jwt

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"media-server/internal/models"

	"time"

	jwtGo "github.com/dgrijalva/jwt-go"
)

var Users map[string]string
var Key []byte

func GenerateNewToken(creds models.Credentials, t int) (string, error) {
	if !CredentialsValidation(creds) {
		return "", fmt.Errorf("Invalid credits")
	}
	if len(Key) == 0 {
		Key = generateNewKey()
	}
	claims := generateClaims(creds, t)
	token := jwtGo.NewWithClaims(jwtGo.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(Key)
	return tokenString, err
}

func generateNewKey() []byte {
	return []byte(base64.RawStdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", rand.Intn(1000000)))))
}

func generateClaims(creds models.Credentials, t int) *models.Claims {
	expirationTime := time.Now().Add(time.Minute * time.Duration(t))
	return &models.Claims{
		Username: creds.Username,
		StandardClaims: jwtGo.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
}

func CredentialsValidation(creds models.Credentials) bool {
	expectedCreds, ok := Users[creds.Username]
	if !ok || creds.Password != expectedCreds {
		return false
	}
	return true
}
