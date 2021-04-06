package query

import (
	"fmt"
	"shed/bookservice/repos/dgraph"
	"shed/bookservice/repos/dgraph/model"

	"github.com/mitchellh/mapstructure"
)

type PostRepo struct {
	client dgraph.Dgraph
}

func NewPostRepo() PostRepo {
	return PostRepo{client: dgraph.Dgraph{
		URL:    "https://billowing-night.ap-south-1.aws.cloud.dgraph.io/graphql",
		Secret: "ZTE4YjRhNGEwYTAxNWRiYjMwMGI4YmEyMjc1ODU5ZmI=",
	}}
}

func (repo PostRepo) GetUserPosts(userId string) model.User {
	client := repo.client
	query := dgraph.Request{
		Query: fmt.Sprintf(`query {
			getUser(username: "%s") {
			  posts{
				  id
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

func (repo PostRepo) GetUserFeeds(userId string) model.User {
	client := repo.client
	query := dgraph.Request{
		Query: fmt.Sprintf(`query {
			getUser(username: "%s") {
				following(offset: 10) {
				  username
				  posts {
					title
					text
				  }
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
