package dto

type Result struct {
	Duration     string        `json:"duration"`
	QueryResults []QueryResult `json:"queryResults"`
}
