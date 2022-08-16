package dto

type QueryResult struct {
	Similarity float64  `json:"similarity"`
	Document   Document `json:"document"`
}
