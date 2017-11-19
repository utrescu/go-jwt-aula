package main

import (
	"os"

	"github.com/naoina/toml"
	"github.com/utrescu/listIP"
)

// AulaResult és per treure resultats sobre una aula
type AulaResult struct {
	Aula    string
	EnMarxa []string
}

type aules struct {
	Aules map[string]AulaInfo
}

// AulesResult retorna les aules
type AulesResult struct {
	Aules []string `json:"aules"`
}

type AulaInfo struct {
	Rang string
	Name string
	Port int
}

func (a *aules) loadConfig(fitxer string) error {
	f, err := os.Open(fitxer)
	if err != nil {
		return err
	}
	defer f.Close()
	return toml.NewDecoder(f).Decode(a)
}

func (a *aules) listAules() AulesResult {
	var resultat AulesResult

	keys := make([]string, len(a.Aules))

	i := 0
	for k := range a.Aules {
		keys[i] = k
		i++
	}
	resultat.Aules = keys
	return resultat
}

func (a *AulaInfo) cercaMaquines(numAula string) (AulaResult, error) {
	var resultat AulaResult
	resultat.Aula = numAula
	enmarxa, _, err := listIP.Check([]string{a.Rang}, a.Port, 64, "100ms")
	resultat.EnMarxa = enmarxa
	return resultat, err
}
