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

type Variables struct {
	Patch interface{} `json:"patch"`
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

func (dgraph *Dgraph) Do(query Request) (map[string]interface{}, error) {

	jsonValue, _ := json.Marshal(query)
	request, err := http.NewRequest("POST", dgraph.URL, bytes.NewBuffer(jsonValue))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Dg-Auth", dgraph.Secret)

	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{Timeout: time.Second * 5}
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

	response, err := dgraph.Do(query)

	if err != nil {
		return User{}
	}

	var user User
	mapstructure.Decode(response["getUser"], &user)
	return user
}

func (dgraph *Dgraph) CreateUser(newUser User) (User, error) {

	query := Request{
		Query: `mutation addUser($patch: [AddUserInput!]!) {
			addUser(input: $patch) {
			  user {
				username
			  }
			}
		  }`, Variables: Variables{Patch: newUser}, Operation: "addUser"}

	response, err := dgraph.Do(query)

	if err != nil {
		return User{}, err
	}

	data := response["addUser"].(map[string]interface{})
	var user []User
	mapstructure.Decode(data["user"], &user)
	return user[0], nil
}
