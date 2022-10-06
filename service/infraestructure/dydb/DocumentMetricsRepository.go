package dydb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"lambda-search-nir/service/application/domain"
	"lambda-search-nir/service/application/repositories"
)

type DocumentMetricsRepository struct {
	AwsSession       *session.Session
	TableName        string
	MemoryRepository repositories.IndexMemoryRepository
}

func NewDocumentMetricsRepository(awsSession *session.Session, tableName string, memoryRepository repositories.IndexMemoryRepository) repositories.DocumentMetricsRepository {
	return DocumentMetricsRepository{
		AwsSession:       awsSession,
		TableName:        tableName,
		MemoryRepository: memoryRepository,
	}
}

func (d DocumentMetricsRepository) FindByDocumentIDs(documentIDs map[string]int8) ([]domain.NormalizedDocument, error) {

	var normalizedDocuments, notfound, err = d.MemoryRepository.LoadMetricsFromCache(documentIDs)
	if err != nil {
		return nil, err
	}

	nocache := make([]string, 0)
	for id, _ := range notfound {
		nocache = append(nocache, id)
	}

	paginator := Paginator(nocache, 100)

	for _, page := range paginator {

		if len(page) == 0 {
			continue
		}

		var filter = expression.Name("Id").Equal(expression.Value(page[0]))
		for _, id := range page[1:] {
			filter = filter.Or(expression.Name("Id").Equal(expression.Value(id)))
		}

		expr, err := expression.NewBuilder().WithFilter(filter).Build()

		if err != nil {
			return nil, err
		}

		params := &dynamodb.ScanInput{
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			FilterExpression:          expr.Filter(),
			ProjectionExpression:      expr.Projection(),
			TableName:                 aws.String(d.TableName),
		}

		svc := dynamodb.New(d.AwsSession)
		result, err := svc.Scan(params)
		if err != nil {
			return nil, err
		}

		for _, item := range result.Items {
			var normalizedDocument domain.NormalizedDocument
			err = dynamodbattribute.UnmarshalMap(item, &normalizedDocument)
			if err != nil {
				return nil, err
			}

			d.MemoryRepository.Save("metrics"+normalizedDocument.Id, normalizedDocument)
			normalizedDocuments = append(normalizedDocuments, normalizedDocument)
		}

	}

	return normalizedDocuments, nil

}

func Paginator(ids []string, ln int) [][]string {
	size := len(ids)
	slice := len(ids) / ln

	if slice == 0 {
		var x = make([][]string, 0)
		return append(x, ids)
	}

	var start = 0
	var end = size / slice
	var list [][]string
	for i := 0; i < slice+1; i++ {
		arr := ids[start:end]

		if len(arr) > 0 {
			list = append(list, arr)
		}

		start = end
		end = end + end

		if end > size {
			end = end - (end - size)
		}
	}

	return list
}
