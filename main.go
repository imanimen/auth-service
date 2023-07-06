package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var SECRET = []byte("super-secret-auth-key")

var api_key = "1234"

func validateJwt(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "" {
			tokenString := r.Header.Get("Authorization")[7:]
			token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
				}
				return []byte(SECRET), nil
			})
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized " + err.Error()))
				return
			}
			if token.Valid {
				next(w, r)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized"))
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
		}
	}
}

func createJwt() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	tokenStr, err := token.SignedString(SECRET)
	if err != nil {
		errMsg := fmt.Sprintf("Error creating JWT: %s", err.Error())
		fmt.Println(errMsg)
		return "", fmt.Errorf(errMsg)
	}
	return tokenStr, nil
}

func getJwt(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") == api_key {
		token, err := createJwt()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			errMsg := fmt.Sprintf("Error creating JWT: %s", err.Error())
			w.Write([]byte(errMsg))
			return
		}
		fmt.Fprint(w, token)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
	}
}

func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "super secret api")
}

func main() {
	http.Handle("/api", validateJwt(Home))
	http.HandleFunc("/jwt", getJwt)
	http.ListenAndServe(":8080", nil)
}