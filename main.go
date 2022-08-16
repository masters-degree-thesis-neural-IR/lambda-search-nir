package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws/session"
	"lambda-search-nir/service/application/exception"
	"lambda-search-nir/service/application/service"
	"lambda-search-nir/service/infraestructure/dto"
	"lambda-search-nir/service/infraestructure/dydb"
	"log"
	"net/http"
)

var TopicArn string
var TableName string
var AwsRegion string

func ErrorHandler(err error) events.APIGatewayProxyResponse {

	switch err.(type) {
	case *exception.ValidationError:

		err, _ := err.(*exception.ValidationError)

		return events.APIGatewayProxyResponse{
			StatusCode: err.StatusCode,
			Headers:    map[string]string{"Content-Type": "text/plain; charset=utf-8"},
			Body:       err.Message,
		}

	default:

		log.Fatalln("Error...: ", err)

		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    map[string]string{"Content-Type": "text/plain; charset=utf-8"},
			Body:       "Internal error",
		}
	}

}

func makeBody(body string) (dto.Document, error) {
	var doc dto.Document
	err := json.Unmarshal([]byte(body), &doc)

	if err != nil {
		return doc, err
	}

	return doc, nil

}

func handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	//if req.HTTPMethod != "POST" {
	//	return events.APIGatewayProxyResponse{
	//		StatusCode: http.StatusBadRequest,
	//		Headers:    map[string]string{"Content-Type": "text/plain; charset=utf-8"},
	//		Body:       "Invalid HTTP Method",
	//	}, nil
	//}
	//
	//document, err := makeBody(req.Body)
	//
	//if err != nil {
	//	return ErrorHandler(err), nil
	//}
	//
	//awsSession, err := session.NewSession(&aws.Config{
	//	Region: aws.String(AwsRegion)},
	//)
	//
	//if err != nil {
	//	return ErrorHandler(err), nil
	//}
	//
	//repository := dydb.NewDocumentRepository(awsSession, TableName)
	//documentEvent := sns.NewDocumentEvent(awsSession, TopicArn)
	//documentService := service.NewDocumentService(documentEvent, repository)
	//err = documentService.CreateDocument(document.Title, document.Body)
	//
	//if err != nil {
	//	return ErrorHandler(err), nil
	//}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Headers:    map[string]string{"Content-Type": "text/plain; charset=utf-8"},
		Body:       "Document created",
	}, nil

}

func main() {
	AwsRegion = "us-east-1"
	TopicArn = "arn:aws:sns:us-east-1:149501088887:mestrado-document-created" //os.Getenv("BAR")
	TableName = "NIR_Index"

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	repository := dydb.NewIndexRepository(sess, TableName)

	service := service.NewSearch(repository)
	query := "plano federal"
	results, err := service.SearchDocument(query)

	if err != nil {
		log.Fatalln(err)
	}

	for _, result := range results {
		log.Println("****** Query ******")
		log.Println("..... Similaridade : ", result.Similarity)
		log.Println("..... ID: ", result.NormalizedDocument.Id)
	}

	//lambda.Start(handler)
}
