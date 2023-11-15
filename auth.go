package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var SECRET = []byte("Do not hardcode stuff like this")

type Claims struct {
	Username string `json:"username"`
}

func CreateJWT(username, role string) (tokenString string, err error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	claims["username"] = username
	claims["role"] = role
	tokenString, err = token.SignedString(SECRET)
	return
}

func ValidateJWT(r *http.Request, method string) (username, role, token string) {
	t := r.FormValue("token")
	if t == "" {
		return
	}
	tok, err := jwt.Parse(t, func(tt *jwt.Token) (interface{}, error) {
		_, ok := tt.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("token error")
		}
		return SECRET, nil
	})
	if err == nil {
		claims, ok := tok.Claims.(jwt.MapClaims)
		if ok && tok.Valid {
			username = claims["username"].(string)
			role = claims["role"].(string)
			token = t
		}
	}
	return
}
