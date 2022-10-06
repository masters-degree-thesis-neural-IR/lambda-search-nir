package repositories

import "lambda-search-nir/service/application/domain"

type IndexMemoryRepository interface {
	FindByTerm(term string) ([]string, error)
	Save(term string, document domain.NormalizedDocument) error
	LoadMetricsFromCache(documentIDs map[string]int8) ([]domain.NormalizedDocument, map[string]int8, error)
}
