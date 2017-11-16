package main

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
)

func initDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil || db == nil {
		return nil, errors.New("DB error")
	}
	return db, nil
}

func recuperarDeBaseDeDades(db *sql.DB, user string) ([]byte, error) {

	sqlSearchUser := `
		SELECT password
		FROM users WHERE
		username = ?`

	rows, err := db.Query(sqlSearchUser)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []byte
	for rows.Next() {
		err2 := rows.Scan(&result)
		if err2 != nil {
			return nil, err2
		}
	}
	return result, nil

}
