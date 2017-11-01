package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
)

// TokenData defineix les dades d'usuari
type TokenData struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// GetTokenHandler Genera un JWT per l'usuari rebut
// ------------------------------------------------------------------------
func GetTokenHandler(w http.ResponseWriter, req *http.Request) {
	var user User
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewDecoder(req.Body).Decode(&user); err == io.EOF {
		json.NewEncoder(w).Encode(Exception{Message: "Incorrect User"})
	} else if err != nil {
		json.NewEncoder(w).Encode(Exception{Message: "Incorrect User"})
	} else if !user.hasValues() {
		json.NewEncoder(w).Encode(Exception{Message: "Incorrect User"})
	} else {
		expireToken := time.Now().Add(time.Hour * 1).Unix()
		expireCookie := time.Now().Add(time.Hour * 1)
		claims := TokenData{user.Username,
			jwt.StandardClaims{
				ExpiresAt: expireToken,
				Issuer:    "localhost:3000",
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, error := token.SignedString(clauDeSignat)
		if error != nil {
			fmt.Println(error)
		}

		cookie := http.Cookie{Name: "Auth", Value: tokenString, Expires: expireCookie, HttpOnly: true}
		http.SetCookie(w, &cookie)

		json.NewEncoder(w).Encode(JwtToken{Token: tokenString})
	}
}

// ValidateToken és un middleware que comprova que el token és correcte
// ------------------------------------------------------------------------
func ValidateToken(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		authorizationHeader := req.Header.Get("authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				token, error := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("There was an error")
					}
					return clauDeSignat, nil
				})
				if error != nil {
					json.NewEncoder(w).Encode(Exception{Message: error.Error()})
					return
				}
				if token.Valid {
					context.Set(req, "decoded", token.Claims)
					next(w, req)
				} else {
					json.NewEncoder(w).Encode(Exception{Message: "Invalid authorization token"})
				}
			}
		} else {
			json.NewEncoder(w).Encode(Exception{Message: "An authorization header is required"})
		}
	})
}
