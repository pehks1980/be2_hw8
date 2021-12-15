package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func LoadData(es *elasticsearch.Client) {
	var spacecrafts []map[string]interface{}
	pageNumber := 0
	totalPages := 5
	totalPages1 := 0
	for {
		response, _ := http.Get("http://stapi.co/api/v1/rest/spacecraft/search?pageSize=100&pageNumber=" + strconv.Itoa(pageNumber))

		body, _ := ioutil.ReadAll(response.Body)

		defer response.Body.Close()

		var result map[string]interface{}

		json.Unmarshal(body, &result)

		page := result["page"].(map[string]interface{})

		totalPages1 = int(page["totalPages"].(float64))

		crafts := result["spacecrafts"].([]interface{})
		for _, craftInterface := range crafts {
			craft := craftInterface.(map[string]interface{})
			spacecrafts = append(spacecrafts, craft)
		}
		pageNumber++
		if pageNumber >= totalPages {
			break
		}
	}

	for _, data := range spacecrafts {
		uid, _ := data["uid"].(string)
		jsonString, _ := json.Marshal(data)
		request := esapi.IndexRequest{Index: "stsc", DocumentID: uid, Body: strings.NewReader(string(jsonString))}
		request.Do(context.Background(), es)
	}
	fmt.Printf(" %d spacecraft read ,(total was %d )", len(spacecrafts), totalPages1 )
}

func Get(id string, es *elasticsearch.Client) {
	id = strings.TrimSuffix(id, "\n")
	request := esapi.GetRequest{Index: "stsc", DocumentID: id}
	response, err := request.Do(context.Background(), es)
	if err != nil {
		log.Printf("error occured at getting craft with ID = %s, %v", id, err )
		return
	}
	log.Printf("resp statuscode = %d", response.StatusCode)
	if response.StatusCode != 200 {
		return
	}
	var results map[string]interface{}
	json.NewDecoder(response.Body).Decode(&results)
	Print(results["_source"].(map[string]interface{}))
}

func Print(spacecraft map[string]interface{}) {
	name := spacecraft["name"]
	status := ""
	if spacecraft["status"] != nil {
		status = "- " + spacecraft["status"].(string)
	}
	registry := ""
	if spacecraft["registry"] != nil {
		registry = "- " + spacecraft["registry"].(string)
	}
	class := ""
	if spacecraft["spacecraftClass"] != nil {
		class = "- " + spacecraft["spacecraftClass"].(map[string]interface{})["name"].(string)
	}
	fmt.Println(name, registry, class, status)
}

func Search(key, value, querytype string, es *elasticsearch.Client) {
	key = strings.TrimSuffix(key, "\n")
	value = strings.TrimSuffix(value, "\n")
	var buffer bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			querytype: map[string]interface{}{
				key: value,
			},
		},
	}

	json.NewEncoder(&buffer).Encode(query)
	response, err := es.Search(es.Search.WithIndex("stsc"), es.Search.WithBody(&buffer))
	if err != nil {
		log.Printf("error occured at find with request  %v %v", buffer, err )
		return
	}
	log.Printf("resp statuscode = %d", response.StatusCode)
	if response.StatusCode != 200 {
		return
	}
	var result map[string]interface{}
	json.NewDecoder(response.Body).Decode(&result)

	for _, hit := range
		result["hits"].(map[string]interface{})["hits"].([]interface{}) {
			craft := hit.(map[string]interface{})["_source"].(map[string]interface{})
			Print(craft)
	}
}