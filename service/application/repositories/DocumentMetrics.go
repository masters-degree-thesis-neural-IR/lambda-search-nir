package repositories

import "lambda-search-nir/service/application/domain"

type DocumentMetricsRepository interface {
	FindByDocumentIDs(documentIDs map[string]int8) ([]domain.NormalizedDocument, error)
}
