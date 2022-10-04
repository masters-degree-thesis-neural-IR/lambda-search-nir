package usecases

type CreateIndexUc interface {
	CreateIndex(id string, title string, body string) error
}
