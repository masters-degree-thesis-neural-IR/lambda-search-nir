package repositories

import "lambda-search-nir/service/application/domain"

type DocumentRepository interface {
	FindById(id string) (*domain.Document, error)
}
