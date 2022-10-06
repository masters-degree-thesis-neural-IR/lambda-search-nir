package dydb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"lambda-search-nir/service/application/domain"
	"lambda-search-nir/service/application/repositories"
	"sync"
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

func (d DocumentRepository) LoadDocument(id string) domain.Document {

	svc := dynamodb.New(d.AwsSession)
	result, _ := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(d.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(id),
			},
		},
	})

	document := domain.Document{}
	_ = dynamodbattribute.UnmarshalMap(result.Item, &document)

	return document

}

func (d DocumentRepository) FindByDocumentIDs(documentIDs []string) (map[string]domain.Document, error) {

	documents := make(map[string]domain.Document)

	documentIDs = append(documentIDs, "drt")

	var wg sync.WaitGroup
	wg.Add(len(documentIDs))

	mutex := sync.RWMutex{}

	for _, id := range documentIDs {

		go func(id string) {
			defer wg.Done()
			doc := d.LoadDocument(id)

			println(doc.Id)

			mutex.Lock()
			documents[doc.Id] = doc
			mutex.Unlock()
		}(id)
	}

	wg.Wait()

	return documents, nil

}
