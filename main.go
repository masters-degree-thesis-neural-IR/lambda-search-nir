package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"lambda-search-nir/service/application/domain"
	"lambda-search-nir/service/application/exception"
	"lambda-search-nir/service/application/service"
	"lambda-search-nir/service/infraestructure/dto"
	"lambda-search-nir/service/infraestructure/dydb"
	"net/http"
	"time"
)

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

		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    map[string]string{"Content-Type": "text/plain; charset=utf-8"},
			Body:       "Internal error",
		}
	}

}

func makeBody(results []domain.QueryResult, duration time.Duration) (string, error) {

	rst := dto.Result{
		Duration:     duration.String(),
		QueryResults: make([]dto.QueryResult, len(results)),
	}

	for i, result := range results {
		rst.QueryResults[i] = dto.QueryResult{
			Similarity: result.Similarity,
			Document: dto.Document{
				Id:    result.Document.Id,
				Title: result.Document.Title,
				Body:  result.Document.Body,
			},
		}
	}

	body, err := json.Marshal(rst)
	return string(body), err

}

func handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	if req.HTTPMethod != "GET" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers:    map[string]string{"Content-Type": "text/plain; charset=utf-8"},
			Body:       "Invalid HTTP Method",
		}, nil
	}

	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(AwsRegion)},
	)

	if err != nil {
		return ErrorHandler(err), nil
	}

	documentRepository := dydb.NewDocumentRepository(awsSession, "NIR_Document")
	indexRepository := dydb.NewIndexRepository(awsSession, TableName)

	service := service.NewSearch(indexRepository, documentRepository)
	query := req.QueryStringParameters["query"]

	start := time.Now()
	results, err := service.SearchDocument(query)
	duration := time.Since(start)

	if err != nil {
		return ErrorHandler(err), nil
	}

	body, err := makeBody(results, duration)

	if err != nil {
		return ErrorHandler(err), nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Headers:    map[string]string{"Content-Type": "application/json; charset=utf-8"},
		Body:       body,
	}, nil

}

func main() {

	AwsRegion = "us-east-1"
	TableName = "NIR_Index"

	lambda.Start(handler)
}

func maker() {
	//sess := session.Must(session.NewSessionWithOptions(session.Options{
	//	SharedConfigState: session.SharedConfigEnable,
	//}))

}
