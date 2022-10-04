package service

import (
	"lambda-search-nir/service/application/domain"
	"lambda-search-nir/service/application/logger"
	"lambda-search-nir/service/application/repositories"
	"lambda-search-nir/service/application/usecases"
)

type DocumentService struct {
	Logger             logger.Logger
	DocumentRepository repositories.DocumentRepository
}

func NewDocumentService(logger logger.Logger, documentRepository repositories.DocumentRepository) usecases.DocumentUc {

	return DocumentService{
		Logger:             logger,
		DocumentRepository: documentRepository,
	}
}

func (s DocumentService) LoadDocuments(scoreResult []domain.ScoreResult) ([]domain.DocumentResult, error) {

	s.Logger.Info("Start LoadDocuments")

	var documentIDs = make([]string, len(scoreResult))
	for i, result := range scoreResult {
		documentIDs[i] = result.DocumentID
	}

	documents, err := s.DocumentRepository.FindByDocumentIDs(documentIDs)

	s.Logger.Info("Total documentos: ", len(documents))

	if err != nil {
		return []domain.DocumentResult{}, err
	}

	documentResults := make([]domain.DocumentResult, len(scoreResult))

	for i, result := range scoreResult {

		document := documents[result.DocumentID]

		documentResult := domain.DocumentResult{
			Similarity: result.Similarity,
			Document:   document,
		}
		documentResults[i] = documentResult

	}

	return documentResults, nil
}
