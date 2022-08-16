package dydb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"lambda-search-nir/service/application/domain"
	"lambda-search-nir/service/application/repositories"
)

type DocumentRepository struct {
	AwsSession *session.Session
	TableName  string
}

func NewDocumentRepository(awsSession *session.Session, tableName string) repositories.DocumentRepository {
	return DocumentRepository{
		AwsSession: awsSession,
		TableName:  tableName,
	}
}

func (d DocumentRepository) FindById(id string) (*domain.Document, error) {
	svc := dynamodb.New(d.AwsSession)

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(d.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(id),
			},
		},
	})

	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, nil
	}

	document := domain.Document{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &document)
	if err != nil {
		return nil, err
	}

	return &document, nil
}
