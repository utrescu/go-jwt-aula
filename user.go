package main

// User defineix les dades d'usuari
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s User) hasValues() bool {
	return s.Username != "" && s.Password != ""
}
