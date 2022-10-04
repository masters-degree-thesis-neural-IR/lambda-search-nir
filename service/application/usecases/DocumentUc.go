package usecases

import "lambda-search-nir/service/application/domain"

type DocumentUc interface {
	LoadDocuments(scoreResult []domain.ScoreResult) ([]domain.DocumentResult, error)
}
