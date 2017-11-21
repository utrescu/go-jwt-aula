package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
)

// JwtToken es fa servir per retornar el token web
type JwtToken struct {
	Token string `json:"token"`
}

// TokenData defineix les dades d'usuari
type TokenData struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// GetTokenHandler Genera un JWT per l'usuari rebut
// ------------------------------------------------------------------------
func GetTokenHandler(user User) (string, error) {

	expireToken := time.Now().Add(time.Hour * 1).Unix()

	claims := TokenData{user.Username,
		jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer:    "localhost:3000",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(clauDeSignat)
}

// ValidateToken és un middleware que comprova que el token és correcte
//
// Modificat per poder fer servir Cookies a més de que s'envïi l'autenticació
// a més de la versió original en el Header...
//
//  Recomanació de: https://stormpath.com/blog/build-secure-user-interfaces-using-jwts
// ------------------------------------------------------------------------
func ValidateToken(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		var tokenRebut string

		w.Header().Set("Content-Type", "application/json")

		// Comprovar si hi ha una Cookie 'Auth'
		// cookie, err := req.Cookie("Auth")
		//if err != nil {
		// Si no hi ha Cookie, mirem les capsaleres
		authorizationHeader := req.Header.Get("authorization")
		if authorizationHeader == "" {
			respondWithError(w, http.StatusUnauthorized, "An authorization token is required")
			return
		}
		bearerToken := strings.Split(authorizationHeader, " ")
		if len(bearerToken) != 2 {
			respondWithError(w, http.StatusUnauthorized, "An authorization token is required")
			return
		}
		tokenRebut = bearerToken[1]

		// } else {
		//	tokenRebut = cookie.Value
		// }

		// token, err := jwt.ParseWithClaims(tokenRebut, &TokenData{}, func(token *jwt.Token) (interface{}, error) {
		token, err2 := jwt.Parse(tokenRebut, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error")
			}
			return clauDeSignat, nil
		})
		if err2 != nil {
			respondWithError(w, http.StatusBadRequest, err2.Error())
			return
		}
		if token.Valid {
			context.Set(req, "decoded", token.Claims)
			next(w, req)
		} else {
			respondWithError(w, http.StatusUnauthorized, "Invalid authorization token")
		}

	})
}
