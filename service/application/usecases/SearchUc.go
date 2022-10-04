package usecases

import "lambda-search-nir/service/application/domain"

type SearchUc interface {
	LexicalSearch(query string) ([]domain.ScoreResult, error)
}
