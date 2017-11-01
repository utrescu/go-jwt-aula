package main

import (
	"encoding/json"
	"net/http"
	"os"

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

// clauDeSignat és la clau que fem servir per signar el Token
var clauDeSignat = []byte("SiLaLletFosXocolataNoCaldriaColacao")

func main() {
	// Iniciar el router gorilla/mux
	router := mux.NewRouter()

	// A l'arrel simplement mostrem una pàgina estàtica i posem els seus recursos a 'static' (en realitat no fa cap falta)
	router.Handle("/", http.FileServer(http.Dir("./views/")))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	// Login no està protegit per JWT
	router.HandleFunc("/login", GetTokenHandler).Methods("POST")

	// Protegim les crides al Handler amb el middleware ValidateToken(*) de jwt.go
	router.HandleFunc("/aula/list", ValidateToken(ListAulesHandler)).Methods("GET")
	router.HandleFunc("/aula/{num}/status", ValidateToken(ListClasse)).Methods("GET")
	router.HandleFunc("/aula/{num}/stop", ValidateToken(NotImplemented)).Methods("POST")

	// Port en el que correrà el servidor
	http.ListenAndServe(":3000", handlers.LoggingHandler(os.Stdout, router))

}

// NotImplemented es fa servir quan algun dels recursos no està definit
// en teoria s'eliminarà en producció
// ------------------------------------------------------------------------
var NotImplemented = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("Not Implemented"))
	json.NewEncoder(w).Encode(Exception{Message: "Not implemented"})
})

// ListAulesHandler retorna la llista d'aules
// ------------------------------------------------------------------------
var ListAulesHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	decoded := context.Get(req, "decoded")
	var token TokenData
	mapstructure.Decode(decoded.(jwt.MapClaims), &token)

	payload, _ := json.Marshal(aules)

	w.Write([]byte(payload))
})

// ListClasse retorna les característiques de la classe
// ------------------------------------------------------------------------
var ListClasse = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

	var aula Aula
	vars := mux.Vars(req)
	num := vars["num"]
	aula.Nom = num

	resposta, _ := json.Marshal(pcEnMarxa[num])

	w.Write([]byte(resposta))

})
