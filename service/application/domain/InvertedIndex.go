package domain

type InvertedIndex struct {
	CorpusSize              int
	Df                      map[string]int
	Idf                     map[string]float64
	NormalizedDocumentFound map[string]NormalizedDocument
}
