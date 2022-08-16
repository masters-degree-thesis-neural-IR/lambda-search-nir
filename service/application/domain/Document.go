package domain

type Document struct {
	Id     string
	Length int
	Tf     map[string]int
}
