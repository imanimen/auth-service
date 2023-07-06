package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var SECRET = []byte("super-secret-auth-key")

var api_key = "1234"

func validateJwt(next func(w http.ResponseWriter, r* http.Request)) http.Handler {
	return http.HandlerFunc( func(w http.ResponseWriter, r *http.Request) {
		if r.Header["token"] != nil {
			token, err := jwt.Parse(r.Header["token"][0], func(t *jwt.Token) (interface{}, error) {
				_, ok := t.Method.(*jwt.SigningMethodHMAC)
				if !ok {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("Unauthorized"))
				}
				return SECRET, nil
			})
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized " + err.Error() ))
			}
			if token.Valid {
				next(w, r)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized"))
			}
		}
	})
}

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



func getJwt(w http.ResponseWriter, r *http.Request) {
	if r.Header["api"] != nil {
		if r.Header["api"][0] == api_key {
			token, err := createJwt()
			if err != nil {
				return 
			}
			fmt.Fprint(w, token)
		}
	}
}

func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "seper secret api")
}

func main() {
	http.Handle("/api", validateJwt(Home))
	http.HandleFunc("/jwt", getJwt)
	http.ListenAndServe(":8080", nil)
}