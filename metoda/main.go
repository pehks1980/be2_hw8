package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	//"github.com/elastic/elastic-transport-go/v7/etransport"
	"github.com/elastic/go-elasticsearch/v7"
)

func ElasticNewClient() *elasticsearch.Client {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://192.168.1.210:9200",
		},
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second,
			DialContext:           (&net.Dialer{Timeout: time.Second}).DialContext,
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS11,
			},
		},
	}

	es, _ := elasticsearch.NewClient(cfg)
	//log.Print("Elastic server =", es.Transport.(*elastictransport.Client).URLs())

	return es
}

// check loaded elastic index http://192.168.1.210:9200/stsc/_search
func main() {
	es := ElasticNewClient()
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("0) Exit")
		fmt.Println("1) Load spacecraft")
		fmt.Println("2) Get spacecraft")
		fmt.Println("3) Find by match")
		fmt.Println("4) Find by prefix")

		option := ReadText(reader, "Enter option\n")
		if option == "0\n" {
			break
		} else if option == "1\n" {
			LoadData(es)
		} else if option == "2\n" {
			id := ReadText(reader, "Enter spacecraft ID:")
			Get(id, es)
		} else if option == "3\n" {
			key := ReadText(reader, "Enter key:")
			value := ReadText(reader, "Enter value:")
			Search(key, value,"match", es)
		} else if option == "4\n" {
			key := ReadText(reader, "Enter key:")
			value := ReadText(reader, "Enter value:")
			Search(key, value,"prefix", es)
		} else {
			fmt.Println("Invalid option")
		}

	}
}

func ReadText(reader *bufio.Reader, s string) string {
	fmt.Print(s)
	text, _ := reader.ReadString('\n')
	return text
}
