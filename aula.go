package main

// Aula defineix les dades bàsiques d'una aula
type Aula struct {
	Nom    string
	Numero int
	Xarxa  string
}

// AulaDetall és per treure resultats sobre les aules individualment
type AulaDetall struct {
	Aula    string
	EnMarxa []string
}

//
// Dades falses per fer proves (a la realitat s'hauran de
// recuperar de la xarxa)
// ------------------------------------------------------------
//  - Aules disponibles
var aules = []Aula{
	Aula{"309", 309, "192.168.9.0/24"},
	Aula{"310", 310, "192.168.10.0/24"},
	Aula{"314", 314, "192.168.16.0/24"},
}

//   - Llista de PC de la xarxa
var pcEnMarxa = map[string]AulaDetall{
	"309": AulaDetall{"309", []string{"i309-01m", "i309-01d", "i309-03e"}},
	"310": AulaDetall{"310", []string{"i310-01e", "i310-02e", "i310-05e", "i310-05m", "i310-05d"}},
	"314": AulaDetall{"314", []string{"i314-01m", "i314-02e", "i314-03m", "i314-05e", "i314-05m"}},
}
