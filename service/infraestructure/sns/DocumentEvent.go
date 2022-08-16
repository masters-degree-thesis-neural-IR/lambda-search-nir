package sns

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"lambda-search-nir/service/application/domain"
	"lambda-search-nir/service/application/ports"
	"log"
)

type DocumentEvent struct {
	AwsSession *session.Session
	TopicArn   string
}

func NewDocumentEvent(awsSession *session.Session, topicArn string) ports.DocumentEvent {
	var c ports.DocumentEvent = DocumentEvent{
		AwsSession: awsSession,
		TopicArn:   topicArn,
	}
	return c
}

func (d DocumentEvent) Created(document domain.Document) error {

	doc, err := json.Marshal(document)

	if err != nil {
		log.Fatalln("Error...: ", err)
		return err
	}

	svc := sns.New(d.AwsSession)
	result, err := svc.Publish(&sns.PublishInput{
		Message:  aws.String(string(doc)),
		TopicArn: aws.String(d.TopicArn),
	})

	if err != nil {
		log.Fatalln("Error...: ", err)
		return err
	}

	log.Println("Message ID...: ", *result.MessageId)

	return nil
}
