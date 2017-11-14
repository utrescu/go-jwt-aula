package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
)

// JwtToken es fa servir per retornar el token web
type JwtToken struct {
	Token string `json:"token"`
}

// Exception es fa servir per retornar missatges
type Exception struct {
	Message string `json:"message"`
}

// Llista amb les aules del sistema
var config aules

// clauDeSignat és la clau que fem servir per signar el Token
var clauDeSignat = []byte("SiLaLletFosXocolataNoCaldriaColacao")

func main() {
	// Iniciar el router gorilla/mux
	router := mux.NewRouter()

	// Carregar la configuració
	err := config.loadConfig("aules.toml")
	if err != nil {
		panic("No s'ha pogut carregar la configuració: aula.toml")
	}

	// A l'arrel simplement mostrem una pàgina estàtica
	// i posem els seus recursos a 'static' (en realitat no fa cap falta)
	router.Handle("/", http.FileServer(http.Dir("./views/")))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	// Login no està protegit per JWT
	router.HandleFunc("/login", LoginHandler).Methods("POST")

	// Protegim les URL amb el middleware ValidateToken(*) de jwt.go
	router.HandleFunc("/aula/list", ValidateToken(ListAulesHandler)).Methods("GET")
	router.HandleFunc("/aula/{num}/status", ValidateToken(ListClasse)).Methods("GET")
	router.HandleFunc("/aula/{num}/stop", ValidateToken(NotImplemented)).Methods("POST")
	router.HandleFunc("/logout", ValidateToken(Logout)).Methods("GET")

	// Port en el que escoltarà el servidor
	http.ListenAndServe(":3000", handlers.LoggingHandler(os.Stdout, router))

}

// LoginHandler intenta capturar el contingut rebut i generar un token
//
// - "application/json": Es pot processar tot sol
// - "form-urlencoded" : S'ha de convertir
// ------------------------------------------------------------------------
var LoginHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	// Sempre generem contingut JSON
	w.Header().Set("Content-Type", "application/json")

	var user User

	switch contentType := req.Header.Get("Content-type"); contentType {
	case "application/json":
		if err := json.NewDecoder(req.Body).Decode(&user); err == io.EOF {
			json.NewEncoder(w).Encode(Exception{Message: "Incorrect User"})
			return
		} else if err != nil {
			json.NewEncoder(w).Encode(Exception{Message: "Incorrect User"})
			return
		}
	case "application/x-www-form-urlencoded":
		err := req.ParseForm()
		if err != nil {
			json.NewEncoder(w).Encode(Exception{Message: "Form data incorrect"})
			return
		}
		user.Username = req.FormValue("username")
		user.Password = req.FormValue("password")
	default:
		json.NewEncoder(w).Encode(Exception{Message: "Content-Type " + contentType + "not implemented"})
		return
	}

	// Comprovar l'usuari
	if !user.hasValues() {
		json.NewEncoder(w).Encode(Exception{Message: "Incorrect User"})
		return
	}

	tokenString, err := GetTokenHandler(user)
	if err != nil {
		json.NewEncoder(w).Encode(Exception{Message: "Error generating token"})
		return
	}

	// Generar el token i la resposta
	expireCookie := time.Now().Add(time.Hour * 1)
	cookie := http.Cookie{Name: "Auth", Value: tokenString, Expires: expireCookie, HttpOnly: true}
	http.SetCookie(w, &cookie)
	json.NewEncoder(w).Encode(JwtToken{Token: tokenString})
})

// Logout es fa servir per desconnectar els clients web
// ------------------------------------------------------------------------
func Logout(res http.ResponseWriter, req *http.Request) {
	deleteCookie := http.Cookie{Name: "Auth", Value: "none", Expires: time.Now()}
	http.SetCookie(res, &deleteCookie)
	return
}

// NotImplemented es fa servir quan algun dels recursos no està definit
// en teoria s'eliminarà en producció
// ------------------------------------------------------------------------
var NotImplemented = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Exception{Message: "Not implemented"})
})

// ListAulesHandler retorna la llista d'aules
// ------------------------------------------------------------------------
var ListAulesHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	decoded := context.Get(req, "decoded")
	var token TokenData
	mapstructure.Decode(decoded.(jwt.MapClaims), &token)

	payload, _ := json.Marshal(config.listAules())
	w.Write([]byte(payload))
})

// ListClasse retorna les característiques d'una classe determinada que
// es rep com a paràmetre
// ------------------------------------------------------------------------
var ListClasse = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	numAula := vars["num"]

	// Si l'aula existeix mira de localitzar les màquines en marxa
	if infoAula, ok := config.Aules[numAula]; ok {

		aula, err := infoAula.cercaMaquines(numAula)
		if err != nil {
			json.NewEncoder(w).Encode(Exception{Message: err.Error()})
			return
		}
		resposta, _ := json.Marshal(aula)
		w.Write([]byte(resposta))
	} else {
		json.NewEncoder(w).Encode(Exception{Message: "Inexistent class"})
	}
})
