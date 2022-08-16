package service

import (
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"lambda-search-nir/service/application/domain"
	"lambda-search-nir/service/application/exception"
	"lambda-search-nir/service/application/repositories"
	"lambda-search-nir/service/application/stopwords"
	"lambda-search-nir/service/application/usecases"
	"strings"
	"unicode"
)

type IndexService struct {
	IndexRepository repositories.IndexRepository
}

func NewIndexService(indexRepository repositories.IndexRepository) usecases.CreateIndexUc {
	var c usecases.CreateIndexUc = IndexService{
		IndexRepository: indexRepository,
	}
	return c
}

func (i IndexService) CreateIndex(id string, title string, body string) error {

	tokens := Tokenizer(body, true)
	normalizedTokens, err := RemoveStopWords(tokens, "pt")

	if err != nil {
		return err
	}

	document := domain.Document{
		Id:     id,
		Length: len(normalizedTokens),
		Tf:     TermFrequency(normalizedTokens),
	}

	for _, term := range normalizedTokens {

		index, err := i.IndexRepository.FindByTerm(term)

		if err != nil {
			return err
		}

		if index != nil {
			documentList := index.Documents
			if NotContains(document, documentList) {
				index.Term = term
				index.Documents = append(documentList, document)
				i.IndexRepository.Update(*index)
			}
		} else {
			index := domain.Index{
				Term:      term,
				Documents: []domain.Document{document},
			}
			i.IndexRepository.Save(index)
		}
	}

	return nil
}

func NotContains(document domain.Document, documents []domain.Document) bool {

	for _, doc := range documents {
		if doc.Id == document.Id {
			return false
		}
	}

	return true
}

func RemoveAccents(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, e := transform.String(t, s)
	if e != nil {
		panic(e)
	}
	return output
}

func Tokenizer(document string, normalize bool) []string {

	var temp = strings.ReplaceAll(document, ",", "")
	temp = strings.ReplaceAll(temp, ".", "")

	fields := strings.Fields(temp)

	if normalize {

		localSlice := make([]string, len(fields))
		for i, token := range fields {
			localSlice[i] = strings.ToLower(RemoveAccents(token))
		}

		return localSlice
	}

	return fields

}

func StopWordLang(lang string) (map[string]bool, error) {

	if lang == "en" {
		return stopwords.English, nil
	}

	if lang == "pt" {
		return stopwords.Portuguese, nil
	}

	return nil, *exception.ThrowValidationError("Not found language from stop word")
}

func RemoveStopWords(tokens []string, lang string) ([]string, error) {

	stopWordLang, err := StopWordLang(lang)

	if err != nil {
		return nil, err
	}

	if len(tokens) == 0 {
		return make([]string, 0), nil
	}

	var localSlice = make([]string, 0)

	for _, token := range tokens {
		if !stopWordLang[token] {
			localSlice = append(localSlice, token)
		}
	}

	return localSlice, nil

}

func TermFrequency(tokens []string) map[string]int {

	localMap := make(map[string]int)

	for _, token := range tokens {

		if localMap[token] == 0 {
			localMap[token] = 1
		} else {
			localMap[token] = localMap[token] + 1
		}
	}

	return localMap

}
