package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type GraphQLRequest struct {
	Query string `json:"query"`
}

func main() {
	fmt.Println("**** PRODUCT-HUNT-GRAPH-VISUALIZE PROJECT ****")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("PRODUCT_HUNT_DEV_TOKEN")
	query := `{
		posts(first: 10) {
			edges {
				node {
					name
					tagline
					url
					votesCount
					thumbnail {
						url
					}
				}
			}
		}
	}`

	requestBody, err := json.Marshal(GraphQLRequest{Query: query})
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://api.producthunt.com/v2/api/graphql", bytes.NewBuffer(requestBody))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var responseData map[string]any
	if err := json.Unmarshal(body, &responseData); err != nil {
		panic(err)
	}

	for _, product := range responseData {
		fmt.Printf("%s", product)
	}
}
