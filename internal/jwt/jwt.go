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

// Generate token is credentials valid
func GenerateNewToken(creds models.Credentials, t int) (string, error) {
	if !CredentialsValidation(creds) {
		return "", fmt.Errorf("Invalid credits")
	}
	claims := generateClaims(creds, t)
	token := jwtGo.NewWithClaims(jwtGo.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(Key)
	return tokenString, err
}

// Check user auth
func Auth(token string) error {
	var claims models.Claims
	tokenValidation, err := jwtGo.ParseWithClaims(token, &claims, func(t *jwtGo.Token) (interface{}, error) { return Key, nil })
	if err != nil {
		return err
	}
	if !tokenValidation.Valid {
		return err
	}
	return nil

}

// Generate new Key for jwt-token
func getKey() []byte {
	if len(Key) == 0 {
		Key = []byte(base64.RawStdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", rand.Intn(1000000)))))
	}
	return Key
}

// Generate claims for jwt-token
func generateClaims(creds models.Credentials, t int) *models.Claims {
	expirationTime := time.Now().Add(time.Minute * time.Duration(t))
	return &models.Claims{
		Username: creds.Username,
		StandardClaims: jwtGo.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
}

// Check jwt credentials valid
func CredentialsValidation(creds models.Credentials) bool {
	expectedCreds, ok := Users[creds.Username]
	if !ok || creds.Password != expectedCreds {
		return false
	}
	return true
}
