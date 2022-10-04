package usecases

type WordEmbeddingUc interface {
	CreateEmbedding(id, title, body string)
}
