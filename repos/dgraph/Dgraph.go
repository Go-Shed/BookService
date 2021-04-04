package dgraph

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/mitchellh/mapstructure"
)

type Patch struct {
	Username string `json:"username"`
}

type Variables struct {
	Patch Patch `json:"patch"`
}

type Request struct {
	Query     string    `json:"query"`
	Variables Variables `json:"variables,omitempty"`
	Operation string    `json:"operation,omitempty"`
}

type Dgraph struct {
	URL    string
	Secret string
}

func (dgraph *Dgraph) GetUsers(userId string) User {

	query := Request{
		Query: fmt.Sprintf(`query {
			getUser(username: "%s") {
			  username
			  posts{
				  title
			  }
			}
		  }`, userId)}

	jsonValue, _ := json.Marshal(query)
	request, err := http.NewRequest("POST", dgraph.URL, bytes.NewBuffer(jsonValue))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Dg-Auth", dgraph.Secret)

	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{Timeout: time.Second * 10}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	result, _ := ioutil.ReadAll(response.Body)
	// fmt.Println(string(result))

	var dat DgraphResponse

	if err := json.Unmarshal([]byte(string(result)), &dat); err != nil {
		panic(err)
	}
	data := dat.Data.(map[string]interface{})
	var user User
	mapstructure.Decode(data["getUser"], &user)
	return user
}
