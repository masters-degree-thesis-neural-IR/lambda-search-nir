package service

import (
	"lambda-search-nir/service/application/domain"
	"lambda-search-nir/service/application/exception"
	"lambda-search-nir/service/application/logger"
	"lambda-search-nir/service/application/nlp"
	"lambda-search-nir/service/application/repositories"
	"lambda-search-nir/service/application/usecases"
	"time"
)

type Search struct {
	Logger                    logger.Logger
	DocumentRepository        repositories.DocumentRepository
	IndexMemoryRepository     repositories.IndexMemoryRepository
	DocumentMetricsRepository repositories.DocumentMetricsRepository
}

func NewSearchService(logger logger.Logger, documentMetricsRepository repositories.DocumentMetricsRepository, indexMemoryRepository repositories.IndexMemoryRepository, documentRepository repositories.DocumentRepository) usecases.SearchUc {
	return Search{
		Logger:                    logger,
		DocumentRepository:        documentRepository,
		IndexMemoryRepository:     indexMemoryRepository,
		DocumentMetricsRepository: documentMetricsRepository,
	}
}

func (s Search) MakeInvertedIndex(localQuery []string, foundDocuments map[string]int8) (domain.InvertedIndex, error) {

	s.Logger.Info("MakeInvertedIndex ", len(foundDocuments))
	normalizedDocuments, err := s.DocumentMetricsRepository.FindByDocumentIDs(foundDocuments)

	s.Logger.Info("normalizedDocuments", len(foundDocuments))

	if err != nil {
		s.Logger.Error(err.Error())
		return domain.InvertedIndex{}, err
	}

	invertedIndex := domain.InvertedIndex{
		Df:                      map[string]int{},
		NormalizedDocumentFound: make(map[string]domain.NormalizedDocument, 0),
		CorpusSize:              0,
	}

	for _, term := range localQuery {
		for _, document := range normalizedDocuments {
			qtd := document.Tf[term]
			invertedIndex.Df[term] += qtd
			//if qtd > 0 {
			//	invertedIndex.Df[term] += 1
			//}

		}
	}

	for _, document := range normalizedDocuments {
		invertedIndex.NormalizedDocumentFound[document.Id] = document
	}

	invertedIndex.CorpusSize = len(invertedIndex.NormalizedDocumentFound)
	invertedIndex.Idf = nlp.CalcIdf(invertedIndex.Df, invertedIndex.CorpusSize)

	return invertedIndex, nil

}

func (s Search) FindDocuments(localQuery []string) map[string]int8 {

	var foundDocuments = make(map[string]int8, 0)

	for _, term := range localQuery {
		documents, _ := s.IndexMemoryRepository.FindByTerm(term)
		for _, docID := range documents {
			foundDocuments[docID] = 0
		}
	}

	return foundDocuments
}

func (s Search) LexicalSearch(query string) ([]domain.ScoreResult, error) {

	s.Logger.Info("LexicalSearch")
	start := time.Now()
	localQuery := nlp.Tokenizer(query, true) //nlp.RemoveStopWords(nlp.Tokenizer(query, true), "en")
	foundDocuments := s.FindDocuments(localQuery)
	duration := time.Since(start)
	s.Logger.Info("Duração: ", duration)

	invertedIndex, err := s.MakeInvertedIndex(localQuery, foundDocuments)

	s.Logger.Info("Total inverted index: ", len(invertedIndex.NormalizedDocumentFound))

	if err != nil {
		s.Logger.Error(err.Error())
		return nil, exception.ThrowValidationError(err.Error())
	}

	scoreResult := nlp.SortDesc(nlp.ScoreBM25(localQuery, &invertedIndex), 10)

	s.Logger.Info("scoreResult", scoreResult)

	return scoreResult, nil

}
