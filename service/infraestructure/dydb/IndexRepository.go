package dydb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"lambda-search-nir/service/application/domain"
	"lambda-search-nir/service/application/repositories"
	"log"
)

type IndexRepository struct {
	AwsSession *session.Session
	TableName  string
}

func NewIndexRepository(awsSession *session.Session, tableName string) repositories.IndexRepository {
	return IndexRepository{
		AwsSession: awsSession,
		TableName:  tableName,
	}
}

func (i IndexRepository) FindByTerm(term string) (*domain.Index, error) {

	svc := dynamodb.New(i.AwsSession)

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(i.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Term": {
				S: aws.String(term),
			},
		},
	})

	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, nil
	}

	index := domain.Index{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &index)
	if err != nil {
		return nil, err
	}

	return &index, nil

}

func (i IndexRepository) Update(index domain.Index) error {

	docs, err := dynamodbattribute.MarshalList(index.Documents)

	if err != nil {
		return err
	}

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":r": {
				L: docs,
			},
		},
		TableName: aws.String(i.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Term": {
				S: aws.String(index.Term),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set Documents = :r"),
	}

	svc := dynamodb.New(i.AwsSession)
	_, err = svc.UpdateItem(input)
	if err != nil {
		log.Fatalln("Error...: ", err)
		return err
	}

	return err

}

func (i IndexRepository) Save(index domain.Index) error {

	item, err := dynamodbattribute.MarshalMap(index)

	if err != nil {
		log.Fatalln("Error...: ", err)
		return err
	}

	svc := dynamodb.New(i.AwsSession)
	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(i.TableName),
	}

	_, err = svc.PutItem(input)

	if err != nil {
		log.Fatalln("Error...: ", err)
		return err
	}

	return nil
}
