package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

// Llista amb les aules del sistema
var config aules

// Base de dades
var db *sql.DB

// clauDeSignat és la clau que fem servir per signar el Token
var clauDeSignat = []byte("SiLaLletFosXocolataNoCaldriaColacao")

func main() {

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

	// Crear i Iniciar el router gorilla/mux
	servidor := Rutes{}
	servidor.Run(":3000")
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"message": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
