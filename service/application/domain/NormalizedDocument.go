package domain

type NormalizedDocument struct {
	Id     string
	Length int
	Tf     map[string]int
}
