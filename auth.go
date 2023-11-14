package main

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var SECRET = []byte("Do not hardcode stuff like this")

func CreateJWT(id int, role string) (tokenString string, err error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	claims["userID"] = id
	claims["role"] = role
	tokenString, err = token.SignedString(SECRET)
	return
}

// func ProfileInfo(server *APIServer, ) http.HandlerFunc {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 		if r.Header["Token"] != nil {
// 			token, err := jwt.Parse(r.Header["Token"][0], func(t *jwt.Token) (interface{}, error) {
// 				_, ok := t.Method.(*jwt.SigningMethodHMAC)
// 				if !ok {
// 					unauthorized(w)
// 				}
// 				return SECRET, nil
// 			})
// 			if err != nil || !token.Valid {
// 				unauthorized(w)
// 			} else if token.Valid {
// 				claims := token.Claims.(jwt.MapClaims)
// 				userID := claims["userID"].(int)
// 				role := claims["role"].(string)
// 				switch role {
// 				case "employer":
// 					server.
// 				}
// 			}
// 		} else {
// 			unauthorized(w)
// 		}
// 	})
// }

// func unauthorized(w http.ResponseWriter) {
// 	w.WriteHeader(http.StatusUnauthorized)
// 	w.Write([]byte("Access not authorized"))
// }
