package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

// User : Form or JSON user data
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u User) hasValues() bool {
	return u.Username != "" && u.Password != ""
}

func (u User) isValid(db *sql.DB) bool {

	if hashFromDatabase, err := recuperarDeBaseDeDades(db, u.Username); err != nil {
		// Comparing the password with the hash
		if err := bcrypt.CompareHashAndPassword(hashFromDatabase, []byte(u.Password)); err != nil {
			return false
		}
		return true
	}
	return false
}
