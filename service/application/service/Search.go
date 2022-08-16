package service

import (
	"lambda-search-nir/service/application/domain"
	"lambda-search-nir/service/application/exception"
	"lambda-search-nir/service/application/nlp"
	"lambda-search-nir/service/application/repositories"
	"lambda-search-nir/service/application/usecases"
	"log"
)

type Search struct {
	IndexRepository    repositories.IndexRepository
	DocumentRepository repositories.DocumentRepository
}

func NewSearch(indexRepository repositories.IndexRepository, documentRepository repositories.DocumentRepository) usecases.SearchUc {
	return Search{
		IndexRepository:    indexRepository,
		DocumentRepository: documentRepository,
	}
}

func (s Search) MakeInvertedIndex(localQuery []string) (domain.InvertedIndex, error) {

	invertedIndex := domain.InvertedIndex{
		Df:                      map[string]int{},
		Terms:                   make(map[string][]domain.NormalizedDocument, 0),
		NormalizedDocumentFound: make(map[string]domain.NormalizedDocument, 0),
		CorpusSize:              0,
	}

	for _, qTerm := range localQuery {

		index, err := s.IndexRepository.FindByTerm(qTerm)

		if err != nil {
			log.Fatalln("Error....: ", err)
			return invertedIndex, err
		}

		if index != nil {
			invertedIndex.Df[index.Term] = len(index.Documents)
			invertedIndex.Terms[index.Term] = index.Documents

			for _, document := range index.Documents {
				invertedIndex.NormalizedDocumentFound[document.Id] = document
			}
		}
	}

	invertedIndex.CorpusSize = len(invertedIndex.NormalizedDocumentFound)
	invertedIndex.Idf = nlp.CalcIdf(invertedIndex.Df, invertedIndex.CorpusSize)

	return invertedIndex, nil

}

func (s Search) SearchDocument(query string) ([]domain.QueryResult, error) {

	localQuery, _ := nlp.RemoveStopWords(nlp.Tokenizer(query, true), "pt")
	invertedIndex, err := s.MakeInvertedIndex(localQuery)

	if err != nil {
		return nil, *exception.ThrowValidationError(err.Error())
	}

	queryResults := nlp.SortDesc(nlp.ScoreBM25(localQuery, &invertedIndex))
	tempQueryResults := make([]domain.QueryResult, len(queryResults))

	for i, queryResult := range queryResults {
		doc, err := s.DocumentRepository.FindById(queryResult.NormalizedDocument.Id)

		if doc == nil {
			continue
		}

		if err != nil {
			return nil, err
		}

		tempQueryResults[i].Similarity = queryResult.Similarity
		tempQueryResults[i].NormalizedDocument = queryResult.NormalizedDocument
		tempQueryResults[i].Document = *doc

	}

	return tempQueryResults, nil
}
