package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"lambda-search-nir/service/application/domain"
	"lambda-search-nir/service/application/exception"
	"lambda-search-nir/service/application/logger"
	"lambda-search-nir/service/application/repositories"
	"lambda-search-nir/service/application/service"
	"lambda-search-nir/service/infraestructure/controller"
	"lambda-search-nir/service/infraestructure/dto"
	"lambda-search-nir/service/infraestructure/dydb"
	zaplog "lambda-search-nir/service/infraestructure/log"
	"lambda-search-nir/service/infraestructure/memory"
	"net/http"
	"time"
)

var TableDocument string
var TableMetrics string
var AwsRegion string

var documentMetricsRepository repositories.DocumentMetricsRepository
var documentRepository repositories.DocumentRepository
var log logger.Logger

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

		log.Error(err.Error())

		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    map[string]string{"Content-Type": "text/plain; charset=utf-8"},
			Body:       "Internal error",
		}
	}

}

func makeBody(documentResults []domain.DocumentResult, duration time.Duration) (dto.Result, error) {

	total := len(documentResults)

	var algorithm = "BM25"

	rst := dto.Result{
		Total:          total,
		Duration:       duration.String(),
		Algorithm:      algorithm,
		SemanticSearch: false,
		QueryResults:   make([]dto.QueryResult, total),
	}

	for i, result := range documentResults {
		rst.QueryResults[i] = dto.QueryResult{
			Similarity: result.Similarity,
			Document: dto.Document{
				Id:    result.Document.Id,
				Title: result.Document.Title,
				Body:  result.Document.Body,
			},
		}
	}

	return rst, nil

}

func handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	memoryRepository := memory.NewSpeedupRepository()
	documentService := service.NewDocumentService(log, documentRepository)
	searchService := service.NewSearchService(log, documentMetricsRepository, memoryRepository, documentRepository)
	controller := controller.NewController(documentService, searchService)

	query := req.QueryStringParameters["query"]

	log.Info("Recebeu a requisição: ", query)

	start := time.Now()
	documentResults, err := controller.SearchDocuments(query)
	duration := time.Since(start)

	if err != nil {
		return ErrorHandler(err), nil
	}

	body, err := makeBody(documentResults, duration)
	if err != nil {
		return ErrorHandler(err), nil
	}

	lbody, err := json.Marshal(body)

	if err != nil {
		log.Error(err.Error())
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    map[string]string{"Content-Type": "application/json; charset=utf-8"},
			Body:       "internal error",
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Headers:    map[string]string{"Content-Type": "application/json; charset=utf-8"},
		Body:       string(lbody),
	}, nil

}

func main() {

	AwsRegion = "us-east-1"
	TableDocument = "NIR_Document"
	TableMetrics = "NIR_Metrics"

	log = zaplog.NewZapLogger()
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(AwsRegion)},
	)

	if err != nil {
		log.Fatal(err.Error())
	}

	documentMetricsRepository = dydb.NewDocumentMetricsRepository(awsSession, TableMetrics)
	documentRepository = dydb.NewDocumentRepository(awsSession, TableDocument)

	lambda.Start(handler)
}

func mainrrt() {

	AwsRegion = "us-east-1"
	TableDocument = "NIR_Document"
	TableMetrics = "NIR_Metrics"

	log = zaplog.NewZapLogger()
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(AwsRegion)},
	)

	if err != nil {
		log.Fatal(err.Error())
	}

	documentMetricsRepository = dydb.NewDocumentMetricsRepository(awsSession, TableMetrics)
	documentRepository = dydb.NewDocumentRepository(awsSession, TableDocument)

	memoryRepository := memory.NewSpeedupRepository()
	documentService := service.NewDocumentService(log, documentRepository)
	searchService := service.NewSearchService(log, documentMetricsRepository, memoryRepository, documentRepository)
	controller := controller.NewController(documentService, searchService)

	query := "thermoelastic interaction problems"

	documentResults, err := controller.SearchDocuments(query)

	fmt.Printf("Result %v", documentResults)

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

func mainrt() {

	ids := []string{"100", "20", "2", "154", "145", "12", "47", "987", "12", "14", "2", "5", "69", "70", "45"}

	println(len(ids))
	fmt.Printf("%v", Paginator(ids, 3))

}

func mainxxx() {

	AwsRegion = "us-east-1"
	TableDocument = "NIR_Document"
	TableMetrics = "NIR_Metrics"

	awsSession, _ := session.NewSession(&aws.Config{
		Region: aws.String(AwsRegion)},
	)

	svc := dynamodb.New(awsSession)

	//ids := []string{"10", "20", "30", "150"}

	//exps := make([]expression.ConditionBuilder, len(ids))
	//for i, id := range ids {
	//	xp := expression.Name("Id").Equal(expression.Value(id))
	//
	//	xs := expression.Or(expression.Name("Id").Equal(expression.Value("20")),
	//		expression.Name("Id").Equal(expression.Value("20")))
	//
	//	exps[i] = xs
	//}

	//xx := expression.Or(, expression.Name("Id").LessThan(expression.Value("8")))

	var filt = expression.Name("Id").Equal(expression.Value("500"))
	filt = filt.Or(expression.Name("Id").Equal(expression.Value("20")))
	filt = filt.Or(expression.Name("Id").Equal(expression.Value("30")))

	expr, err := expression.NewBuilder().WithFilter(filt).Build()

	if err != nil {
		println(err.Error())
	}

	//println()

	//Id := aws.String("Id")

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(TableMetrics),
	}

	result, err := svc.Scan(params)
	if err != nil {
		println(err.Error())
	}

	println(len(result.Items))
}
