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

func (a *aules) listAules() []string {
	keys := make([]string, len(a.Aules))

	i := 0
	for k := range a.Aules {
		keys[i] = k
		i++
	}
	return keys
}

func (a *AulaInfo) cercaMaquines(numAula string) (AulaResult, error) {
	var resultat AulaResult
	resultat.Aula = numAula
	enmarxa, _, err := listIP.Check([]string{a.Rang}, a.Port, 64, "100ms")
	resultat.EnMarxa = enmarxa
	return resultat, err
}
