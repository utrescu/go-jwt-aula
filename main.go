package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"html/template"
	"log"
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

// Base de dades
var db *sql.DB

// clauDeSignat és la clau que fem servir per signar el Token
var clauDeSignat = []byte("SiLaLletFosXocolataNoCaldriaColacao")

func main() {
	// Iniciar el router gorilla/mux
	router := mux.NewRouter()

	// Carregar la configuració
	err := config.loadConfig("config/aules.toml")
	if err != nil {
		panic("aules.toml " + err.Error())
	}

	// Obrir la base de dades
	db, err = initDB("config/usuaris.db")
	if err != nil {
		panic("Base de dades no trobada")
	}
	defer db.Close()

	// router.HandleFunc("/base", BaseHandler).Methods("GET")
	// Login no està protegit per JWT
	router.HandleFunc("/login", LoginHandler).Methods("POST")

	// Protegim les URL amb el middleware ValidateToken(*) de jwt.go
	router.HandleFunc("/login", ToLoginHandler).Methods("GET")
	router.HandleFunc("/help", HelpHandler).Methods("GET")
	router.HandleFunc("/aula/list", ValidateToken(ListAulesHandler)).Methods("GET")
	router.HandleFunc("/aula/{num}/status", ValidateToken(ListClasse)).Methods("GET")
	router.HandleFunc("/aula/{num}/stop", ValidateToken(NotImplemented)).Methods("POST")
	router.HandleFunc("/logout", ValidateToken(Logout)).Methods("GET")

	// A l'arrel simplement mostrem una pàgina estàtica
	// i posem els seus recursos a 'static' (en realitat no fa cap falta)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./views/")))
	http.Handle("/", router)

	// Port en el que escoltarà el servidor
	http.ListenAndServe(":3000", handlers.LoggingHandler(os.Stdout, router))

}

// ToLoginHandler mostra la pàgina de login a menys que ja tingui la cookie
// ------------------------------------------------------------------------
var ToLoginHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

	if correctCookie(req) {
		t, _ := template.ParseFiles("templates/base.html")
		t.Execute(w, config.Aules)
		// http.Redirect(w, req, "/base", http.StatusSeeOther)
	} else {
		log.Println("/login -> cookie no correcta")
		http.ServeFile(w, req, "./views/login.html")
	}
})

// HelpHandler mostra la pàgina d'error en format HTML
// ------------------------------------------------------------------------
var HelpHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "./views/help.html")
})

// LoginHandler intenta capturar el contingut rebut i generar un token
//
// - "application/json": Es pot processar tot sol
// - "form-urlencoded" : S'ha de convertir
// ------------------------------------------------------------------------
var LoginHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	// Sempre generem contingut JSON
	w.Header().Set("Content-Type", "application/json")

	var user User
	var err error

	switch contentType := req.Header.Get("Content-type"); contentType {
	case "application/json":
		err = json.NewDecoder(req.Body).Decode(&user)

	case "application/x-www-form-urlencoded":
		err = req.ParseForm()
		if err != nil {
			break
		}
		user.Username = req.FormValue("username")
		user.Password = req.FormValue("password")
	default:
		err = errors.New("Content-Type error")
		return
	}
	// Error or no user ...
	if err != nil || !user.hasValues() || !user.hasCorrectPassword(db) {
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
