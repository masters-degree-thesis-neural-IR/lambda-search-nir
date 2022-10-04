package repositories

import "lambda-search-nir/service/application/domain"

type DocumentRepository interface {
	FindByDocumentIDs(documentIDs []string) (map[string]domain.Document, error)
}
