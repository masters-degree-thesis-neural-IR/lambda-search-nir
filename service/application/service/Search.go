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
	IndexRepository repositories.IndexRepository
}

func NewSearch(indexRepository repositories.IndexRepository) usecases.SearchUc {
	return Search{
		IndexRepository: indexRepository,
	}
}

func (s Search) MakeInvertedIndex(localQuery []string) (*domain.InvertedIndex, error) {

	invertedIndex := domain.InvertedIndex{
		Df:                      map[string]int{},
		Terms:                   make(map[string][]domain.NormalizedDocument, 0),
		NormalizedDocumentFound: make(map[string]domain.NormalizedDocument, 0),
		CorpusSize:              0,
	}

	normalizedDocumentFound := make(map[string]domain.NormalizedDocument)

	for _, qTerm := range localQuery {

		index, err := s.IndexRepository.FindByTerm(qTerm)

		if err != nil {
			log.Fatalln("Error....: ", err)
			return nil, err
		}

		log.Println("Achou o index? ", index != nil)

		if index != nil {
			invertedIndex.Df[index.Term] = len(index.Documents)
			invertedIndex.Terms[index.Term] = index.Documents

			for _, document := range index.Documents {
				normalizedDocumentFound[document.Id] = document
			}
		}
	}

	invertedIndex.CorpusSize = len(normalizedDocumentFound)
	invertedIndex.Idf = nlp.CalcIdf(invertedIndex.Df, invertedIndex.CorpusSize)
	return &invertedIndex, nil

}

func (s Search) SearchDocument(query string) ([]domain.QueryResult, error) {

	localQuery, _ := nlp.RemoveStopWords(nlp.Tokenizer(query, true), "pt")
	invertedIndex, err := s.MakeInvertedIndex(localQuery)

	if err != nil {
		return nil, *exception.ThrowValidationError(err.Error())
	}

	return nlp.ScoreBM25(localQuery, invertedIndex), nil
}
