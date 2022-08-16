package usecases

import "lambda-search-nir/service/application/domain"

type SearchUc interface {
	SearchDocument(query string) ([]domain.QueryResult, error)
	//MakeInvertedIndex(localQuery []string) (domain.InvertedIndex, error)
}
