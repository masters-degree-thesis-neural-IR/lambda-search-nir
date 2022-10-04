package repositories

type IndexMemoryRepository interface {
	FindByTerm(term string) ([]string, error)
}
