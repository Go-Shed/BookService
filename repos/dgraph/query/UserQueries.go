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
	return UserRepo{client: dgraph.Dgraph{}}
}

func (repo UserRepo) GetUsers(token string) model.User {
	client := repo.client
	query := dgraph.Request{
		Query: fmt.Sprintf(`query {
			getUser(token: "%s") {
			  username
			  posts{
				  title
				  text
			  }
			}
		  }`, token)}

	response, err := client.Do(query)

	if err != nil {
		return model.User{}
	}

	var user model.User
	mapstructure.Decode(response["getUser"], &user)
	return user
}

func (repo UserRepo) FollowUser(user, userToFollow string) error {

	client := repo.client
	query := dgraph.Request{
		Query: fmt.Sprintf(`mutation {
			updateUser(input: {filter: {username: {eq: "%s"}}, set: {following: {username: "%s"}}}){
			  user{
				username
			  }
			}
		  }
		  `, user, userToFollow), Operation: "updateUser"}

	_, err := client.Do(query)
	if err != nil {
		return err
	}
	return nil
}

func (repo UserRepo) UnFollowUser(user, userToUnFollow string) error {

	client := repo.client
	query := dgraph.Request{
		Query: fmt.Sprintf(`mutation {
			updateUser(input: {filter: {username: {eq: "%s"}}, remove: {following: {username: "%s"}}}){
			  user{
				username
			  }
			}
		  }
		  `, user, userToUnFollow), Operation: "updateUser"}

	_, err := client.Do(query)
	if err != nil {
		return err
	}
	return nil
}

func (repo UserRepo) CreateUser(newUser model.User) (model.User, error) {

	client := repo.client
	query := dgraph.Request{
		Query: `mutation addUser($patch: [AddUserInput!]!) {
			addUser(input: $patch) {
			  user {
				userName
			  }
			}
		  }`, Variables: dgraph.Variables{Patch: newUser}, Operation: "addUser"}

	response, err := client.Do(query)

	if err != nil {
		return model.User{}, err
	}

	if response["addUser"] == nil {
		return model.User{}, fmt.Errorf("username already exists")
	}
	data := response["addUser"].(map[string]interface{})
	var user []model.User
	mapstructure.Decode(data["user"], &user)
	return user[0], nil
}

func (repo UserRepo) SearchUser(username string) model.User {

	client := repo.client
	query := dgraph.Request{
		Query: fmt.Sprintf(`{
			queryUser(filter: {userName: {regexp: "/.*%s.*/i"}}){
			 userId
			  userName
			}
		  }
		  `, username)}

	response, err := client.Do(query)

	if err != nil {
		return model.User{}
	}

	var user model.User
	mapstructure.Decode(response["queryUser"], &user)
	return user
}
