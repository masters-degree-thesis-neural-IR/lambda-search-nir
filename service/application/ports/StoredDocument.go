package ports

import "lambda-search-nir/service/application/domain"

type Store interface {
	StoreDocument(document domain.Document) error
}
