package ports

import "lambda-search-nir/service/application/domain"

type DocumentEvent interface {
	Created(document domain.Document) error
}
