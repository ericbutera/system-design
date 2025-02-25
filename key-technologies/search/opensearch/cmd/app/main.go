package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/opensearch-project/opensearch-go"
)

func main() {
	client, err := opensearch.NewClient(opensearch.Config{
		Addresses: []string{"http://opensearch:9200"},
	})
	if err != nil {
		log.Fatalf("Error creating OpenSearch client: %s", err)
	}

	// createIndex(client)
	insertAsset(client)
	searchAssets(client, "msdef-601")
}

func insertAsset(client *opensearch.Client) {
	asset := map[string]interface{}{
		"vendor_id":   "msdef-601",
		"ip_address":  "192.168.50.10",
		"hostname":    "laptop-user1",
		"os":          "Windows 10",
		"integration": "defender",
	}

	jsonData, _ := json.Marshal(asset)

	res, err := client.Index("assets", bytes.NewReader(jsonData), client.Index.WithDocumentID("msdef-601"))
	if err != nil {
		log.Fatalf("Error indexing document: %s", err)
	}
	defer res.Body.Close()

	fmt.Println("Document inserted:", res)
}

func searchAssets(client *opensearch.Client, vendorID string) {
	query := fmt.Sprintf(`{"query": {"match": {"vendor_id": "%s"}}}`, vendorID)

	res, err := client.Search(
		client.Search.WithIndex("assets"),
		client.Search.WithBody(bytes.NewReader([]byte(query))),
	)
	if err != nil {
		log.Fatalf("Error searching assets: %s", err)
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println("Search results:", string(body))
}

func createIndex(client *opensearch.Client) {
	indexName := "assets"

	mapping := `{
        "settings": {
            "number_of_shards": 1,
            "number_of_replicas": 0
        },
        "mappings": {
            "properties": {
                "vendor_id": { "type": "keyword" },
                "ip_address": { "type": "ip" },
                "hostname": { "type": "text" },
                "os": { "type": "text" },
                "integration": { "type": "keyword" }
            }
        }
    }`

	res, err := client.Indices.Create(indexName, client.Indices.Create.WithBody(bytes.NewReader([]byte(mapping))))
	if err != nil {
		log.Fatalf("Error creating index: %s", err)
	}
	defer res.Body.Close()

	fmt.Println("Index creation response:", res)
}
