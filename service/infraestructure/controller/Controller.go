package controller

import (
	"fmt"
	"lambda-search-nir/service/application/domain"
	"lambda-search-nir/service/application/usecases"
)

type Controller struct {
	DocumentService usecases.DocumentUc
	Search          usecases.SearchUc
}

func NewController(documentService usecases.DocumentUc, search usecases.SearchUc) Controller {
	return Controller{
		DocumentService: documentService,
		Search:          search,
	}
}

func (c *Controller) SearchDocuments(query string) ([]domain.DocumentResult, error) {

	var err error
	scoreResult, err := c.Search.LexicalSearch(query)
	if err != nil {
		return []domain.DocumentResult{}, err
	}

	for _, res := range scoreResult {
		println(res.DocumentID)
		fmt.Printf("%v\n", res.Similarity)
	}

	return c.DocumentService.LoadDocuments(scoreResult)

}
