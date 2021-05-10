package dgraph

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Variables struct {
	Patch interface{} `json:"patch"`
}

type Request struct {
	Query     string    `json:"query"`
	Variables Variables `json:"variables,omitempty"`
	Operation string    `json:"operation,omitempty"`
}

type DgraphResponse struct {
	Data interface{} `json:"data"`
}

type Dgraph struct {
}

const URL = "https://patient-resonance.ap-south-1.aws.cloud.dgraph.io/graphql"
const secret = "MTMyYWIwY2E3YTc0ODQ3MTJkNWRjMDkxYzA4MWJmN2U="

func (dgraph *Dgraph) Do(query Request) (map[string]interface{}, error) {

	jsonValue, _ := json.Marshal(query)
	request, err := http.NewRequest("POST", URL, bytes.NewBuffer(jsonValue))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Dg-Auth", secret)

	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{Timeout: time.Second * 10}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	result, _ := ioutil.ReadAll(response.Body)
	// fmt.Println(string(result))

	var dgraphResponse DgraphResponse
	if err := json.Unmarshal([]byte(string(result)), &dgraphResponse); err != nil {
		panic(err)
	}
	data := dgraphResponse.Data.(map[string]interface{})
	return data, nil
}
