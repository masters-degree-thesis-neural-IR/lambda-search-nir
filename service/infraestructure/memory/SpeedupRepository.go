package memory

import (
	"context"
	"encoding/json"
	"google.golang.org/grpc"
	"lambda-search-nir/service/application/domain"
	"lambda-search-nir/service/application/exception"
	"lambda-search-nir/service/application/repositories"
	"lambda-search-nir/service/infraestructure/speedup"
	"log"
	"time"
)

type SpeedupRepository struct {
	Host string
}

func NewSpeedupRepository(host string) repositories.IndexMemoryRepository {

	return &SpeedupRepository{
		Host: host,
	}
}

func (r *SpeedupRepository) FindByTerm(term string) ([]string, error) {

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(r.Host, grpc.WithInsecure())

	if err != nil {
		log.Println(err.Error())
		return nil, exception.ThrowValidationError("Not is possible connect to RCP Server.")
	}
	defer conn.Close()

	client := speedup.NewDataServiceClient(conn)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	response, err := client.GetData(ctx, &speedup.RequestDataKey{
		Key: term,
	})

	if ctx.Err() == context.Canceled {
		return nil, exception.ThrowValidationError("RPC Client cancelled, abandoning.")
	}

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	if response.GetException() != "" {
		log.Println(err.Error())
		return nil, exception.ThrowValidationError(response.GetException())
	}

	var locDocuments []string
	json.Unmarshal([]byte(response.GetValue()), &locDocuments)
	return locDocuments, nil

}

func (r *SpeedupRepository) Save(term string, document domain.NormalizedDocument) error {

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(r.Host, grpc.WithInsecure())
	if err != nil {
		log.Println(err.Error())
		return exception.ThrowValidationError("Not is possible connect to RCP Server.")
	}
	defer conn.Close()

	client := speedup.NewDataServiceClient(conn)

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	value, err := json.Marshal(document)
	response, err := client.SetData(ctx, &speedup.RequestDataKeyValue{
		Key:   term,
		Value: string(value),
	})

	if ctx.Err() == context.Canceled {
		log.Println(err.Error())
		return exception.ThrowValidationError("RPC Client cancelled, abandoning.")
	}

	if err != nil {
		return err
	}

	if response.GetException() != "" {
		log.Println(err.Error())
		return exception.ThrowValidationError(response.GetException())
	}

	return nil

}

func (r *SpeedupRepository) LoadMetricsFromCache(documentIDs map[string]int8) ([]domain.NormalizedDocument, map[string]int8, error) {

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(r.Host, grpc.WithInsecure())
	notFound := make(map[string]int8)
	cacheList := make(map[string]int8)

	for key, value := range documentIDs {
		cacheList["metrics"+key] = value
		notFound[key] = value
	}

	if err != nil {
		return nil, nil, exception.ThrowValidationError("Not is possible connect to RCP Server.")
	}

	defer conn.Close()
	client := speedup.NewDataServiceClient(conn)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	keys := make([]*speedup.RequestDataKey, len(cacheList))

	var i = 0
	for key, _ := range cacheList {
		keys[i] = &speedup.RequestDataKey{
			Key: key,
		}
		i++
	}

	response, err := client.GetsData(ctx, &speedup.RequestDataKeyList{
		RequestDataKeyList: keys,
	})

	if len(response.ResponseDataValueList) == 0 {
		return nil, notFound, nil
	}

	if ctx.Err() == context.Canceled {
		return nil, nil, exception.ThrowValidationError("RPC Client cancelled, abandoning.")
	}

	if err != nil {
		return nil, nil, err
	}

	if response.GetException() != "" {
		return nil, nil, exception.ThrowValidationError(response.GetException())
	}

	normalizedDocuments := make([]domain.NormalizedDocument, 0)
	for _, value := range response.ResponseDataValueList {
		var normalizedDocument domain.NormalizedDocument
		err := json.Unmarshal([]byte(value.Value), &normalizedDocument)
		if err != nil {
			log.Fatal(err.Error())
			return nil, nil, err
		}

		if &normalizedDocument != nil {
			delete(notFound, normalizedDocument.Id)
			normalizedDocuments = append(normalizedDocuments, normalizedDocument)
		}

	}
	return normalizedDocuments, notFound, nil

}
