package ports

type WordEmbedding interface {
	Generate(sentence string) ([]float64, error)
}
