package score

import (
	"lambda-search-nir/service/application/domain"
)

func BM25(query []string, document *domain.NormalizedDocument, idf map[string]float64, corpusSize int, b float64, k1 float64) float64 {

	var score = 0.0
	docLength := float64(document.Length)
	frequencies := document.Tf
	avgDocLen := docLength / float64(corpusSize)

	for _, term := range query {

		if frequencies[term] == 0 {
			continue
		}

		freq := float64(frequencies[term])
		numerator := idf[term] * freq * (k1 + 1)
		denominator := freq + k1*(1-b+b*docLength/avgDocLen)
		score += numerator / denominator
	}

	return score
}
