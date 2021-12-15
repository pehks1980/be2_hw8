package repository

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/estransport"
	"github.com/google/uuid"
	uu "github.com/satori/go.uuid"
	"os"
	"pehks1980/be2_hw81/internal/pkg/model"

	"net"
	"net/http"
	"pehks1980/be2_hw81/internal/app/endpoint"
	"time"
)

// PgElastic - init elastic struct holds ptr to es client
type EsRepo struct {
	URL string
	ES 	*elasticsearch.Client
	IDX string
}


// New Init of pg driver
func (esr *EsRepo) New(ctx context.Context, filename, filename1 string) endpoint.RepoIf {
	// Строка адреса подключения к эластику + рабочий индекс
	url := filename // "http://192.168.1.210:9200"

	cfg := elasticsearch.Config{
		Addresses: []string{
			url,
		},
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second,
			DialContext:           (&net.Dialer{Timeout: time.Second}).DialContext,
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS11,
			},
		},
		Logger: &estransport.ColorLogger{
			Output: os.Stdout,
		},

	}

	es, _ := elasticsearch.NewClient(cfg)
	log.Print("Elastic server =", es.Transport.(*estransport.Client).URLs())
	log.Print("Using Index = ", filename1)

	return &EsRepo{
		URL:    url,
		ES: 	es,
		IDX:	filename1, //elastic index name
	}
}

func (esr *EsRepo) CloseConn() {
	panic("implement me")
}

// crud for goods as per be2_hw8

//FindGood - find data using query elastic
func (esr *EsRepo) FindGood(ctx context.Context, key string) ([]model.Good, error) {
	var buffer bytes.Buffer

	query1 := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query": key,
				"type": "best_fields",
				"fields": []string{"name", "id", "qty", "price"},
				"tie_breaker": 0.3,
				"lenient": true, //You can set lenient to true in your multi_match query to ignore exceptions caused by data-type mismatches
			}, // eg english lenient - mild smooth easy
		},
	}

	json.NewEncoder(&buffer).Encode(query1)
	/*
		body := fmt.Sprintf(
			`{ "query": { "multi_match" : { "query": 'socks', "type": "best_fields", "fields": [ "name", "qty", "price" ], "tie_breaker": 0.3 } } }`,
			`{"query": {"match": {"%s" : "%s"}}}`,
			)
		log.Printf("q1=%s,q2=%s",query1,body)*/
	response, err := esr.ES.Search(
		//esr.ES.Search.WithContext(ctx),
		esr.ES.Search.WithIndex(esr.IDX),
		esr.ES.Search.WithBody(&buffer),
		esr.ES.Search.WithPretty(), //nice output of json
	)
	if err != nil{
		log.Printf("elastic search error: %v",err)
		return nil, err
	}
	if response.StatusCode != 200 {
		log.Printf("elastic search status code: %d resposne %v", response.StatusCode, response)
		return nil, nil
	}

	var result map[string]interface{}
	json.NewDecoder(response.Body).Decode(&result)

	var goods []model.Good
	for _, hit := range result["hits"].(map[string]interface{})["hits"].([]interface{}) {
		// goodif elastic result one row
		// repack to struct goods - array of structs
		goodif := hit.(map[string]interface{})["_source"].(map[string]interface{})
		id, _ := uu.FromString(goodif["id"].(string))
		var good = model.Good{
			ID: uuid.UUID(id),
			Name: goodif["name"].(string),
			Price: int(goodif["price"].(float64)),
			Qty: int(goodif["qty"].(float64)),
		}
		goods = append(goods, good)
		log.Printf("found good = %v",good)
	}
	return goods, nil
}

// AddUpdGood - add change good item
func (esr *EsRepo) AddUpdGood(ctx context.Context, good model.Good) (string, error) {
	jsonString, _ := json.Marshal(good)

	request := esapi.IndexRequest{
		Index: esr.IDX, DocumentID: good.ID.String(),
		Body: strings.NewReader(string(jsonString)),
	}

	res, err := request.Do(ctx, esr.ES)
	if err != nil {
		log.Printf("error add put to elastic %v",err)
		return res.String(), err
	}

	return res.String(), nil
}

//GetGood - get good from elastic
func (esr *EsRepo) GetGood(ctx context.Context, title string) (model.Good, error) {
	request := esapi.GetRequest{Index: esr.IDX, DocumentID: title}
	response, err := request.Do(ctx, esr.ES)
	if err != nil{
		log.Printf("get elastic error: %v",err)
		return model.Good{}, err
	}
	var goodif, results map[string]interface{}
	json.NewDecoder(response.Body).Decode(&results)
	goodif = results["_source"].(map[string]interface{})

	id, _ := uu.FromString(goodif["id"].(string))
	fmt.Println(id)
	var good = model.Good{
		ID: uuid.UUID(id),
		Name: goodif["name"].(string),
		Price: int(goodif["price"].(float64)),
		Qty: int(goodif["qty"].(float64)),
	}

	log.Printf("struct good= %v",good)

	return good,nil
}

//DelGood - delete good from elastic
func (esr *EsRepo) DelGood(ctx context.Context, id uuid.UUID) error {
	request := esapi.DeleteRequest{Index: esr.IDX, DocumentID: id.String()}
	response, err := request.Do(ctx, esr.ES)
	if err != nil{
		log.Printf("get elastic error: %v",err)
		return err
	}
	log.Printf("elastic delete response=%v",response)
	return nil
}

func (esr *EsRepo) AuthUser(ctx context.Context, user model.User) (string, error) {
	panic("implement me")
}

func (esr *EsRepo) GetUser(ctx context.Context, name string) (model.User, error) {
	panic("implement me")
}

func (esr *EsRepo) AddUpdUser(ctx context.Context, user model.User) (string, error) {
	panic("implement me")
}

func (esr *EsRepo) DelUser(ctx context.Context, id uuid.UUID) error {
	panic("implement me")
}

func (esr *EsRepo) GetUserEnvs(ctx context.Context, name string) (model.Envs, error) {
	panic("implement me")
}

func (esr *EsRepo) AddUpdEnv(ctx context.Context, env model.Environment) (string, error) {
	panic("implement me")
}

func (esr *EsRepo) GetEnv(ctx context.Context, title string) (model.Environment, error) {
	panic("implement me")
}

func (esr *EsRepo) DelEnv(ctx context.Context, id uuid.UUID) error {
	panic("implement me")
}

func (esr *EsRepo) GetEnvUsers(ctx context.Context, title string) (model.Users, error) {
	panic("implement me")
}
