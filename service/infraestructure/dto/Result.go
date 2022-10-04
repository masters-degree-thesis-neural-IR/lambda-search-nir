package dto

type Result struct {
	Total          int           `json:"total"`
	Duration       string        `json:"duration"`
	Algorithm      string        `json:"algorithm"`
	SemanticSearch bool          `json:"semanticSearch"`
	QueryResults   []QueryResult `json:"queryResults"`
}
