package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var SECRET = []byte("super-secret-auth-key")


func createJwt() (string, error) {
	token := jwt.New(jwt.SigningMethodHS384)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	tokenStr, err := token.SignedString(SECRET)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	return tokenStr, nil
}

func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "seper secret api")
}

func main() {
	http.HandleFunc("/api", Home)
	http.ListenAndServe(":8080", nil)
}