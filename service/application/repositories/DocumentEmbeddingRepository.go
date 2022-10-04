package repositories

import "lambda-search-nir/service/application/domain"

type DocumentEmbeddingRepository interface {
	FindByDocumentIDs(documentIDs map[string]int8) ([]domain.DocumentEmbedding, error)
}
