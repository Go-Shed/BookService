package query

import (
	"fmt"
	"shed/bookservice/repos/dgraph"
	"shed/bookservice/repos/dgraph/model"

	"github.com/mitchellh/mapstructure"
)

type UserRepo struct {
	client dgraph.Dgraph
}

func NewUserRepo() UserRepo {
	return UserRepo{client: dgraph.Dgraph{
		URL:    "https://billowing-night.ap-south-1.aws.cloud.dgraph.io/graphql",
		Secret: "ZTE4YjRhNGEwYTAxNWRiYjMwMGI4YmEyMjc1ODU5ZmI=",
	}}
}

func (repo UserRepo) GetUsers(userId string) model.User {
	client := repo.client
	query := dgraph.Request{
		Query: fmt.Sprintf(`query {
			getUser(username: "%s") {
			  username
			  posts{
				  title
				  text
			  }
			}
		  }`, userId)}

	response, err := client.Do(query)

	if err != nil {
		return model.User{}
	}

	var user model.User
	mapstructure.Decode(response["getUser"], &user)
	return user
}

func (repo UserRepo) CreateUser(newUser model.User) (model.User, error) {

	client := repo.client
	query := dgraph.Request{
		Query: `mutation addUser($patch: [AddUserInput!]!) {
			addUser(input: $patch) {
			  user {
				username
			  }
			}
		  }`, Variables: dgraph.Variables{Patch: newUser}, Operation: "addUser"}

	response, err := client.Do(query)

	if err != nil {
		return model.User{}, err
	}

	data := response["addUser"].(map[string]interface{})
	var user []model.User
	mapstructure.Decode(data["user"], &user)
	return user[0], nil
}
