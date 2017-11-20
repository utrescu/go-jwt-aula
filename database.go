package main

import (
	"database/sql"
	"errors"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func initDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil || db == nil {
		return nil, errors.New("DB error")
	}
	return db, nil
}

func recuperaPasswordDeBaseDeDades(db *sql.DB, user string) ([]byte, error) {

	row := db.QueryRow("SELECT password FROM users WHERE username=?", user)
	// defer row.Close()

	var result []byte
	err := row.Scan(&result)
	if err != nil {
		log.Printf("DB: %s", err.Error())
		return nil, err
	}
	return result, nil
}
