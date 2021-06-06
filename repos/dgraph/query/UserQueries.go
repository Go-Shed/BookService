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

func (repo UserRepo) FetchProfile(profileUserId, userId string) (model.User, error) {
	client := repo.client
	query := dgraph.Request{
		Query: fmt.Sprintf(`query {
			getUser(userId: "%s") {
    userId
    userName
    userPhoto
    email
    followersAggregate{
      count
    }
    followingAggregate{
      count
    }
    followers(filter: {userId: {eq: "%s"}}){
				  userId
				}
			  
	}
}`, profileUserId, userId)}

	response, err := client.Do(query)

	if err != nil {
		return model.User{}, err
	}

	var user model.User
	mapstructure.Decode(response["getUser"], &user)
	return user, nil
}

func (repo UserRepo) FollowUser(user, userToFollow string) error {

	client := repo.client
	query := dgraph.Request{
		Query: fmt.Sprintf(`mutation {
			updateUser(input: {filter: {userId: {eq: "%s"}}, set: {following: {userId: "%s"}}}){
			  user{
				userId
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
			updateUser(input: {filter: {userId: {eq: "%s"}}, remove: {following: {userId: "%s"}}}){
			  user{
				userName
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

func (repo UserRepo) SearchUser(userId, username string) []model.User {

	client := repo.client
	query := dgraph.Request{
		Query: fmt.Sprintf(`{
			queryUser(filter: {userName: {regexp: "/.*%s.*/i"}}){
			 userId
			  userName
			  userPhoto
			  followers(filter: {userId: {eq: "%s"}}){
				userId
			  }
			}
		  }
		  `, username, userId)}

	response, err := client.Do(query)

	if err != nil {
		return []model.User{}
	}

	var user []model.User
	mapstructure.Decode(response["queryUser"], &user)
	return user
}

func (repo UserRepo) GetFollowers(userId string) (model.User, error) {

	client := repo.client
	query := dgraph.Request{
		Query: fmt.Sprintf(`query MyQuery {
			getUser(userId: "%s") {
			  followers{
				userId
				userName
				userPhoto
			  }
			}
		  }
		  
		  `, userId)}

	response, err := client.Do(query)

	if err != nil {
		return model.User{}, err
	}

	var user model.User
	mapstructure.Decode(response["getUser"], &user)
	return user, nil
}

func (repo UserRepo) GetFollowing(userId string) (model.User, error) {

	client := repo.client
	query := dgraph.Request{
		Query: fmt.Sprintf(`query MyQuery {
			getUser(userId: "%s") {
			  following{
				userId
				userName
				userPhoto
			  }
			}
		  }
		  `, userId)}

	response, err := client.Do(query)

	if err != nil {
		return model.User{}, err
	}

	var user model.User
	mapstructure.Decode(response["getUser"], &user)
	return user, nil
}
