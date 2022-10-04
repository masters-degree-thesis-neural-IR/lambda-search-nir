package nlp

import (
	"fmt"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"lambda-search-nir/service/application/domain"
	"lambda-search-nir/service/application/exception"
	"lambda-search-nir/service/application/nlp/score"
	"lambda-search-nir/service/application/nlp/stopwords"
	"log"
	"math"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

//func NotContains(document domain.DocumentID, documents []domain.DocumentID) bool {
//
//	for _, doc := range documents {
//		if doc.Id == document.Id {
//			return false
//		}
//	}
//
//	return true
//}

func NotContains(documentID string, documents []string) bool {

	for _, id := range documents {
		if documentID == id {
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

func CleanSpecialCharacters(word string) string {

	reg, err := regexp.Compile("[\\p{P}\\p{S}]+")

	if err != nil {
		log.Fatal(err)
	}

	return reg.ReplaceAllString(word, "")
}

func Tokenizer(document string, normalize bool) []string {

	fields := strings.Fields(document)

	if normalize {

		var localSlice = make([]string, 0)
		for _, token := range fields {
			var tempToken = strings.ToLower(RemoveAccents(CleanSpecialCharacters(token)))
			tempToken = strings.TrimSpace(tempToken)
			if len(tempToken) > 2 {
				localSlice = append(localSlice, tempToken)
			}
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

	return nil, exception.ThrowValidationError("Not found language from stop word")
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

func calcIdf(df map[string]int, corpusSize int) map[string]float64 {

	epsilon := 0.25

	idf := make(map[string]float64)
	var negativeIdfs = make([]string, 0)
	var idfSum float64 = 0

	for term, frequency := range df {

		corpusSize := float64(corpusSize)
		freq := float64(frequency)
		lidf := math.Log(1 + (corpusSize-freq+0.5)/freq + 0.5)
		if math.IsInf(lidf, 1) {
			idfSum += 0
		} else {
			idfSum += lidf
		}
		idf[term] = lidf
		if lidf < 0 {
			println("Tem negativo")
			negativeIdfs = append(negativeIdfs, term)
		}
	}

	n := len(negativeIdfs)
	fmt.Printf("IDF sum: %v", idfSum)
	fmt.Printf("N sum: %v", n)
	averageIdf := idfSum / float64(n)
	eps := (epsilon * averageIdf) * -1

	fmt.Printf("Tem eps: %v", eps)

	for _, term := range negativeIdfs {
		print("Tem negativo")
		idf[term] = eps
	}

	//fmt.Printf("%v", idf)

	return idf

}

func CalcIdf(df map[string]int, corpusSize int) map[string]float64 {

	idf := make(map[string]float64)

	for term, frequency := range df {
		//idf[term] = math.Log(1 + (corpus_size-freq+0.5)/(freq+0.5))
		freq := float64(frequency)
		corpusSize := float64(corpusSize)
		idf[term] = math.Log(1 + (corpusSize-freq+0.5)/(freq+0.5))
	}

	return idf

}

func ScoreCosineSimilarity(query []float64, documentsEmbedding []domain.DocumentEmbedding) []domain.ScoreResult {

	queryResults := make([]domain.ScoreResult, len(documentsEmbedding))

	for i, document := range documentsEmbedding {
		score := score.CosineSimilarity(query, document.Embedding)

		queryResults[i] = domain.ScoreResult{
			Similarity: score,
			DocumentID: document.Id,
		}
		i++
	}

	return queryResults
}

func ScoreBM25(query []string, invertedIndex *domain.InvertedIndex) []domain.ScoreResult {

	queryResults := make([]domain.ScoreResult, invertedIndex.CorpusSize)

	var i = 0
	for _, doc := range invertedIndex.NormalizedDocumentFound {

		score := score.BM25(query, &doc, invertedIndex.Idf, invertedIndex.CorpusSize, 0.75, 1.2)

		queryResults[i] = domain.ScoreResult{
			Similarity: score,
			DocumentID: doc.Id,
		}
		i++

	}

	return queryResults
}

func SortDesc(results []domain.ScoreResult, top int) []domain.ScoreResult {

	sort.Slice(results, func(i, j int) bool {
		return results[i].Similarity > results[j].Similarity
	})

	var maxScore = make([]domain.ScoreResult, 0)

	count := 0
	for _, document := range results {
		if count == top {
			break
		}
		maxScore = append(maxScore, document)
		count++
	}

	return maxScore
}
