package main

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var SECRET = []byte("Do not hardcode stuff like this")

func CreateJWT(username, role string) (tokenString string, err error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	claims["username"] = username
	claims["role"] = role
	tokenString, err = token.SignedString(SECRET)
	return
}

func ValidateJWT(r *http.Request, method string) (username, role string) {
	return
}
