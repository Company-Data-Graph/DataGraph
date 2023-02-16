package jwt

import (
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"media-server/src/models"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var Users map[string]string

var Key []byte

func GenerateNewToken(creds models.Credentials, t int) (string, error) {
	if len(Key) == 0 {
		log.Println("Key not found! Generate new key!")
		Key = []byte(base64.RawStdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", rand.Intn(1000000)))))
		log.Println("New key:", string(Key))
	}
	expirationTime := time.Now().Add(time.Minute * time.Duration(t))
	claims := &models.Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(Key)
	return tokenString, err
}

func CheckToken() {

}

func CredentialsValidation(creds models.Credentials) bool {
	expectedCreds, ok := Users[creds.Username]
	if !ok || creds.Password != expectedCreds {
		return false
	}
	return true
}
