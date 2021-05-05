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
	return PostRepo{client: dgraph.Dgraph{}}
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

func (repo PostRepo) GetUserFeedHomeScreen(userId string) (model.User, error) {
	client := repo.client
	query := dgraph.Request{
		Query: fmt.Sprintf(`query MyQuery {
			getUser(userId: "%s") {
			  following {
				userId
				userName
				userPhoto
				posts {
				  id
				  text
				  color
				  createdAt
				  likes(filter: {userId: {eq: "%s"}}) {
					userId
				  }
				  likesAggregate{
					  count
				  }
				}
			  }
			}
		}`, userId, userId)}

	response, err := client.Do(query)

	if err != nil {
		return model.User{}, nil
	}

	var user model.User
	mapstructure.Decode(response["getUser"], &user)
	return user, nil
}

func (repo PostRepo) CreatePost(post model.Post) (model.Post, error) {

	client := repo.client
	query := dgraph.Request{
		Query: `mutation addPost($patch: [AddPostInput!]!) {
			addPost(input: $patch) {
			  post {
				id
			  }
			}
		  }`, Variables: dgraph.Variables{Patch: post}, Operation: "addPost"}

	response, err := client.Do(query)

	if err != nil {
		return model.Post{}, err
	}

	data := response["addPost"].(map[string]interface{})
	var createdPost []model.Post
	mapstructure.Decode(data["post"], &createdPost)
	return createdPost[0], nil
}

func (repo PostRepo) UpdatePost(post model.Post) (model.Post, error) {

	client := repo.client
	query := dgraph.Request{
		Query: fmt.Sprintf(`mutation{
			updatePost(input: {filter: {id: "%s"}, set: {text: "%s", updatedAt: "%s"}}){
			  post{
				id
			  }
			}
		  }`, post.Id, post.Text, post.UpdatedAt), Operation: "updateUser"}

	response, err := client.Do(query)

	if err != nil {
		return model.Post{}, err
	}

	data := response["updatePost"].(map[string]interface{})
	var updatedPost []model.Post
	mapstructure.Decode(data["post"], &updatedPost)
	return updatedPost[0], nil
}
