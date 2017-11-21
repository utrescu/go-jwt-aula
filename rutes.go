package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
)

// Rutes : defineix les rutes del servei
type Rutes struct {
	Router *mux.Router
}

// Run : Executa el programa
func (a *Rutes) Run(addr string) {
	a.Router = mux.NewRouter()
	a.inicialitzaRoutes()
	// log.Fatal(http.ListenAndServe(":8000", a.Router))

	http.ListenAndServe(addr, handlers.LoggingHandler(os.Stdout, handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(a.Router)))
}

func (a *Rutes) inicialitzaRoutes() {

	a.Router.HandleFunc("/login", a.entrar).Methods("POST")
	a.Router.HandleFunc("/help", a.mostraAjuda).Methods("GET")

	// Protegim les URL amb el middleware ValidateToken(*) de jwt.go
	a.Router.HandleFunc("/aula/list", ValidateToken(a.llistaAules)).Methods("GET")
	a.Router.HandleFunc("/aula/{num}/status", ValidateToken(a.llistaClasse)).Methods("GET")
	a.Router.HandleFunc("/aula/{num}/stop", ValidateToken(a.noImplementat)).Methods("POST")
	a.Router.HandleFunc("/logout", ValidateToken(a.sortir)).Methods("GET")

	// A l'arrel simplement mostrem una pàgina estàtica
	// i posem els seus recursos a 'static' (en realitat no fa cap falta)
	a.Router.PathPrefix("/").Handler(http.FileServer(http.Dir("./views/")))
	http.Handle("/", a.Router)

}

// // ToLoginHandler mostra la pàgina de login a menys que ja tingui la cookie
// // ------------------------------------------------------------------------
// var ToLoginHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

// 	if correctCookie(req) {
// 		t, _ := template.ParseFiles("templates/base.html")
// 		t.Execute(w, config.Aules)
// 		// http.Redirect(w, req, "/base", http.StatusSeeOther)
// 	} else {
// 		http.ServeFile(w, req, "./views/login.html")
// 	}
// })

// HelpHandler mostra la pàgina d'error en format HTML
// ------------------------------------------------------------------------
func (a *Rutes) mostraAjuda(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "./views/help.html")
}

// LoginHandler intenta capturar el contingut rebut i generar un token
//
// - "application/json": Es pot processar tot sol
// - "form-urlencoded" : S'ha de convertir
// ------------------------------------------------------------------------
func (a *Rutes) entrar(w http.ResponseWriter, req *http.Request) {
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
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(Exception{Message: "Incorrect User"})
		return
	}

	tokenString, err := GetTokenHandler(user)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(Exception{Message: "Error generating token"})
		return
	}

	// Cookie
	// expireCookie := time.Now().Add(time.Hour * 1)
	// cookie := http.Cookie{Name: "Auth", Value: tokenString, Expires: expireCookie, HttpOnly: true}
	// http.SetCookie(w, &cookie)

	// Generar el token i la resposta
	json.NewEncoder(w).Encode(JwtToken{Token: tokenString})
}

// Logout es fa servir per desconnectar els clients web
// ------------------------------------------------------------------------
func (a *Rutes) sortir(w http.ResponseWriter, req *http.Request) {
	// deleteCookie := http.Cookie{Name: "Auth", Value: "none", Expires: time.Now()}
	// http.SetCookie(res, &deleteCookie)

	// Expire Token?
	return
}

// noImplementat es fa servir quan algun dels recursos no està definit
// en teoria s'eliminarà en producció
// ------------------------------------------------------------------------
func (a *Rutes) noImplementat(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	json.NewEncoder(w).Encode(Exception{Message: "Not implemented"})
}

// llistaAules retorna la llista d'aules
// ------------------------------------------------------------------------
func (a *Rutes) llistaAules(w http.ResponseWriter, req *http.Request) {
	decoded := context.Get(req, "decoded")
	var token TokenData
	mapstructure.Decode(decoded.(jwt.MapClaims), &token)

	payload, _ := json.Marshal(config.listAules())
	w.Write([]byte(payload))
}

// llistaClasse retorna les característiques d'una classe determinada que
// es rep com a paràmetre
// ------------------------------------------------------------------------
func (a *Rutes) llistaClasse(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	numAula := vars["num"]

	// Si l'aula existeix mira de localitzar les màquines en marxa
	if infoAula, ok := config.Aules[numAula]; ok {

		aula, err := infoAula.cercaMaquines(numAula)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(Exception{Message: err.Error()})
			return
		}
		resposta, _ := json.Marshal(aula)
		w.Write([]byte(resposta))
	} else {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(Exception{Message: "Inexistent class"})
	}
}
