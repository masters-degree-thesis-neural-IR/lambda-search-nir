package dydb

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"lambda-search-nir/service/application/domain"
	"lambda-search-nir/service/application/repositories"
)

type DocumentRepository struct {
	Cache      map[string]*domain.Document
	AwsSession *session.Session
	TableName  string
}

func NewDocumentRepository(awsSession *session.Session, tableName string) repositories.DocumentRepository {
	return DocumentRepository{
		AwsSession: awsSession,
		TableName:  tableName,
		Cache:      make(map[string]*domain.Document),
	}
}

func (d DocumentRepository) FindByDocumentIDs(documentIDs []string) (map[string]domain.Document, error) {

	documents := make(map[string]domain.Document)
	var nocache []string

	fmt.Printf("Documents IDs %v\n", len(documentIDs))

	//verify documents in local cache
	for _, id := range documentIDs {

		document := d.Cache[id]
		if document != nil {
			documents[id] = *document
		} else {
			nocache = append(nocache, id)
		}
	}

	var filter = expression.Name("Id").Equal(expression.Value(nocache[0]))
	for _, id := range nocache[1:] {
		filter = filter.Or(expression.Name("Id").Equal(expression.Value(id)))
	}

	expr, err := expression.NewBuilder().WithFilter(filter).Build()

	if err != nil {
		println(err.Error())
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
		println(err.Error())
	}

	for _, item := range result.Items {
		document := domain.Document{}
		err = dynamodbattribute.UnmarshalMap(item, &document)
		if err != nil {
			return nil, err
		}

		d.Cache[document.Id] = &document
		documents[document.Id] = document
	}

	return documents, nil

}
