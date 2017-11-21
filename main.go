package main

import (
	"database/sql"
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

	// Iniciar el router gorilla/mux
	servidor := Rutes{}

	servidor.Run(":3000")
}
